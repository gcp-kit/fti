// Package inserter - Firestore にダミーデータを追加するためのパッケージ
package inserter

// JSONModel - JSONのモデル
type JSONModel struct {
	Version string          `json:"version"`
	Items   []JSONModelItem `json:"items"`
}

// JSONModelItem - JsonModel のアイテム
type JSONModelItem struct {
	Ref            string                     `json:"ref"`
	Payload        map[string]interface{}     `json:"payload"`
	SubCollections map[string][]JSONModelItem `json:"sub_collections"`
}
