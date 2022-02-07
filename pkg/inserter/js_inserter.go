// Package inserter - Firestore にダミーデータを追加するためのパッケージ
package inserter

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
	v8 "rogchap.com/v8go"
)

// JSInserter - Inserter for .js
type JSInserter struct {
	ci *CommonInserter
}

// NewJSInserter - JSInserter constructor
func NewJSInserter(ci *CommonInserter) *JSInserter {
	return &JSInserter{
		ci: ci,
	}
}

// Execute - .js を読み込んで登録する
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

	jms := make([]ModelItem, 0)
	err = json.Unmarshal(jb, &jms)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal json: %w", err)
	}

	docPath := make([]string, 0)
	err = j.CreateItem(ctx, append(docPath, cn), jms, make([]int, 0))
	if err != nil {
		return xerrors.Errorf("failed to create item: %w", err)
	}

	return nil
}

//nolint:dupl
// CreateItem - Firestore にアイテムを生成する
func (j *JSInserter) CreateItem(ctx context.Context, path []string, items []ModelItem, collectionIndexes []int) error {
	for idx, parentItem := range items {
		nowIndexes := append(collectionIndexes, idx)
		docPath := strings.Join(path, "/")
		err := j.ci.CreateItem(ctx, docPath, parentItem.Ref, parentItem.Payload)
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
			err := j.CreateItem(ctx, append(path, j.ci.refIDs[parentItem.Ref], collectionName), subItems, nowIndexes)
			if err != nil {
				return xerrors.Errorf("failed to create item in array: %w", err)
			}
		}
	}

	return nil
}
