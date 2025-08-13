//go:build linux

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Simple CLI tool on Linux that reads ~/.config/volume-speeder/config for settings
// Supported keys:
// - multiplier: integer >=1 used for how many times to repeat volume key press
// Commands:
// - get: prints current multiplier
// - set <n>: sets multiplier and saves config
// - inc: increments multiplier by 1 and saves
// - dec: decrements multiplier by 1 (min 1) and saves
// - help: prints usage

func main() {
	cfg, _ := loadConfig()
	args := os.Args[1:]
	if len(args) == 0 || args[0] == "help" || args[0] == "-h" || args[0] == "--help" {
		printHelp()
		return
	}

	switch args[0] {
	case "get":
		fmt.Println(cfg.Multiplier)
	case "set":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "set requires a value")
			os.Exit(2)
		}
		n, err := strconv.Atoi(args[1])
		if err != nil || n < 1 {
			fmt.Fprintln(os.Stderr, "multiplier must be an integer >= 1")
			os.Exit(2)
		}
		cfg.Multiplier = n
		if err := saveConfig(cfg); err != nil {
			fmt.Fprintln(os.Stderr, "failed to save:", err)
			os.Exit(1)
		}
		fmt.Println(n)
	case "inc":
		cfg.Multiplier++
		if cfg.Multiplier < 1 {
			cfg.Multiplier = 1
		}
		_ = saveConfig(cfg)
		fmt.Println(cfg.Multiplier)
	case "dec":
		cfg.Multiplier--
		if cfg.Multiplier < 1 {
			cfg.Multiplier = 1
		}
		_ = saveConfig(cfg)
		fmt.Println(cfg.Multiplier)
	default:
		printHelp()
	}
}

type Config struct {
	Multiplier int
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil { return "", err }
	return filepath.Join(home, ".config", "volume-speeder", "config"), nil
}

func loadConfig() (Config, error) {
	cfg := Config{Multiplier: 3}
	path, err := configPath()
	if err != nil { return cfg, err }
	f, err := os.Open(path)
	if err != nil { return cfg, err }
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 { continue }
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if strings.EqualFold(k, "multiplier") {
			if n, err := strconv.Atoi(v); err == nil && n >= 1 { cfg.Multiplier = n }
		}
	}
	return cfg, nil
}

func saveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil { return err }
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil { return err }
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	_, err = fmt.Fprintf(f, "multiplier=%d\n", cfg.Multiplier)
	return err
}

func printHelp() {
	fmt.Println("Usage: volume-speeder [get|set <n>|inc|dec|help]")
}
