package models

type Votes struct {
	Voice  int32 `json:"voice"`
	User   int64 `json:"user"`
	Thread int64 `json:"thread"`
}
