package media

import (
	"fmt"
	"strings"
)

func GetPostPath(path string) (string) {
	suffix := "?embed=1&mode=tme&single=1"
	path = strings.TrimPrefix(path, "media/")
	lastDotIndex := strings.LastIndex(path, ".")
	if lastDotIndex == -1 {
		return fmt.Sprintf("%s%s", path, suffix)
	}
	identifier := path[:lastDotIndex]
	return fmt.Sprintf("%s%s", identifier, suffix);
}
