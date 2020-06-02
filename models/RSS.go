package models

type RSSData struct {
	UUID    string `gorm:"primary_key"`
	Title   string
	Source  string
	PubDate string
}
