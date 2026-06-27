package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"path/filepath"
)

type Config struct {
	Port string
	Address string
	DefaultLimit int
}

func GetConfig() (Config, error) {
	c := Config{
		Port: ":7087",
		Address: "localhost:6600",
		DefaultLimit: 100,
	};

	base, err := GetBaseDir();
	if err != nil {
		return Config{}, err;
	}

	path := filepath.Join(base, "config");

	err = parseConfig(path, &c);
	if err != nil {
		return Config{}, err;
	}

	return c, nil;
}

func parseConfig(path string, c *Config) (error) {
	bt, err := os.ReadFile(path);
	if err != nil {
		return err;
	}

	lines := strings.Split(string(bt), "\n");

	for i := range lines {
		line := strings.TrimSpace(lines[i]);
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "port: "):
			c.Port = ":"+line[5:];

		case strings.HasPrefix(line, "address: "):
			c.Address = line[9:];

		case strings.HasPrefix(line, "default-limit: "):
			converted, err := strconv.Atoi(line[15:]);
			if err != nil {
				return err;
			}
			c.DefaultLimit = converted;
		default:
			return fmt.Errorf("unknown option on config: %v", line[i]);
		}
	}
	return nil;
}

func GetBaseDir() (string, error) {
	// xdgConfigHome := os.Getenv("XDG_CONFIG_HOME");
	//
	// dir := "scrob";
	// path := filepath.Join(xdgConfigHome, dir);
	//
	path := "examples";
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err;
	}

	return path, nil;
}

