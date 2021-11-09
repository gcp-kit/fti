package inserter

type JsonModel struct {
	Version string      `json:"version"`
	Ref     string      `json:"ref"`
	Payload interface{} `json:"payload"`
}
