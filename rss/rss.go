package rss

import (
	"time"
	"github.com/gorilla/feeds"
)

type Item struct {
	Title       string
	Link        string
	Description string
	Content		string
	Created 	time.Time
}

type Channel struct {
	Title       string
	Link        string
	Description string
	Items       []Item
}

func GenerateFeed(data Channel) (string, error) {
	feed := &feeds.Feed{
		Title: data.Title,
		Link: &feeds.Link{Href: data.Link},
		Description: data.Description,
	}
	
	feed.Items = make([]*feeds.Item, len(data.Items))
	for i, item := range data.Items {
		feed.Items[i] = &feeds.Item{
			Title: item.Title,
			Link: &feeds.Link{Href: item.Link},
			Description: item.Description,
			Content: item.Content,
			Created: item.Created,
		}
	}

	rss, err := feed.ToRss()
	if err != nil {
		return "", err
	}

	return rss, nil
}
