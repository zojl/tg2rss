package parser

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/zojl/tg2rss/rss"
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
		channel.Items = append(channel.Items, getItem(post))
	})

	return channel, nil
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
	
	media := post.Find("a.tgme_widget_message_photo_wrap, a.tgme_widget_message_video_player")
	hasMedia := media.Length() > 0
	
	if hasMedia {
		post.Find("a.tgme_widget_message_photo_wrap").Each(func(j int, photoContent *goquery.Selection) {
			styleRaw, _ := photoContent.Attr("style")
			linkHref, _ := photoContent.Attr("href")
			item.Content = item.Content + fmt.Sprintf(
				"<a href='%s'><img src='%s' alt='post image'></a><br>",
				linkHref,
				getBackgroundImage(styleRaw),
			)
		})
		
		post.Find("a.tgme_widget_message_video_player").Each(func(j int, previewContent *goquery.Selection) {
			if (previewContent.HasClass("not_supported")) {
				styleRaw, _ := previewContent.Find("i.tgme_widget_message_video_thumb").First().Attr("style")
				linkHref, _ := previewContent.Attr("href")
				item.Content = item.Content + fmt.Sprintf(
					"<a href='%s'><img src='%s' alt='video preview'></a><br>",
					linkHref,
					getBackgroundImage(styleRaw),
				)
				return
			}

			videoPlayer, _ := previewContent.Find(".tgme_widget_message_video_wrap").First().Html()
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
	
	messageDate := post.Find(".tgme_widget_message_date").First()
	item.Link, _ = messageDate.Attr("href")
	
	postTime, _ := messageDate.Find("time").First().Attr("datetime")
	item.Created, _ = time.Parse(timeLayout, postTime)
	
	unsupported := post.Find(".message_media_not_supported_label")
	if (unsupported.Length() != 0) && (len(item.Content) == 0) {
		item.Content = fmt.Sprintf("Unsupported post, <a href='%s'>view in Telegram</a>", item.Link)
		item.Description = "Unsopported post: " + item.Link
	}
	
	title := item.Description
	if hasMedia {
		title = "🖼️" + title
	}
	if len([]rune(title)) > maxTitleLength {
		item.Title = string([]rune(title)[:maxTitleLength - 3]) + "..."
	} else {
		item.Title = title
	}

	return item
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