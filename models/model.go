package models

//改改改！！！
import (
	"errors"
	"time"
)

//zyj
type Card struct {
	Id         string `orm:"pk"`
	UserId     string `orm:"column(user_id);"` //rel(fk)
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
	UserId  string `orm:"column(user_id)"`
	Remark  string
	DelTime time.Time `orm:"column(del_time)"`
}

type ParseStruct struct {
	EnterpriseMap map[string]string
	KindMap       map[string]string
	StateMap      map[string]string
	CityMap       map[string]string
}

var ParseMaps = ParseStruct{
	map[string]string{
		"001": "ANZ",
		"002": "Calvin Klein",
		"003": "Starbucks",
		"004": "Subway",
	},
	map[string]string{
		"1": "Recharge",
		"2": "Integrate",
		"3": "Discount",
		"4": "RechargeIntegral",
		"5": "RechargeDiscount",
		"6": "IntegralDiscount",
		"7": "RechargeIntegralDiscount",
	},
	map[string]string{
		"1": "New South Wales",
		"2": "Queensland",
		"3": "South Australia",
		"4": "Tasmania",
		"5": "Victoria",
		"6": "Western Australia",
		"7": "Australia Capital Territory",
		"8": "Northern Territory",
	},
	map[string]string{
		"1001": "Sydney",
		"1002": "Wollongong",
		"1003": "Newcastle",
		"2001": "Brisbane",
		"2002": "Gold Coast",
		"2003": "Caloundra",
		"2004": "Townsville",
		"2005": "Cairns",
		"2006": "Toowoomba",
		"3001": "Adelaide",
		"4001": "Hobart",
		"5001": "Melbourne",
		"5002": "Geelong",
		"6001": "Perth",
		"7001": "Canberra",
		"7002": "Jervis Bay",
		"8001": "Darwin",
	},
}

//将card结构中的Id解析出对应的含义赋值给card的其他导出属性
func (card *Card) CardParse() error {
	if len(card.Id) != 16 {
		return errors.New("INVALID LENGTH CARD ID")
	}
	var ok bool
	card.EName, ok = ParseMaps.EnterpriseMap[card.Id[0:3]]
	card.Kind, ok = ParseMaps.KindMap[card.Id[3:4]]
	card.Style = card.Id[4:6]
	card.State, ok = ParseMaps.StateMap[card.Id[6:7]]
	card.City, ok = ParseMaps.CityMap[card.Id[6:10]]
	card.FactoryNum = card.Id[10:12]
	card.BatchNum = card.Id[12:13]
	card.SerialNum = card.Id[13:16]
	if !ok {
		return errors.New("INVALID CONTENT CARD ID")
	}
	return nil
}

// func (delCard *DelCard) GetTime() {
// 	delCard.DelTime = time.Now();
// }
