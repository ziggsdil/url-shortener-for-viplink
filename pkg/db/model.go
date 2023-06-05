package db

type Link struct {
	ShortSuffix string `json:"short_suffix"`
	Link        string `json:"link"`
	SecretKey   string `json:"secret_key"`
	Clicks      int    `json:"clicks"`
}
