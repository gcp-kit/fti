package inserter

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/xerrors"
)

type JSONInserter struct {
	client *firestore.Client
	refIDs map[string]string
}

func NewJSONInserter(client *firestore.Client) *JSONInserter {
	return &JSONInserter{
		client: client,
		refIDs: map[string]string{},
	}
}

func (j *JSONInserter) executeJSON(ctx context.Context, cn, path string) error {
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
			err := j.createItem(ctx, cn, jm.Ref, mp)
			if err != nil {
				return xerrors.Errorf("failed to create item in array (index=%d): %w", idx, err)
			}
		}
	} else if payload, ok := jm.Payload.(map[string]interface{}); ok {
		err := j.createItem(ctx, cn, jm.Ref, payload)
		if err != nil {
			return xerrors.Errorf("failed to create item: %w", err)
		}
	} else {
		// print log or error?
	}

	return nil
}

func (j *JSONInserter) createItem(ctx context.Context, cn, refID string, item map[string]interface{}) error {
	item = j.tryParseDate(item)
	item = j.setRefs(item)

	d := j.client.Collection(cn).NewDoc()
	_, err := d.Create(ctx, item)
	if err != nil {
		return xerrors.Errorf("failed to create item: %w", err)
	}

	if refID != "" {
		if _, ok := j.refIDs[refID]; ok {
			return xerrors.Errorf("already ref id: %s", refID)
		}
		j.refIDs[refID] = d.ID
	}

	return nil
}

func (j *JSONInserter) tryParseDate(item map[string]interface{}) map[string]interface{} {
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
			item[k] = j.tryParseDate(vt)
		}
	}

	return item
}

func (j *JSONInserter) setRefs(item map[string]interface{}) map[string]interface{} {
	for k, v := range item {
		switch vt := v.(type) {
		case string:
			if strings.HasPrefix(vt, "$") {
				refID := strings.TrimPrefix(vt, "$")
				rv, ok := j.refIDs[refID]
				if !ok {
					log.Printf("%s was not found", refID)
					continue
				}
				item[k] = rv
			}
		}
	}

	return item
}
