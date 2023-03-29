package api

// Base contains possible fields for any response
// non-required fields are optional and may not exist
type Base struct {
	ID          int    `json:"id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	By          string `json:"by"`
	Time        int    `json:"time"`
	Kids        []Kid  `json:"kids"`
	Dead        bool   `json:"dead"`
	Deleted     bool   `json:"deleted"`
	Descendants int    `json:"descendants"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

type Kid struct {
	Base   Base
	Parent int    `json:"parent"`
	Text   string `json:"text"`
}
