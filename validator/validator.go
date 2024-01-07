package validator

import (
	"regexp"
)

func Validate(channelName string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_]{1,32}$`, channelName)
	return match
}
