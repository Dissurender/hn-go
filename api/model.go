package api

// Base contains possible fields for any response
// non-required fields are optional and may not exist
type Base struct {
	ID          int    `json:"id" binding:"required"`
	Type        string `json:"type" binding:"required"`
	By          string `json:"by"`
	Time        int    `json:"time"`
	Kids        []int  `json:"kids,omitempty"`
	Dead        bool   `json:"-"`
	Deleted     bool   `json:"-"`
	Descendants int    `json:"descendants"`
	Score       int    `json:"score"`
	Title       string `json:"title,omitempty"`
	Parent      int    `json:"parent,omitempty"`
	URL         string `json:"url,omitempty"`
	Text        string `json:"text,omitempty"`
}

type BaseWithKids struct {
	ID          int           `json:"id" binding:"required"`
	Type        string        `json:"type" binding:"required"`
	By          string        `json:"by"`
	Time        int           `json:"time"`
	Kids        []interface{} `json:"kids,omitempty"`
	Dead        bool          `json:"-"`
	Deleted     bool          `json:"-"`
	Descendants int           `json:"descendants"`
	Score       int           `json:"score"`
	Title       string        `json:"title,omitempty"`
	Parent      int           `json:"parent,omitempty"`
	URL         string        `json:"url,omitempty"`
	Text        string        `json:"text,omitempty"`
}

type Comment struct {
	ID          int           `json:"id" binding:"required"`
	Type        string        `json:"type" binding:"required"`
	By          string        `json:"by"`
	Time        int           `json:"time"`
	Kids        []interface{} `json:"kids,omitempty"`
	Dead        bool          `json:"-"`
	Deleted     bool          `json:"-"`
	Descendants int           `json:"descendants"`
	Score       int           `json:"score"`
	Parent      int           `json:"parent,omitempty"`
	Text        string        `json:"text,omitempty"`
}
