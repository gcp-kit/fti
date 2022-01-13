package inserter

import (
	"context"
	"encoding/json"
	"os"

	"golang.org/x/xerrors"
)

type JSONInserter struct {
	ci *CommonInserter
}

func NewJSONInserter(ci *CommonInserter) *JSONInserter {
	return &JSONInserter{
		ci: ci,
	}
}

func (j *JSONInserter) Execute(ctx context.Context, cn, path string) error {
	jb, err := os.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("failed to read json file: %+v", err)
	}

	jm := new(Model)
	err = json.Unmarshal(jb, jm)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal json: %w", err)
	}

	for idx, item := range jm.Items {
		err := j.ci.CreateItem(ctx, cn, item.Ref, item.Payload)
		if err != nil {
			return xerrors.Errorf("failed to create item in array (index=%d): %w", idx, err)
		}
	}

	return nil
}
