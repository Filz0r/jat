package database

type ApplicationStatusKind string

const (
	Applied     ApplicationStatusKind = "applied"
	Rejected                          = "rejected"
	Ghosted                           = "ghosted"
	Interviewed                       = "interviewed"
	Irrelevant                        = "irrelevant"
	Accepted                          = "accepted"
)

func (k ApplicationStatusKind) String() string {
	return string(k)
}

func (k ApplicationStatusKind) Valid() bool {
	switch k {
	case Applied, Rejected, Ghosted, Interviewed, Irrelevant, Accepted:
		return true
	default:
		return false
	}
}
