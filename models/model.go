package models

import "time"

type Card struct {
	Id         string `orm:"pk"`
	UserId     string `orm:"column(user_id);rel(fk)"`
	Kind       string `orm:"column(type)"`
	Style      string
	Remark     string
	EName      string `orm:"column(e_name)"`
	State      string
	City       string
	FactoryNum string `orm:"column(factory_num)"` //印刷厂编号
	BatchNum   string `orm:"column(batch_num)"`   //印刷批次
	SerialNum  string `orm:"column(serial_num)"`  //同批次的卡片编号
}

type User struct {
	Id       string `orm:"pk"`
	Tel      string
	Mail     string
	Password string
}

type Enterprise struct {
	Id      string `orm:"pk"`
	Name    string
	HelpMsg string `orm:"column(help_msg)"`
	Website string
}

type DelCard struct {
	CardId  string `orm:"pk;column(card_id)"`
	UserId  string `orm:"column(user_id);rel(fk)"`
	Remark  string
	DelTime string `orm:"column(del_type)"`
}

func (this *DelCard) setDelTime() {
	t := time.Now().Format("2006-01-02 15:04:05")
	this.DelTime = t
}
