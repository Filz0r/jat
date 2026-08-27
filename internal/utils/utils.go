package utils

import "os"

type CompanyCounts struct {
	UserCount  int64
	TotalCount int64
}

func IsDevMode() bool {
	return os.Getenv("JAT_DEV") != ""
}
