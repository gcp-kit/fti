// Package inserter - Firestore にダミーデータを追加するためのパッケージ
package inserter

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/xerrors"
)

// CommonInserter - Inserterの共通部分
type CommonInserter struct {
	client *firestore.Client
	refIDs map[string]string
}

// NewCommonInserter - CommonInserter constructor
func NewCommonInserter(client *firestore.Client) *CommonInserter {
	return &CommonInserter{
		client: client,
		refIDs: map[string]string{},
	}
}

// CreateItem - item を Firestore に作る
func (c *CommonInserter) CreateItem(ctx context.Context, cn, ID, refID string, item map[string]interface{}) error {
	item = c.tryParseDate(item)
	item = c.setRefs(item)
	log.Print(item)

	var d *firestore.DocumentRef
	if ID == "" {
		d = c.client.Collection(cn).NewDoc()
	} else {
		d = c.client.Collection(cn).Doc(ID)
	}
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

func (c *CommonInserter) tryParseDate(item map[string]interface{}) map[string]interface{} {
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
			item[k] = c.tryParseDate(vt)
		}
	}

	return item
}

func (c *CommonInserter) replaceMultiRefs(src string, reg *regexp.Regexp) string {
	ms := reg.FindAllString(src, -1)
	if len(ms) == 0 {
		return ""
	}
	for _, m := range ms {
		refID := m[2 : len(m)-1]
		rv, ok := c.refIDs[refID]
		if !ok {
			log.Printf("%s was not found", refID)
			return ""
		}
		src = strings.Replace(src, m, rv, 1)
	}
	return src
}

func (c *CommonInserter) setRefs(item map[string]interface{}) map[string]interface{} {
	reg := regexp.MustCompile(`\#\{.*?\}`)
	for k, v := range item {
		switch vt := v.(type) {
		case []interface{}:
			new := make([]string, len(vt))
			isStr := true
			for i, vtv := range vt {
				vStr, ok := vtv.(string)
				if !ok {
					isStr = false
					break
				}
				if strings.HasPrefix(vStr, "$") && !reg.MatchString(vStr) {
					refID := strings.TrimPrefix(vStr, "$")
					rv, ok := c.refIDs[refID]
					if !ok {
						log.Printf("%s was not found", refID)
					} else {
						new[i] = rv
					}
					continue
				}
				n := c.replaceMultiRefs(vStr, reg)
				if n != "" {
					new[i] = n
					continue
				}
				new[i] = vStr
			}
			if isStr {
				item[k] = new
			}
		case map[string]interface{}:
			for vtk, vtv := range vt {
				if strings.HasPrefix(vtk, "$") && !reg.MatchString(vtk) {
					refID := strings.TrimPrefix(vtk, "$")
					rk, ok := c.refIDs[refID]
					if !ok {
						log.Printf("%s was not found", refID)
					} else {
						vt[rk] = vtv
						delete(vt, vtk)
					}
				}
				nk := c.replaceMultiRefs(vtk, reg)
				if nk != "" {
					vt[nk] = vtv
					delete(vt, vtk)
				}
			}
		case string:
			if strings.HasPrefix(vt, "$") && !reg.MatchString(vt) {
				refID := strings.TrimPrefix(vt, "$")
				rv, ok := c.refIDs[refID]
				if !ok {
					log.Printf("%s was not found", refID)
				} else {
					item[k] = rv
				}
			}
			n := c.replaceMultiRefs(vt, reg)
			if n != "" {
				item[k] = n
			}
		}
	}

	return item
}
