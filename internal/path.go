package internal

import (
	"os/user"
	"path/filepath"
	"strings"
)

func ResolvePath(filePath string) (string, error) {
	userInfo, err := user.Current()
	if err != nil {
		return "", err
	}
	switch {
	case filePath == "~":
		return userInfo.HomeDir, nil
	case strings.HasPrefix(filePath, "~/"):
		return filepath.Join(userInfo.HomeDir, filePath[2:]), nil
	default:
		return filepath.Abs(filePath)
	}
}
