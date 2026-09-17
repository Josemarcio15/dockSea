package connection

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// ExecCommand executa um comando no servidor com timeout padrão seguro de 2 minutos
func (c *Client) ExecCommand(cmd string, useSudo bool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	return c.ExecCommandContext(ctx, cmd, useSudo)
}

// ExecCommandContext executa um comando no servidor respeitando o cancelamento/timeout do Context
func (c *Client) ExecCommandContext(ctx context.Context, cmd string, useSudo bool) (string, error) {
	finalCmd, sudoPassword := c.prepareCmd(cmd, useSudo)

	if isLocal(c.server) {
		cmdObj := exec.CommandContext(ctx, "bash", "-c", finalCmd)
		if sudoPassword != "" {
			cmdObj.Stdin = strings.NewReader(sudoPassword + "\n")
		}
		out, err := cmdObj.CombinedOutput()
		if err != nil {
			return string(out), fmt.Errorf("comando local falhou (%w): %s", err, string(out))
		}
		return string(out), nil
	}

	sc := c.GetSSHClient()
	if sc == nil {
		return "", fmt.Errorf("cliente SSH não conectado")
	}

	session, err := sc.NewSession()
	if err != nil {
		return "", fmt.Errorf("falha ao criar sessão SSH: %w", err)
	}
	defer session.Close()

	// Injeta a senha do sudo pelo stdin com segurança sem expor na tabela de processos (/proc)
	if sudoPassword != "" {
		stdinPipe, err := session.StdinPipe()
		if err != nil {
			return "", fmt.Errorf("falha ao abrir stdin para sudo: %w", err)
		}
		go func() {
			_, _ = io.WriteString(stdinPipe, sudoPassword+"\n")
			_ = stdinPipe.Close()
		}()
	}

	// Goroutine para escutar cancelamento de contexto (SSH-003 / SSH-008)
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			_ = session.Signal(ssh.SIGKILL)
			_ = session.Close()
		case <-done:
		}
	}()

	out, err := session.CombinedOutput(finalCmd)
	if ctx.Err() != nil {
		return "", fmt.Errorf("comando cancelado ou timeout atingido: %w", ctx.Err())
	}
	if err != nil {
		return string(out), fmt.Errorf("comando falhou (%w): %s", err, string(out))
	}

	return string(out), nil
}

// StartCommandOutput inicia um comando e retorna seu stdout como io.Reader para streaming contínuo
func (c *Client) StartCommandOutput(ctx context.Context, cmd string, useSudo bool) (io.Reader, func() error, error) {
	finalCmd, sudoPassword := c.prepareCmd(cmd, useSudo)

	if isLocal(c.server) {
		cmdObj := exec.CommandContext(ctx, "bash", "-c", finalCmd)
		if sudoPassword != "" {
			cmdObj.Stdin = strings.NewReader(sudoPassword + "\n")
		}
		stdout, err := cmdObj.StdoutPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao abrir stdout pipe local: %w", err)
		}
		stderr, err := cmdObj.StderrPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao abrir stderr pipe local: %w", err)
		}
		if err := cmdObj.Start(); err != nil {
			return nil, nil, fmt.Errorf("falha ao iniciar comando local: %w", err)
		}
		waitFn := func() error {
			return cmdObj.Wait()
		}
		return io.MultiReader(stdout, stderr), waitFn, nil
	}

	sc := c.GetSSHClient()
	if sc == nil {
		return nil, nil, fmt.Errorf("cliente SSH não conectado")
	}

	session, err := sc.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao criar sessão SSH: %w", err)
	}

	// Injeta a senha do sudo de forma segura se necessário
	if sudoPassword != "" {
		stdinPipe, err := session.StdinPipe()
		if err == nil {
			go func() {
				_, _ = io.WriteString(stdinPipe, sudoPassword+"\n")
				_ = stdinPipe.Close()
			}()
		}
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao abrir stdout pipe SSH: %w", err)
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao abrir stderr pipe SSH: %w", err)
	}

	pr, pw := io.Pipe()

	if err := session.Start(finalCmd); err != nil {
		_ = session.Close()
		_ = pr.Close()
		_ = pw.Close()
		return nil, nil, fmt.Errorf("falha ao iniciar comando SSH: %w", err)
	}

	// Goroutine monitorando cancelamento de contexto para encerrar o processo remoto
	cancelDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = session.Signal(ssh.SIGKILL)
			_ = session.Close()
		case <-cancelDone:
		}
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		_, _ = io.Copy(pw, stdout)
	}()

	go func() {
		defer wg.Done()
		_, _ = io.Copy(pw, stderr)
	}()

	go func() {
		wg.Wait()
		_ = pw.Close()
	}()

	waitFn := func() error {
		defer close(cancelDone)
		defer session.Close()
		defer pr.Close()
		if err := session.Wait(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("execução cancelada: %w", ctx.Err())
			}
			return err
		}
		return nil
	}

	return pr, waitFn, nil
}

