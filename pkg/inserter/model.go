package inserter

// CollectionName - Name of collection
type CollectionName string

// Collection - JSONやYAMLのモデル
type Collection struct {
	Version string     `json:"version" yaml:"version"`
	Items   []Document `json:"items" yaml:"version"`
}

// Document - Modelが持つアイテム
type Document struct {
	ID             string                        `json:"id" yaml:"id"`
	Ref            string                        `json:"ref" yaml:"ref"`
	Payload        map[string]interface{}        `json:"payload" yaml:"payload"`
	SubCollections map[CollectionName][]Document `json:"sub_collections"`
}
