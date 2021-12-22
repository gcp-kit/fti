package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	jm := new(JsonModel)
	err = json.Unmarshal(jb, jm)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal json: %w", err)
	}

	err = j.CreateItem(ctx, cn, jm.Items, make([]int, 0))
	if err != nil {
		return xerrors.Errorf("failed to create item: %w", err)
	}

	return nil
}

func (j *JSONInserter) CreateItem(ctx context.Context, cn string, items []JsonModelItem, collectionIndexes []int) error {
	for idx, parentItem := range items {
		nowIndexes := append(collectionIndexes, idx)
		err := j.ci.CreateItem(ctx, cn, parentItem.Ref, parentItem.Payload)
		if err != nil {
			errorIndexes := make([]string, 0)
			for _, v := range nowIndexes {
				errorIndexes = append(errorIndexes, strconv.Itoa(v))
			}
			return xerrors.Errorf("failed to create item in array (index=%s): %w", strings.Join(errorIndexes, "/"), err)
		}
		if parentItem.SubCollections == nil || len(parentItem.SubCollections) == 0 {
			continue
		}
		for collectionName, subItems := range parentItem.SubCollections {
			err := j.CreateItem(ctx, fmt.Sprintf("%s/%s/%s", cn, j.ci.refIDs[parentItem.Ref], collectionName), subItems, nowIndexes)
			if err != nil {
				return xerrors.Errorf("failed to create item in array: %w", err)
			}
		}
	}

	return nil
}
