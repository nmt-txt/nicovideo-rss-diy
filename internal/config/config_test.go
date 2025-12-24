package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	path := writeConfigTempFile(t, `{
		"searchQueries": [
	        {"query": "  VOCALOID  "},
	        {"query": "ソフトウェアトーク劇場"}
	    ],
	    "log": "DEBUG"
	}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.Log != "debug" {
		t.Fatalf("expected log debug, got %q", cfg.Log)
	}

	if cfg.SearchQueries[0].Query != "VOCALOID" {
		t.Fatalf("expected trimmed query VOCALOID, got %q", cfg.SearchQueries[0].Query)
	}
}

func TestLoadConfig_InvalidLog(t *testing.T) {
	path := writeConfigTempFile(t, `{
	    "searchQueries": [{"query": "foo"}],
	    "log": "trace"
	}`)
	if _, err := LoadConfig(path); err == nil {
		t.Fatalf("expected error for invalid log")
	}
}

func TestLoadConfig_InvalidQuery(t *testing.T) {
	path := writeConfigTempFile(t, `{
	    "searchQueries": [{"query": ""}],
	    "log": "info"
	}`)
	if _, err := LoadConfig(path); err == nil {
		t.Fatalf("expected error for invalid query")
	}
}

func TestLoadConfig_DefaultLogUsed(t *testing.T) {
	path := writeConfigTempFile(t, `{
	    "searchQueries": [{"query": "foo"}]
	}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}

	if cfg.Log != "info" {
		t.Fatalf("expected log default to info, got %q", cfg.Log)
	}
}

func writeConfigTempFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	return path
}
