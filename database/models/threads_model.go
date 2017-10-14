package models

import "time"

type Threads struct {
	TID         int64     `json:"tid"`
	Author      int64     `json:"author"`
	Created     time.Time `json:"created"`
	Forum       int64     `json:"forum"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Votes       int32     `json:"votes"`
}
