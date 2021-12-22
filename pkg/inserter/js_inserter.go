package inserter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
	v8 "rogchap.com/v8go"
)

type JSInserter struct {
	ci *CommonInserter
}

func NewJSInserter(ci *CommonInserter) *JSInserter {
	return &JSInserter{
		ci: ci,
	}
}

func (j *JSInserter) Execute(ctx context.Context, cn, path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("failed to js read from file: %w", err)
	}

	v8ctx, err := v8.NewContext()
	if err != nil {
		return xerrors.Errorf("failed to create v8 context: %w", err)
	}

	val, err := v8ctx.RunScript(string(b), path)
	if err != nil {
		return xerrors.Errorf("failed to run script: %s\n%w", path, err)
	}

	if !val.IsArray() {
		return xerrors.Errorf("js returned data must be array")
	}
	obj, err := val.AsObject()
	if err != nil {
		return xerrors.Errorf("failed to convert object by returned value: %w", err)
	}

	jb, err := obj.MarshalJSON()
	if err != nil {
		return xerrors.Errorf("failed to marshal json of returned value: %w", err)
	}

	jms := make([]JsonModelItem, 0)
	err = json.Unmarshal(jb, &jms)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal json: %w", err)
	}

	err = j.CreateItem(ctx, cn, jms, make([]int, 0))
	if err != nil {
		return xerrors.Errorf("failed to create item: %w", err)
	}

	return nil
}

func (j *JSInserter) CreateItem(ctx context.Context, cn string, items []JsonModelItem, collectionIndexes []int) error {
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
		if len(parentItem.SubCollection) == 0 {
			continue
		}
		for collectionName, subItems := range parentItem.SubCollection {
			err := j.CreateItem(ctx, fmt.Sprintf("%s/%s/%s", cn, j.ci.refIDs[parentItem.Ref], collectionName), subItems, nowIndexes)
			if err != nil {
				return xerrors.Errorf("failed to create item in array: %w", err)
			}
		}
	}

	return nil
}
