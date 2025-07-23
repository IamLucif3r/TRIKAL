package models

import "time"

type Article struct {
	Title       string
	Link        string
	SummaryRaw  string
	Source      string
	PublishedAt *time.Time
}
