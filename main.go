package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v5"
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
		if (!isAuthorized(r)) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

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

func isAuthorized(r *http.Request) bool {
	if r.Host == os.Getenv("SAFE_HOST") || os.Getenv("SAFE_HOST") == "" || os.Getenv("HOST_SECRET") == "" {
		return true
	}

	if (len(r.URL.Query()["token"]) == 0) {
		return false
	}

	requestToken := r.URL.Query()["token"][0]
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("HOST_SECRET")), nil
	})

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if pathClaim, ok := claims["path"].(string); ok {
			return r.URL.Path == pathClaim
		}
	}

	return false
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
