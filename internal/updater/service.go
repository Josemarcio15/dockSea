// Package updater implements DockSea self-update on Linux.
//
// Fluxo: consulta as releases do GitHub, baixa o .deb compatível em Go
// (net/http), confere o SHA256 (digest do asset ou arquivo .sha256 publicado)
// e instala via `pkexec apt-get install`, deixando a senha com o polkit.
// Se a elevação falhar, o pacote fica salvo em ~/Downloads como plano B.
package updater

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	repoOwner = "Josemarcio15"
	repoName  = "dockSea"
	apiURL    = "https://api.github.com/repos/" + repoOwner + "/" + repoName + "/releases"
)

// UpdateInfo is returned to the frontend by CheckForUpdates.
type UpdateInfo struct {
	HasUpdate      bool   `json:"hasUpdate"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	ReleaseNotes   string `json:"releaseNotes"`
	ReleaseURL     string `json:"releaseUrl"`
	PublishedAt    string `json:"publishedAt"`
	AssetName      string `json:"assetName"`
	AssetSize      int64  `json:"assetSize"`
	Supported      bool   `json:"supported"`
}

// InstallResult describes the outcome of InstallUpdate.
type InstallResult struct {
	Success bool `json:"success"`
	// Cancelled is true when the user dismissed the polkit dialog.
	Cancelled bool `json:"cancelled"`
	// FallbackPath is where the .deb was saved when automatic install failed.
	FallbackPath string `json:"fallbackPath"`
	Message      string `json:"message"`
}

type ghAsset struct {
	Name   string `json:"name"`
	URL    string `json:"browser_download_url"`
	Size   int64  `json:"size"`
	Digest string `json:"digest"`
}

type ghRelease struct {
	TagName     string    `json:"tag_name"`
	Body        string    `json:"body"`
	HTMLURL     string    `json:"html_url"`
	PublishedAt string    `json:"published_at"`
	Draft       bool      `json:"draft"`
	Assets      []ghAsset `json:"assets"`
}

type UpdaterService struct {
	version string
	client  *http.Client

	mu      sync.Mutex
	pending *ghRelease
	asset   *ghAsset
	busy    bool
}

func NewUpdaterService(version string) *UpdaterService {
	return &UpdaterService{
		version: strings.TrimSpace(version),
		client:  &http.Client{Timeout: 10 * time.Minute},
	}
}

// GetVersion returns the running application version.
func (s *UpdaterService) GetVersion() string { return s.version }

// CheckForUpdates queries GitHub for the newest release that ships a .deb for this arch.
func (s *UpdaterService) CheckForUpdates() (*UpdateInfo, error) {
	info := &UpdateInfo{
		CurrentVersion: s.version,
		ReleaseURL:     "https://github.com/" + repoOwner + "/" + repoName + "/releases",
		Supported:      isSupported(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "DockSea-Updater")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("não foi possível acessar o GitHub: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub respondeu %s", resp.Status)
	}

	var releases []ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("resposta inválida do GitHub: %w", err)
	}

	var best *ghRelease
	var bestAsset *ghAsset
	for i := range releases {
		r := &releases[i]
		if r.Draft {
			continue
		}
		a := findDebAsset(r.Assets)
		if a == nil {
			continue
		}
		if best == nil || compareVersions(r.TagName, best.TagName) > 0 {
			best, bestAsset = r, a
		}
	}
	if best == nil {
		return info, nil
	}

	info.LatestVersion = strings.TrimPrefix(strings.TrimPrefix(best.TagName, "v"), "V")
	info.ReleaseNotes = strings.TrimSpace(best.Body)
	info.PublishedAt = best.PublishedAt
	info.AssetName = bestAsset.Name
	info.AssetSize = bestAsset.Size
	if best.HTMLURL != "" {
		info.ReleaseURL = best.HTMLURL
	}
	info.HasUpdate = compareVersions(best.TagName, s.version) > 0

	s.mu.Lock()
	s.pending, s.asset = best, bestAsset
	s.mu.Unlock()
	return info, nil
}

// InstallUpdate downloads, verifies and installs the release found by CheckForUpdates.
func (s *UpdaterService) InstallUpdate() (*InstallResult, error) {
	if !isSupported() {
		return nil, errors.New("auto-update disponível apenas no Linux (Debian/Ubuntu)")
	}
	s.mu.Lock()
	if s.busy {
		s.mu.Unlock()
		return nil, errors.New("uma atualização já está em andamento")
	}
	release, asset := s.pending, s.asset
	if release == nil || asset == nil {
		s.mu.Unlock()
		return nil, errors.New("verifique se há atualizações antes de instalar")
	}
	s.busy = true
	s.mu.Unlock()
	defer func() { s.mu.Lock(); s.busy = false; s.mu.Unlock() }()

	dir, err := os.MkdirTemp("", "docksea-update-")
	if err != nil {
		return nil, fmt.Errorf("não foi possível criar pasta temporária: %w", err)
	}
	defer os.RemoveAll(dir)
	// apt roda o download como usuário _apt; precisa conseguir ler o arquivo.
	_ = os.Chmod(dir, 0o755)

	debPath := filepath.Join(dir, sanitizeName(asset.Name))
	sum, err := s.download(asset.URL, debPath, asset.Size)
	if err != nil {
		return nil, err
	}
	_ = os.Chmod(debPath, 0o644)

	expected, err := s.expectedChecksum(release, asset)
	if err != nil {
		return nil, err
	}
	if expected != "" && !strings.EqualFold(expected, sum) {
		return nil, fmt.Errorf("checksum não confere (esperado %s, obtido %s); download descartado", expected, sum)
	}

	installErr, cancelled := runPkexecInstall(debPath)
	if installErr == nil {
		return &InstallResult{Success: true, Message: "Atualização instalada. Reinicie o DockSea para usar a nova versão."}, nil
	}
	if cancelled {
		return &InstallResult{Cancelled: true, Message: "Autenticação cancelada."}, nil
	}

	// Plano B: guarda o pacote para instalação manual.
	fallback, ferr := saveFallback(debPath, asset.Name)
	if ferr != nil {
		return nil, fmt.Errorf("falha na instalação: %v (e não foi possível salvar o pacote: %v)", installErr, ferr)
	}
	return &InstallResult{
		FallbackPath: fallback,
		Message:      fmt.Sprintf("Instalação automática falhou: %v", installErr),
	}, nil
}

// Restart launches the (new) binary and quits the current process.
func (s *UpdaterService) Restart() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	// Após o dpkg substituir o arquivo, /proc/self/exe aponta para "(deleted)".
	exe = strings.TrimSuffix(exe, " (deleted)")
	cmd := exec.Command(exe, os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("não foi possível reiniciar: %w", err)
	}
	go func() {
		time.Sleep(300 * time.Millisecond)
		os.Exit(0)
	}()
	return nil
}

func (s *UpdaterService) download(url, dst string, size int64) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "DockSea-Updater")
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("falha no download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("falha no download: %s", resp.Status)
	}

	f, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("não foi possível gravar o pacote: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), resp.Body)
	if err != nil {
		return "", fmt.Errorf("download interrompido: %w", err)
	}
	if size > 0 && n != size {
		return "", fmt.Errorf("download incompleto (%d de %d bytes)", n, size)
	}
	if err := f.Sync(); err != nil {
		return "", fmt.Errorf("falha ao gravar o pacote: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// expectedChecksum prefers GitHub's asset digest, falling back to a published .sha256 file.
func (s *UpdaterService) expectedChecksum(r *ghRelease, a *ghAsset) (string, error) {
	if d, ok := strings.CutPrefix(a.Digest, "sha256:"); ok && d != "" {
		return d, nil
	}
	for _, other := range r.Assets {
		if other.Name != a.Name+".sha256" && !strings.EqualFold(other.Name, "SHA256SUMS") {
			continue
		}
		req, err := http.NewRequest(http.MethodGet, other.URL, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", "DockSea-Updater")
		resp, err := s.client.Do(req)
		if err != nil {
			return "", fmt.Errorf("falha ao baixar checksum: %w", err)
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		resp.Body.Close()
		if err != nil || resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("falha ao baixar checksum")
		}
		for _, line := range strings.Split(string(body), "\n") {
			fields := strings.Fields(line)
			if len(fields) == 0 {
				continue
			}
			if len(fields) == 1 || strings.TrimPrefix(fields[len(fields)-1], "*") == a.Name {
				return fields[0], nil
			}
		}
	}
	return "", nil
}

// runPkexecInstall returns (err, cancelled).
func runPkexecInstall(debPath string) (error, bool) {
	pkexec, err := exec.LookPath("pkexec")
	if err != nil {
		return errors.New("pkexec não encontrado (instale o policykit-1)"), false
	}
	var args []string
	if apt, err := exec.LookPath("apt-get"); err == nil {
		args = []string{apt, "install", "-y", "--allow-downgrades", debPath}
	} else if dpkg, err := exec.LookPath("dpkg"); err == nil {
		args = []string{dpkg, "-i", debPath}
	} else {
		return errors.New("apt-get/dpkg não encontrados; sistema não baseado em Debian"), false
	}

	cmd := exec.Command(pkexec, args...)
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil, false
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		switch exitErr.ExitCode() {
		case 126:
			return errors.New("autenticação cancelada"), true
		case 127:
			return errors.New("não foi possível autenticar (nenhum agente polkit disponível?)"), false
		}
	}
	msg := strings.TrimSpace(string(out))
	if lines := strings.Split(msg, "\n"); len(lines) > 6 {
		msg = strings.Join(lines[len(lines)-6:], "\n")
	}
	if msg == "" {
		msg = err.Error()
	}
	return errors.New(msg), false
}

func saveFallback(src, name string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "Downloads")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	dst := filepath.Join(dir, sanitizeName(name))
	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return "", err
	}
	return dst, nil
}

func isSupported() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	_, err := exec.LookPath("dpkg")
	return err == nil
}

var archAliases = map[string][]string{
	"amd64": {"amd64", "x86_64"},
	"arm64": {"arm64", "aarch64"},
}

func findDebAsset(assets []ghAsset) *ghAsset {
	aliases := archAliases[runtime.GOARCH]
	if aliases == nil {
		aliases = []string{runtime.GOARCH}
	}
	var generic *ghAsset
	for i := range assets {
		a := &assets[i]
		name := strings.ToLower(a.Name)
		if !strings.HasSuffix(name, ".deb") {
			continue
		}
		for _, al := range aliases {
			if strings.Contains(name, al) {
				return a
			}
		}
		// .deb sem arquitetura no nome: aceita só se não mencionar outra arch.
		if generic == nil && !mentionsOtherArch(name) {
			generic = a
		}
	}
	return generic
}

func mentionsOtherArch(name string) bool {
	for arch, als := range archAliases {
		if arch == runtime.GOARCH {
			continue
		}
		for _, al := range als {
			if strings.Contains(name, al) {
				return true
			}
		}
	}
	return false
}

var unsafeChars = regexp.MustCompile(`[^A-Za-z0-9._-]`)

func sanitizeName(name string) string {
	n := unsafeChars.ReplaceAllString(filepath.Base(name), "_")
	if n == "" || n == "." || n == ".." {
		return "docksea-update.deb"
	}
	return n
}
