package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// filters追加など拡張性確保のため
type SearchQuery struct {
	Query string `json:"query"`
}

type System struct {
	Version string
}

type Config struct {
	SearchQueries []SearchQuery `json:"searchQueries"`
	Log           string        `json:"log,omitempty"`
	System        System        `json:"-"`
}

// LoadConfig 設定ファイルを読み込み、検証して返す。
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルの解析に失敗しました: %w", err)
	}

	logLevel := strings.ToLower(strings.TrimSpace(cfg.Log))
	if logLevel == "" {
		logLevel = "info"
	}
	cfg.Log = logLevel
	switch logLevel {
	case "debug", "info", "error":
	default:
		return nil, fmt.Errorf("logはdebug/info/errorのいずれかである必要があります。")
	}

	for i := range cfg.SearchQueries {
		trimmed := strings.TrimSpace(cfg.SearchQueries[i].Query)
		if trimmed == "" {
			return nil, fmt.Errorf("検索タグ内容を空にすることはできません。APIガイドを参照してください(https://site.nicovideo.jp/search-api-docs/snapshot)。(任意のfilters併用は未対応です)")
		}
		cfg.SearchQueries[i].Query = trimmed
	}

	cfg.System.Version = "1.0.0"
	return &cfg, nil
}
