package models

//改改改！！！
import (
	"errors"
	_ "github.com/astaxie/beego/validation"
	"strconv"
	"time"
)

//mysql
//结构体首字母要大写，小写的成员转化为json数据时会直接被忽略

type Card struct {
	CardId      string    `orm:"pk;size(16)" valid:"Required;Length(16)"` //CardId 编码暂时按照上学期的编码
	UserId      string    `orm:"size(13)" valid:"Required;Length(13)"`    //UserId 必须是由用户给出的，因为CardId中不包含UserId
	CardType    string    `valid:"Required"`                              //卡的类型
	Enterprise  string    `valid:"Required"`
	State       string    `valid:"Required"`
	City        string    `valid:"Required"`
	Money       int       `orm:"default(0)"`
	Score       int    `orm:"null"`
	CouponsNum  int    `orm:"null"` //每一种种类的数量，数量与数量之间用空格隔开
	Coupons     string    `orm:"null"` //描述优惠的方法
	ExpireTime  time.Time `valid:"Required"`
	DelTime     time.Time `orm:"null"`
	CardOrder   int       `valid:"Required"` //该商家合作以来发布的第N条卡片
	FactoryNum  int       `valid:"Required"`
	BatchNum    int       `valid:"Required"`
	SerialNum   int       `valid:"Required"`
}

type CardParseStruct struct {
	EnterpriseMap map[string]string `orm:"-"`
	KindMap       map[string]string `orm:"-"`
	StateMap      map[string]string `orm:"-"`
	CityMap       map[string]string `orm:"-"`
}

type Enterprise struct {
	Id          string `orm:"unique"`
	IsLocal     string `orm:"column(is_local)"`
	Type        string
	RegisterNum string `orm:"column(register_num)"`
	Name        string `orm:"pk"`
	HelpMsg     string `orm:"column(help_msg)"`
	Website     string
}

type EnterpriseParseStruct struct {
	IsLocalMap map[string]string
	TypeMap    map[string]string
}

type User struct {
	Id       string `orm:"pk;size(13)" valid:"Required"` 
	Tel      string `orm:"null"` 
	Mail     string `orm:"null"`
	Password string `valid:"Required"`
	LoginMonth string `valid:"max(2)"`
	LoginYear  string `valid:"max(4)"`
	LoginNum int `valid:"MaxSize(6)"`
}

// type DelCard struct {
// 	CardId  string `orm:"pk;column(card_id)"`
// 	UserId  string `orm:"column(user_id)"`
// 	Remark  string
// 	DelTime time.Time `orm:"column(del_time)"`
// }

var CardParseMaps = CardParseStruct{
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

var EnterpriseParseMaps = EnterpriseParseStruct{
	map[string]string{
		"1": "True",
		"2": "False",
	},
	map[string]string{
		"1": "Bank",
		"2": "Supermarket",
		"3": "Store",
	},
}

//将card结构中的Id解析出对应的含义赋值给card的其他导出属性
func (card *Card) CardParse() error {
	if len(card.CardId) != 16 {
		return errors.New("INVALID LENGTH CARD ID")
	}
	var ok bool
	var err error
	card.Enterprise, ok = CardParseMaps.EnterpriseMap[card.CardId[0:3]]
	card.CardType, ok = CardParseMaps.KindMap[card.CardId[3:4]]
	card.CardOrder, err = strconv.Atoi(card.CardId[4:6])
	card.State, ok = CardParseMaps.StateMap[card.CardId[6:7]]
	card.City, ok = CardParseMaps.CityMap[card.CardId[6:10]]
	card.FactoryNum, err = strconv.Atoi(card.CardId[10:12])
	card.BatchNum, err = strconv.Atoi(card.CardId[12:13])
	card.SerialNum, err = strconv.Atoi(card.CardId[13:])
	if !ok && err != nil {
		return errors.New("INVALID CARD ID")
	}
	return nil
}

//ZYJ 解析生成用户ID
func (user *User) UserParse() error {
	if len(user.Id) != 16 {
		return errors.New("INVALID LENGTH USER ID")
	}
	var err error
	user.LoginYear = user.Id[0:4]
	user.LoginMonth = user.Id[4:6]
	user.LoginNum,err = strconv.Atoi(user.Id[6:])
	if err!=nil{
		return errors.New("INVALID USER")
	}
	return nil
}

func (enterprise *Enterprise) EnterpriseParse() error {
	if len(enterprise.Id) != 5 {
		return errors.New("INVALID LENGTH ENTERPRISE ID")
	}
	var flag bool
	enterprise.IsLocal, flag = EnterpriseParseMaps.IsLocalMap[enterprise.Id[0:1]]
	enterprise.Type, flag = EnterpriseParseMaps.TypeMap[enterprise.Id[1:2]]
	enterprise.RegisterNum = enterprise.Id[2:]
	if !flag {
		return errors.New("INVALID CONTENT ENTERPRISE ID")
	}
	return nil
}

// func (delCard *DelCard) GetTime() {
// 	delCard.DelTime = time.Now();
// }
