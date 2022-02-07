package inserter

type CollectionName string

// Collection - JSONやYAMLのモデル
type Collection struct {
	Version string      `json:"version" yaml:"version"`
	Items   []Document `json:"items" yaml:"version"`
}

// Document - Modelが持つアイテム
type Document struct {
	Ref            string                 `json:"ref" yaml:"ref"`
	Payload        map[string]interface{} `json:"payload" yaml:"payload"`
	SubCollections map[CollectionName][]Document `json:"sub_collections"`
}
