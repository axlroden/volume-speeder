//go:build windows || darwin

package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type Config struct {
    Multiplier int
}

func configPath() (string, error) {
    // Prefer OS config dir when available
    if base, err := os.UserConfigDir(); err == nil && base != "" {
        return filepath.Join(base, "volume-speeder", "config"), nil
    }
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".config", "volume-speeder", "config"), nil
}

func loadConfig() (Config, error) {
    cfg := Config{Multiplier: 3}
    path, err := configPath()
    if err != nil {
        return cfg, err
    }
    f, err := os.Open(path)
    if err != nil {
        return cfg, err
    }
    defer f.Close()
    s := bufio.NewScanner(f)
    for s.Scan() {
        line := strings.TrimSpace(s.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.SplitN(line, "=", 2)
        if len(parts) != 2 {
            continue
        }
        k := strings.TrimSpace(parts[0])
        v := strings.TrimSpace(parts[1])
        if strings.EqualFold(k, "multiplier") {
            if n, err := strconv.Atoi(v); err == nil && n >= 1 {
                cfg.Multiplier = n
            }
        }
    }
    return cfg, nil
}

func saveConfig(cfg Config) error {
    path, err := configPath()
    if err != nil {
        return err
    }
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        return err
    }
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer f.Close()
    _, err = fmt.Fprintf(f, "multiplier=%d\n", cfg.Multiplier)
    return err
}
