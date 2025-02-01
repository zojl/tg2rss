package pyBridge

import (
    "encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/zojl/tg2rss/fetcher"
	"github.com/zojl/tg2rss/rss"
)

type Item struct {
    Id           string `json:"id"`
    Date         string `json:"date"`
    Html         string `json:"html"`
    Title        string `json:"title"`
    Views        string `json:"views"`
}

const timeLayout = "2024-05-14T21:52:34"

func GetPost(postId string) (rss.Item, error) {
	item := rss.Item{}
	host := os.Getenv("PYROGRAM_BRIDGE_HOST")
	url := fmt.Sprintf("%s/json%s", host, postId)
	post, err := fetcher.FetchStream(url)
	if err != nil {
		return item, err
	}

	pybItem := Item{}
	err = json.Unmarshal(post, &pybItem)
	if err != nil {
		return item, err
	}

	item.Title = pybItem.Title
	item.Link = fmt.Sprintf("https://t.me%s", postId)
	item.Description = pybItem.Html
	item.Content = pybItem.Html
	item.Created, _ = time.Parse(timeLayout, pybItem.Date)

	return item, nil
}