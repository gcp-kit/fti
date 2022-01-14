package inserter

type Model struct {
	Version string      `json:"version" yaml:"version"`
	Items   []ModelItem `json:"items" yaml:"version"`
}

type ModelItem struct {
	Ref     string                 `json:"ref" yaml:"ref"`
	Payload map[string]interface{} `json:"payload" yaml:"payload"`
}
