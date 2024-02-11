package validator

import (
	"regexp"
)

func ValidateChannel(channelName string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]{1,32}$`, channelName)
	return match
}

func ValidateMedia(mediaPath string) bool {
	match, _ := regexp.MatchString(`^media/[a-zA-Z0-9_]{1,32}/[0-9]*\.[a-zA-Z0-9]{3,6}$`, mediaPath)
	return match
}
