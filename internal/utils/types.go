package utils

import "github.com/filz0r/jat/internal/database"

type SuggestionRecord struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type CompanyCounts struct {
	UserCount  int64
	TotalCount int64
}

type KindCount struct {
	Kind  database.ApplicationStatusKind
	Count int64
}

type ApplicationStatusKindCounts struct {
	Applied     int64 `json:"applied" validate:"required"`
	Ghosted     int64 `json:"ghosted" validate:"required"`
	Rejected    int64 `json:"rejected" validate:"required"`
	Accepted    int64 `json:"accepted" validate:"required"`
	Irrelevant  int64 `json:"irrelevant" validate:"required"`
	Interviewed int64 `json:"interviewed" validate:"required"`
}
