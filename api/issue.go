package api

type Issue struct {
	Id                 string
	Status             string
	Severity           string
	Description        string
	CreatedAt          string
	UpdatedAt          string
	ResolvedAt         string
	DueAt              string
	StatusChangedAt    string
	RejectionExpiredAt string
}
