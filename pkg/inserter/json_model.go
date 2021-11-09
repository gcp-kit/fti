package inserter

type JsonModel struct {
	Version string          `json:"version"`
	Items   []JsonModelItem `json:"items"`
}

type JsonModelItem struct {
	Ref     string                 `json:"ref"`
	Payload map[string]interface{} `json:"payload"`
}
