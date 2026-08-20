package utils

import "os"

func IsDevMode() bool {
	return os.Getenv("JAT_DEV") != ""
}
