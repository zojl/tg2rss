package parser

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zojl/tg2rss/rss"
	"github.com/zojl/tg2rss/parser/pyBridge"
	"github.com/PuerkitoBio/goquery"
)

const expectedItemsCount = 15
const timeLayout = "2006-01-02T15:04:05Z07:00"

func ParseHTML(html string) (rss.Channel, error) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	channel := rss.Channel{
		Items: make([]rss.Item, 0, expectedItemsCount),
	}

	titleSelection := doc.Find(".tgme_channel_info_header_title").First()
	channel.Title = titleSelection.Text()

	if len(channel.Title) == 0 {
		return rss.Channel{}, errors.New("Not a public channel")
	}

	descriptionSelection := doc.Find(".tgme_channel_info_description").First()
	replaceLineBreaks(descriptionSelection)
	channel.Description = descriptionSelection.Text()

	doc.Find(".tgme_widget_message").Each(func(i int, post *goquery.Selection) {
		newItem := getItem(post)
		channel.Items = append(channel.Items, newItem)
	})

	return channel, nil
}

func ParseMedia(html string) (string, error) {
	post, _ := goquery.NewDocumentFromReader(strings.NewReader(html))

	photos := post.Find("a.tgme_widget_message_photo_wrap")
	if photos.Length() > 0 {
		photo := photos.First()
		styleRaw, _ := photo.Attr("style")
		return getBackgroundImage(styleRaw), nil
	}

	videos := post.Find("a.tgme_widget_message_video_player")
	if videos.Length() > 0 && !videos.First().HasClass("not_supported") {
		video := videos.First()
		videoWrap := video.Find(".tgme_widget_message_video_wrap").First()
		videoLink, _ := videoWrap.Find("video").First().Attr("src")
		return videoLink, nil
	}

	videoPreviews := post.Find(".tgme_widget_message_video_thumb")
	if videoPreviews.Length() > 0 {
		videoPreview := videoPreviews.First()
		styleRaw, _ := videoPreview.Attr("style")
		return getBackgroundImage(styleRaw), nil
	}

	return "", errors.New("post not found or has no media")
}

func replaceLineBreaks(selection *goquery.Selection) {
	selection.Find("br").Each(func(j int, replaceLineBreakselection *goquery.Selection) {
		_ = replaceLineBreakselection.ReplaceWithHtml("\n")
	})
}

func getItem(post *goquery.Selection) rss.Item {
	maxTitleLength, err := strconv.Atoi(os.Getenv("MAX_TITLE_LENGTH"))
	if maxTitleLength < 3 || err != nil {
		maxTitleLength = 3
	}

	item := rss.Item{}
	messageDate := post.Find(".tgme_widget_message_date").First()
	item.Link, _ = messageDate.Attr("href")

	media := post.Find("a.tgme_widget_message_photo_wrap, a.tgme_widget_message_video_player")
	hasMedia := media.Length() > 0

	if hasMedia {
		post.Find("a.tgme_widget_message_photo_wrap").Each(func(j int, photoContent *goquery.Selection) {
			item.Content = item.Content + makePhotoBlock(photoContent)
		})

		post.Find("a.tgme_widget_message_video_player").Each(func(j int, previewContent *goquery.Selection) {
			if (previewContent.HasClass("not_supported")) {
				item.Content = item.Content + makeUnsupportedVideoPreview(previewContent)
				return
			}

			videoPlayer := makeVideoPlayer(previewContent)
			item.Content = item.Content + videoPlayer + "<br>"
		})
	}

	post.Find(".tgme_widget_message_text").Each(func(j int, text *goquery.Selection) {
		if text.Find(".tgme_widget_message_text").Length() > 0 {
			return
		}

		textHtml, _ := text.Html()
		item.Content = item.Content + textHtml
		replaceLineBreaks(text)
		item.Description = item.Description + text.Text()
	})

	post.Find(".tgme_widget_message_poll").Each(func(j int, poll *goquery.Selection) {
		item.Description = item.Description + "📊 Poll: "
		poll.Find(".tgme_widget_message_poll_question, .tgme_widget_message_poll_option").Each(func(k int, pollLine *goquery.Selection) {
			item.Description = item.Description + pollLine.Text() + "<br>"
		})
	})

	unsupported := post.Find(".message_media_not_supported_label")
	if ((!hasMedia && item.Description == "") || (unsupported.Length() != 0 && len(item.Content) == 0)) {
		if (os.Getenv("PYROGRAM_BRIDGE_HOST") != "") {
			postId, _ := getPostIdentifier(item.Link)
			pyBridgeItem, _ := pyBridge.GetPost(postId)
			if err != nil {
				log.Println(err)
			} else {
				return pyBridgeItem
			}
		}
		item.Content = fmt.Sprintf("Unsupported post, <a href='%s'>view in Telegram</a>", item.Link)
		item.Description = "Unsupported post: " + item.Link
	}

	postTime, _ := messageDate.Find("time").First().Attr("datetime")
	item.Created, _ = time.Parse(timeLayout, postTime)

	title := item.Description
	if hasMedia {
		title = "🖼️" + title
	}
	if len([]rune(title)) > maxTitleLength {
		item.Title = string([]rune(title)[:maxTitleLength-3]) + "..."
	} else {
		item.Title = title
	}

	return item
}

