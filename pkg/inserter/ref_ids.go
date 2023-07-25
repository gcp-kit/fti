package inserter

import (
	"os"

	"github.com/goccy/go-yaml"
	"golang.org/x/xerrors"
)

type RefIDs map[string]string

// LoadFromFile - yamlファイルからrefIDを読み込む
func (r RefIDs) LoadFromFile(path string) error {
	// pathのyamlファイル読み込む
	yb, err := os.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("failed to read yaml file: %+v", err)
	}

	// yamlファイルを構造体に変換する
	ym := new(RefIDs)
	err = yaml.Unmarshal(yb, ym)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal yaml: %w", err)
	}
	// 構造体の値をrにコピーする
	for k, v := range *ym {
		r[k] = v
	}

	return nil
}

// SaveToFile - yamlファイルにrefIDを書き込む
func (r RefIDs) SaveToFile(path string) error {
	// yamlに変換する
	yb, err := yaml.Marshal(r)
	if err != nil {
		return xerrors.Errorf("failed to marshal yaml: %w", err)
	}

	// yamlファイルに書き込む
	err = os.WriteFile(path, yb, 0666)
	if err != nil {
		return xerrors.Errorf("failed to write yaml file: %w", err)
	}

	return nil
}
