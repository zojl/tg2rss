package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/zojl/tg2rss/fetcher"
	"github.com/zojl/tg2rss/media"
	"github.com/zojl/tg2rss/parser"
	"github.com/zojl/tg2rss/rss"
	"github.com/zojl/tg2rss/validator"
)

const urlPrefix = "https://t.me/"

func main() {
	godotenv.Load()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if len(path) == 0 {
			w.Write([]byte(""))
			return
		}

		if validator.ValidateChannel(path) {
			handleChannel(w, path)
			return
		}

		if validator.ValidateMedia(path) {
			handleMedia(w, r, path)
			return
		}

		http.Error(w, "Invalid url", 404)
	})

	port := os.Getenv("LISTEN_PORT")
	log.Printf("Starting server at %s...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func handleChannel(w http.ResponseWriter, channelName string) {
	html, err := fetcher.FetchHTML(fmt.Sprintf("%ss/%s", urlPrefix, channelName))
	if err != nil {
		http.Error(w, "Channel not found", 404)
		return
	}
	data, err := parser.ParseHTML(html)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}

	data.Link = urlPrefix + channelName
	rssFeed, _ := rss.GenerateFeed(data)
	w.Write([]byte(rssFeed))
}

func handleMedia(w http.ResponseWriter, r *http.Request, path string) {
	postPath := media.GetPostPath(path)
	html, err := fetcher.FetchHTML(urlPrefix + postPath)

	if (err != nil) {
		http.Error(w, err.Error(), 404)
		return
	}

	mediaPath, err := parser.ParseMedia(html)
	if (err != nil) {
		http.Error(w, err.Error(), 404)
		return
	}

	http.Redirect(w, r, mediaPath, http.StatusSeeOther)
}
