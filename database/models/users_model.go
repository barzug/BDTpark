package models

type Users struct {
	UID      int64  `json:"uid"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
}
