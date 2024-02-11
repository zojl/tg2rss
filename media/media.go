package media

import (
	"strings"
)

func GetPostPath(path string) (string) {
	path = strings.TrimPrefix(path, "media/")
	lastDotIndex := strings.LastIndex(path, ".")
	if lastDotIndex == -1 {
		// Если точка не найдена, возвращаем исходный путь
		return path
	}
	identifier := path[:lastDotIndex]
	return identifier
}
