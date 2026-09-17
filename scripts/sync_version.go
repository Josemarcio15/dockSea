package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	// Read root VERSION file
	versionBytes, err := os.ReadFile("VERSION")
	if err != nil {
		fmt.Printf("Error reading VERSION file: %v\n", err)
		os.Exit(1)
	}
	version := strings.TrimSpace(string(versionBytes))
	if version == "" {
		fmt.Println("VERSION file is empty")
		os.Exit(1)
	}

	// Update build/config.yml
	configPath := "build/config.yml"
	configBytes, err := os.ReadFile(configPath)
	if err == nil {
		content := string(configBytes)
		reVersion := regexp.MustCompile(`(?m)^version:\s*".*"`)
		content = reVersion.ReplaceAllString(content, fmt.Sprintf(`version: "%s"`, version))

		reProductVersion := regexp.MustCompile(`(?m)^(\s*productVersion:\s*)".*"`)
		content = reProductVersion.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, version))

		_ = os.WriteFile(configPath, []byte(content), 0644)
	}

	// Update build/windows/info.json if exists
	infoPath := "build/windows/info.json"
	infoBytes, err := os.ReadFile(infoPath)
	if err == nil {
		content := string(infoBytes)
		reFileVersion := regexp.MustCompile(`"file_version":\s*"[^"]*"`)
		content = reFileVersion.ReplaceAllString(content, fmt.Sprintf(`"file_version": "%s"`, version))

		reProductVersion := regexp.MustCompile(`"product_version":\s*"[^"]*"`)
		content = reProductVersion.ReplaceAllString(content, fmt.Sprintf(`"product_version": "%s"`, version))

		rePropProductVersion := regexp.MustCompile(`"ProductVersion":\s*"[^"]*"`)
		content = rePropProductVersion.ReplaceAllString(content, fmt.Sprintf(`"ProductVersion": "%s"`, version))

		_ = os.WriteFile(infoPath, []byte(content), 0644)
	}

	// Update build/windows/nsis/wails_tools.nsh if exists
	nshPath := "build/windows/nsis/wails_tools.nsh"
	nshBytes, err := os.ReadFile(nshPath)
	if err == nil {
		content := string(nshBytes)
		reNshVersion := regexp.MustCompile(`(?m)!define\s+INFO_PRODUCTVERSION\s+".*"`)
		content = reNshVersion.ReplaceAllString(content, fmt.Sprintf(`!define INFO_PRODUCTVERSION "%s"`, version))
		_ = os.WriteFile(nshPath, []byte(content), 0644)
	}

	// Update build/linux/nfpm/nfpm.yaml if exists
	nfpmPath := "build/linux/nfpm/nfpm.yaml"
	nfpmBytes, err := os.ReadFile(nfpmPath)
	if err == nil {
		content := string(nfpmBytes)
		reNfpm := regexp.MustCompile(`(?m)^version:\s*".*"`)
		content = reNfpm.ReplaceAllString(content, fmt.Sprintf(`version: "%s"`, version))

		_ = os.WriteFile(nfpmPath, []byte(content), 0644)
	}

	// Update build/windows/wails.exe.manifest if exists
	manifestPath := "build/windows/wails.exe.manifest"
	manifestBytes, err := os.ReadFile(manifestPath)
	if err == nil {
		content := string(manifestBytes)
		// Fix XML header if previously replaced
		content = strings.Replace(content, `<?xml version="0.0.9-alpha"`, `<?xml version="1.0"`, 1)
		content = strings.Replace(content, `Microsoft.Windows.Common-Controls" version="0.0.9-alpha"`, `Microsoft.Windows.Common-Controls" version="6.0.0.0"`, 1)
		reManifest := regexp.MustCompile(`(<assemblyIdentity type="win32" name="com\.samabe\.docksea" version=")[^"]*(")`)
		content = reManifest.ReplaceAllString(content, fmt.Sprintf(`${1}%s${2}`, version))
		_ = os.WriteFile(manifestPath, []byte(content), 0644)
	}

	// Update build/windows/nsis/project.nsi if exists
	nsiPath := "build/windows/nsis/project.nsi"
	nsiBytes, err := os.ReadFile(nsiPath)
	if err == nil {
		content := string(nsiBytes)
		// Extract numeric version (e.g. 0.0.9-alpha -> 0.0.9.0)
		numParts := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`).FindStringSubmatch(version)
		numVer := "0.0.0.0"
		if len(numParts) == 4 {
			numVer = fmt.Sprintf("%s.%s.%s.0", numParts[1], numParts[2], numParts[3])
		}
		reNumVer := regexp.MustCompile(`(?m)!define\s+INFO_NUMERIC_VERSION\s+".*"`)
		content = reNumVer.ReplaceAllString(content, fmt.Sprintf(`!define INFO_NUMERIC_VERSION "%s"`, numVer))
		_ = os.WriteFile(nsiPath, []byte(content), 0644)
	}

	// Output version as the result for Taskfile consumption
	fmt.Print(version)
}
