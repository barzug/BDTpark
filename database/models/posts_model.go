package models

import "time"

type Posts struct {
	PID      int64     `json:"pid"`
	Author   int64     `json:"author"`
	Created  time.Time `json:"created"`
	Forum    int64     `json:"forum"`
	IsEdited bool      `json:"isedited"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parents"`
	Threads  int64     `json:"threads"`
}
