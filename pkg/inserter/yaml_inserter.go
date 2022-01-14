package inserter

import (
	"context"
	"os"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v2"
)

// YAMLInserter - InserterのYAML実装
type YAMLInserter struct {
	ci *CommonInserter
}

// NewYAMLInserter - YAMLInserter constructor
func NewYAMLInserter(ci *CommonInserter) *YAMLInserter {
	return &YAMLInserter{
		ci: ci,
	}
}

// Execute - YAMLを読み込んでアイテムを生成する
func (y *YAMLInserter) Execute(ctx context.Context, cn, path string) error {
	yb, err := os.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("failed to read yaml file: %+v", err)
	}

	ym := new(Model)
	err = yaml.Unmarshal(yb, ym)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal yaml: %w", err)
	}

	for idx, item := range ym.Items {
		err := y.ci.CreateItem(ctx, cn, item.Ref, item.Payload)
		if err != nil {
			return xerrors.Errorf("failed to create item in array (index=%d): %w", idx, err)
		}
	}

	return nil
}
