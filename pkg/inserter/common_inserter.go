package inserter

import (
	"context"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/xerrors"
)

type CommonInserter struct {
	client *firestore.Client
	refIDs map[string]string
}

func NewCommonInserter(client *firestore.Client) *CommonInserter {
	return &CommonInserter{
		client: client,
		refIDs: map[string]string{},
	}
}

func (c *CommonInserter) CreateItem(ctx context.Context, cn, refID string, item map[string]interface{}) error {
	item = c.tryParseDate(item)
	item = c.setRefs(item)

	d := c.client.Collection(cn).NewDoc()
	_, err := d.Create(ctx, item)
	if err != nil {
		return xerrors.Errorf("failed to create item: %w", err)
	}

	if refID != "" {
		if _, ok := c.refIDs[refID]; ok {
			return xerrors.Errorf("already ref id: %s", refID)
		}
		c.refIDs[refID] = d.ID
	}

	return nil
}

func (ci *CommonInserter) tryParseDate(item map[string]interface{}) map[string]interface{} {
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
			item[k] = ci.tryParseDate(vt)
		}
	}

	return item
}

func (ci *CommonInserter) setRefs(item map[string]interface{}) map[string]interface{} {
	for k, v := range item {
		switch vt := v.(type) {
		case string:
			if strings.HasPrefix(vt, "$") {
				refID := strings.TrimPrefix(vt, "$")
				rv, ok := ci.refIDs[refID]
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
