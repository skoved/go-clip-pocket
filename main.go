// Copyright skoved
// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	// set at compile time (treat as const)
	version string
	cmdName string

	// cmd flags
	helpFlag    bool
	versionFlag bool
)

const (
	// flag usage info
	helpUsage    = "print the help info"
	versionUsage = "print the version"

	clipboardFilePath = ".local/share/clip-pocket"
)

// Set cmd flags and parse the values
func cmdFlags() {
	flag.BoolVar(&helpFlag, "help", false, helpUsage)
	flag.BoolVar(&helpFlag, "h", false, helpUsage+"(shorthand)")
	flag.BoolVar(&versionFlag, "version", false, versionUsage)
	flag.BoolVar(&versionFlag, "v", false, versionUsage+"(shorthand)")
}

// Print the version and copyright and license
func printVersion() {
	fmt.Fprintln(flag.CommandLine.Output(), cmdName, version)
	fmt.Fprintln(flag.CommandLine.Output(), "Copyright (C) 2025 Sam Koved")
	fmt.Fprintln(flag.CommandLine.Output(), "License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>")
	fmt.Fprintln(
		flag.CommandLine.Output(),
		"This program comes with ABSOLUTELY NO WARRANTY;\nThis is free software, and you are welcome to redistribute it",
	)
}

func getclipBoardFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + clipboardFilePath, nil
}

func main() {
	cmdFlags()
	flag.Parse()

	if helpFlag {
		flag.Usage()
		os.Exit(0)
	}
	if versionFlag {
		printVersion()
		os.Exit(0)
	}

	clipboardFile, err := getclipBoardFile()
	if err != nil {
		fmt.Fprintln(flag.CommandLine.Output(), "Could not determine home directory:", err)
	}

	clFile, openErr := os.OpenFile(clipboardFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if openErr != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "could not open or create clipboard file %s\nError:\n%s\n", clipboardFile, openErr)
	}

	var buf bytes.Buffer

	copied, readErr := io.ReadAll(os.Stdin)
	if readErr != nil {
		fmt.Fprintln(flag.CommandLine.Output(), "Could not read from stdin:", readErr)
	}
	if bytes.Equal(copied, []byte("")) {
		os.Exit(0)
	}
	buf.Write(copied)
	buf.WriteString("\n")

	if _, writeErr := clFile.Write(buf.Bytes()); writeErr != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "could not write to clipboard file %s\nError:\n%s\n", clipboardFile, writeErr)
	}
}
