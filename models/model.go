package models

import "time"

type Card struct {
	Id     string `orm:"pk"`
	UserId string `orm:"column(user_id)"`
	Kind   string `orm:"column(type)"`
	Remark string
}

type User struct {
	Id       string `orm:"pk"`
	Tel      string
	Mail     string
	Password string
}

type DelCard struct {
	Id      string `orm:"pk"`
	UserId  string
	Kind    string `orm:"column(type)"`
	Remark  string
	DelTime time.Time `orm:"column(del_type)"`
}