func makeUnsupportedVideoPreview(previewContent *goquery.Selection) string {
	linkHref, _ := previewContent.Attr("href")
	outImageSrc := ""
	if os.Getenv("PROXY_MEDIA") == "true" {
		identifier, _ := getPostIdentifier(linkHref)
		outImagePath := fmt.Sprintf("/media%s.%s.jpg", identifier)
		outImageSrc = os.Getenv("MEDIA_HOST") + outImagePath

		if (os.Getenv("HOST_SECRET") != "") {
			token, err := generateJWT(outImagePath)
			if (err != nil) {
				log.Println(err)
			}
			outImageSrc = outImageSrc + "?token=" + token
		}
	} else {
		styleRaw, _ := previewContent.Find("i.tgme_widget_message_video_thumb").First().Attr("style")
		outImageSrc = getBackgroundImage(styleRaw)
	}
	return fmt.Sprintf(
		"<a href=\"%s\" style=\"filter: blur(15px)\"><img src=\"%s\" alt=\"video preview\"></a><br>",
		linkHref,
		outImageSrc,
	)
}

func makePhotoBlock(photoContent *goquery.Selection) string {
	styleRaw, _ := photoContent.Attr("style")
	linkHref, _ := photoContent.Attr("href")
	backgroundImageSrc := getBackgroundImage(styleRaw)
	outImageSrc := backgroundImageSrc
	if os.Getenv("PROXY_MEDIA") == "true" {
		extension := extractExtension(backgroundImageSrc)
		identifier, _ := getPostIdentifier(linkHref)
		outImagePath := fmt.Sprintf("/media%s.%s", identifier, extension)
		outImageSrc = os.Getenv("MEDIA_HOST") + outImagePath

		if (os.Getenv("HOST_SECRET") != "") {
			token, err := generateJWT(outImagePath)
			if (err != nil) {
				log.Println(err)
			}
			outImageSrc = outImageSrc + "?token=" + token
		}
	}

	return fmt.Sprintf(
		"<a href='%s'><img src='%s' alt='post image'></a><br>",
		linkHref,
		outImageSrc,
	)
}

func makeVideoPlayer(previewContent *goquery.Selection) string {
	videoWrap := previewContent.Find(".tgme_widget_message_video_wrap").First()
	videoWrapStyle, _ := videoWrap.Attr("style")
	var videoLink string
	if os.Getenv("PROXY_MEDIA") == "true" {
		linkHref, _ := previewContent.Attr("href")
		identifier, _ := getPostIdentifier(linkHref)
		videoPath := fmt.Sprintf("/media%s.%s", identifier, "mp4")
		videoLink = os.Getenv("MEDIA_HOST") + videoPath

		if (os.Getenv("HOST_SECRET") != "") {
			token, err := generateJWT(videoPath)
			if (err != nil) {
				log.Println(err)
			}
			videoLink = videoLink + "?token=" + token
		}
	} else {
		videoLink, _ = videoWrap.Find("video").First().Attr("src")
	}
	videoPlayer := fmt.Sprintf("<div style=\"%s\"><video src=\"%s\" width=\"100%%\" height=\"100%%\"></video></div>", videoWrapStyle, videoLink)
	return videoPlayer
}

func getBackgroundImage(inline string) string {
	pattern := `background-image:url\('(.*)'\)`
	r := regexp.MustCompile(pattern)
	matches := r.FindStringSubmatch(inline)

	if len(matches) > 1 {
	    return matches[1]
	}

	return ""
}

func getPostIdentifier(postUrl string) (string, error) {
	parsedURL, err := url.Parse(postUrl)
	if err != nil {
		return "", err
	}
	path := parsedURL.Path
	return path, nil
}

func extractExtension(url string) string {
	parts := strings.Split(url, "/")
	file := parts[len(parts)-1]
	fileParts := strings.Split(file, ".")
	extension := fileParts[len(fileParts)-1]
	return extension
}

func generateJWT(userPath string) (string, error) {
	claims := jwt.MapClaims{
		"path": userPath,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("HOST_SECRET")))
}