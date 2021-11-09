package inserter

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-generalize/fti/pkg/config"
	"github.com/go-generalize/fti/pkg/files"
	"golang.org/x/xerrors"
)

type Inserter struct {
	client *firestore.Client
	refIDs map[string]string
}

func NewInserter(client *firestore.Client) *Inserter {
	return &Inserter{
		client: client,
	}
}

func (i *Inserter) Execute(ctx context.Context, cfg *config.Config) error {
	targetDir := cfg.Targets
	for _, t := range targetDir {
		if !files.Exists(t) {
			return xerrors.Errorf("cannot find target directory: %s", targetDir)
		}

		err := filepath.Walk(t, i.executeFile(ctx))
		if err != nil {
			return xerrors.Errorf("failed to execute file: %w", err)
		}
	}

	return nil
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
			err = i.executeJSON(ctx, cn, path)
			if err != nil {
				log.Printf("failed to insert json file: %s\n%+v", path, err)
			}
		case strings.HasSuffix(path, ".js"):
			// TODO: implements here for js api
		}

		return nil
	}
}

func (i *Inserter) executeJSON(ctx context.Context, cn, path string) error {
	jb, err := os.ReadFile(path)
	if err != nil {
		return xerrors.Errorf("failed to read json file: %+v", err)
	}

	jm := new(JsonModel)
	err = json.Unmarshal(jb, jm)
	if err != nil {
		return xerrors.Errorf("failed to unmarshal json: %w", err)
	}

	if payload, ok := jm.Payload.([]interface{}); ok {
		for idx, p := range payload {
			mp, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			err := i.createItem(ctx, cn, mp)
			if err != nil {
				return xerrors.Errorf("failed to create item in array (index=%d): %w", idx, err)
			}
		}
	} else if payload, ok := jm.Payload.(map[string]interface{}); ok {
		err := i.createItem(ctx, cn, payload)
		if err != nil {
			return xerrors.Errorf("failed to create item: %w", err)
		}
	} else {
		// print log or error?
	}

	return nil
}

func (i *Inserter) createItem(ctx context.Context, cn string, item map[string]interface{}) error {
	item = i.tryParseDate(item)

	log.Printf("collection: %s, item: %+v", cn, item)
	_, err := i.client.Collection(cn).NewDoc().Create(ctx, item)
	if err != nil {
		return xerrors.Errorf("failed to create item: %w", err)
	}
	return nil
}

func (i *Inserter) tryParseDate(item map[string]interface{}) map[string]interface{} {
	for k, v := range item {
		switch vt := v.(type) {
		case string:
			pt, err := time.Parse(time.RFC3339, vt)
			if err != nil {
				// print log?
				continue
			}
			item[k] = pt

		case map[string]interface{}:
			item[k] = i.tryParseDate(vt)
		}
	}

	return item
}
