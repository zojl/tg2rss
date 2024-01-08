package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/zojl/tg2rss/fetcher"
	"github.com/zojl/tg2rss/parser"
	"github.com/zojl/tg2rss/rss"
	"github.com/zojl/tg2rss/validator"
)

const urlPrefix = "https://t.me/s/"

func main() {
	godotenv.Load()

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

	port := os.Getenv("LISTEN_PORT")
	log.Printf("Starting server at %s...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