// StartCommandInput inicia um comando e retorna seu stdin como io.WriteCloser para recepção de dados via streaming
func (c *Client) StartCommandInput(ctx context.Context, cmd string, useSudo bool) (io.WriteCloser, func() error, error) {
	finalCmd, sudoPassword := c.prepareCmd(cmd, useSudo)

	if isLocal(c.server) {
		cmdObj := exec.CommandContext(ctx, "bash", "-c", finalCmd)
		stdin, err := cmdObj.StdinPipe()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao abrir stdin pipe local: %w", err)
		}
		if err := cmdObj.Start(); err != nil {
			return nil, nil, fmt.Errorf("falha ao iniciar comando receptor local: %w", err)
		}
		waitFn := func() error {
			return cmdObj.Wait()
		}
		return stdin, waitFn, nil
	}

	sc := c.GetSSHClient()
	if sc == nil {
		return nil, nil, fmt.Errorf("cliente SSH não conectado")
	}

	session, err := sc.NewSession()
	if err != nil {
		return nil, nil, fmt.Errorf("falha ao criar sessão SSH: %w", err)
	}

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf

	stdin, err := session.StdinPipe()
	if err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao abrir stdin pipe SSH: %w", err)
	}

	if err := session.Start(finalCmd); err != nil {
		_ = session.Close()
		return nil, nil, fmt.Errorf("falha ao iniciar comando receptor SSH: %w", err)
	}

	// Goroutine monitorando cancelamento de contexto
	cancelDone := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = session.Signal(ssh.SIGKILL)
			_ = session.Close()
		case <-cancelDone:
		}
	}()

	waitFn := func() error {
		defer close(cancelDone)
		defer session.Close()
		if err := session.Wait(); err != nil {
			if ctx.Err() != nil {
				return fmt.Errorf("execução cancelada: %w", ctx.Err())
			}
			if stderrBuf.Len() > 0 {
				return fmt.Errorf("%w: %s", err, stderrBuf.String())
			}
			return err
		}
		return nil
	}

	// Se houver sudoPassword e o comando for receptor de streaming, o streaming direto em stdin deve considerar isso
	_ = sudoPassword
	return stdin, waitFn, nil
}

// prepareCmd sanitiza o comando e determina se precisa de injeção de senha via stdin (SSH-005)
func (c *Client) prepareCmd(cmd string, useSudo bool) (string, string) {
	finalCmd := cmd
	var sudoPassword string

	if useSudo {
		username := strings.TrimSpace(strings.ToLower(c.server.Username))
		if username == "root" {
			finalCmd = cmd
		} else if strings.TrimSpace(c.server.SudoPassword) != "" {
			// sudo -S lê a senha do STDIN com prompt vazio, sem expor a senha em /proc/<pid>/cmdline
			finalCmd = fmt.Sprintf("sudo -S -p '' bash -c %s", escapeShell(cmd))
			sudoPassword = c.server.SudoPassword
		} else {
			if !strings.HasPrefix(cmd, "sudo ") {
				finalCmd = fmt.Sprintf("sudo -n bash -c %s", escapeShell(cmd))
			}
		}
	}
	return finalCmd, sudoPassword
}
