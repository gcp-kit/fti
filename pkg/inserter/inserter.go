// Package inserter - Firestore にダミーデータを追加するためのパッケージ
package inserter

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gcp-kit/fti/pkg/config"
	"github.com/gcp-kit/fti/pkg/files"
	"golang.org/x/xerrors"
)

// Inserter - Inserter
type Inserter struct {
	commonInserter *CommonInserter

	jsonInserter *JSONInserter
	yamlInserter *YAMLInserter
	jsInserter   *JSInserter
}

// NewInserter - Inserter constructor
func NewInserter(client *firestore.Client, refIDs RefIDs) *Inserter {
	ci := NewCommonInserter(client, refIDs)
	return &Inserter{
		commonInserter: ci,
		jsonInserter:   NewJSONInserter(ci),
		jsInserter:     NewJSInserter(ci),
		yamlInserter:   NewYAMLInserter(ci),
	}
}

// Execute - ダミーデータソースの読み込みと実行処理呼び出しを行う
func (i *Inserter) Execute(ctx context.Context, cfg *config.Config) (RefIDs, error) {
	targetDir := cfg.Targets
	for _, t := range targetDir {
		if !files.Exists(t) {
			return nil, xerrors.Errorf("cannot find target directory: %s", targetDir)
		}

		err := filepath.Walk(t, i.executeFile(ctx))
		if err != nil {
			return nil, xerrors.Errorf("failed to execute file: %w", err)
		}
	}

	refIDs := i.commonInserter.RefIDs()

	return refIDs, nil
}

func (i *Inserter) executeFile(ctx context.Context) func(path string, info os.FileInfo, _ error) error {
	return func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		cn := filepath.Base(filepath.Dir(path))

		var err error
		switch {
		case strings.HasSuffix(path, ".json"):
			err = i.jsonInserter.Execute(ctx, cn, path)
			if err != nil {
				log.Printf("failed to insert json file: %s\n%+v", path, err)
			}
		case strings.HasSuffix(path, ".yaml"):
		case strings.HasSuffix(path, ".yml"):
			err = i.yamlInserter.Execute(ctx, cn, path)
			if err != nil {
				log.Printf("failed to insert yaml file: %s\n%+v", path, err)
			}
		case strings.HasSuffix(path, ".js"):
			err = i.jsInserter.Execute(ctx, cn, path)
			if err != nil {
				log.Printf("failed to insert js file: %s\n%+v", path, err)
			}
		}

		return nil
	}
}
