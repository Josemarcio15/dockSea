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
		reNsh := regexp.MustCompile(`!define INFO_PRODUCTVERSION "[^"]*"`)
		content = reNsh.ReplaceAllString(content, fmt.Sprintf(`!define INFO_PRODUCTVERSION "%s"`, version))

		_ = os.WriteFile(nshPath, []byte(content), 0644)
	}

	// Output version as the result for Taskfile consumption
	fmt.Print(version)
}
