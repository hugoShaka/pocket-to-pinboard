package common

import "time"

type Article struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Tags        []string  `json:"tags"`
	Date        time.Time `json:"date"`
	Favourite   bool      `json:"favourite"`
}

type Articles []Article
