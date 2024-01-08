package main

import (
	"log"
	"net/http"
	"github.com/zojl/tg2rss/fetcher"
	"github.com/zojl/tg2rss/parser"
	"github.com/zojl/tg2rss/rss"
	"github.com/zojl/tg2rss/validator"
)

const urlPrefix = "https://t.me/s/"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		channelName := r.URL.Path[1:]
		if (len(channelName) == 0) {
			w.Write([]byte(""))
			return
		}

		if !validator.Validate(channelName) {
			http.Error(w, "Invalid channel name", 404)
			return
		}

		html, err := fetcher.FetchHTML(urlPrefix + channelName)
		if err != nil {
			http.Error(w, "Channel not found", 404)
			return
		}

		data, err := parser.ParseHTML(html)
		if err != nil {
			http.Error(w, err.Error(), 404)
			return
		}
		data.Link = urlPrefix + channelName;

		rssFeed, _ := rss.GenerateFeed(data)
		w.Write([]byte(rssFeed))
	})

	log.Println("Starting server...")
	http.ListenAndServe(":80", nil)
}
