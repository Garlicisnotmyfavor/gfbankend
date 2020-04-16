package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/astaxie/beego/validation"
	"strconv"
	"time"
)

//mysql
//结构体首字母要大写，小写的成员转化为json数据时会直接被忽略

type Card struct {
	CardId string `orm:"pk;size(16)" valid:"Required;Length(16)"` //CardId 编码暂时按照上学期的编码
	// UserId      *User     `orm:"rel(fk)"`
	UserId     string    `orm:"size(13)"`   //UserId是依据时间生成的，因为CardId中不包含UserId
	CardType   string    `valid:"Required"` //卡的类型
	Enterprise string    `valid:"Required"`
	State      string    `valid:"Required"`
	City       string    `valid:"Required"`
	Money      int       `orm:"default(0)"`
	Score      int       `orm:"null"`
	CouponsNum int       `orm:"null"` //每一种种类的数量，数量与数量之间用空格隔开
	Coupons    string    `orm:"null"` //描述优惠的方法
	ExpireTime time.Time `valid:"Required"`
	DelTime    time.Time `orm:"null"`
	CardOrder  int       `valid:"Required"` //该商家合作以来发布的第N条卡片
	FactoryNum int       `valid:"Required"`
	BatchNum   int       `valid:"Required"`
	SerialNum  int       `valid:"Required"`
}

type CardDemo struct {
	ID         int    `orm:"pk;auto"`
	CardType   string `valid:"Required"`
	Enterprise string `valid:"Required"`
	State      string `valid:"Required"`
	City       string `valid:"Required"`
	Coupons    string `orm:"null"`
}

type Enterprise struct {
	Id          string `orm:"unique"`
	Password    string
	Addr        string `orm:"column(addr)"`
	IsLocal     bool   `orm:"column(is_local)"`
	Type        string
	RegisterNum int64  `orm:"column(register_num)"`
	Name        string `orm:"pk"`
	HelpMsg     string `orm:"column(help_msg)"`
	Website     string
	LicenseId   string
}

type Manager struct {
	Name       string //管理员名称
	ID         string `orm:"pk;column(id)"` //身份证号(保证唯一）
	Enterprise string //与企业关联, n..1关系
	Phone      string //手机号，登录用
	Password   string
}

type User struct {
	Id         string `orm:"pk;size(13)" valid:"Required"`
	Tel        string `orm:"null"`
	Mail       string `orm:"null"`
	Password   string `valid:"Required"`
	LoginMonth string `valid:"max(2)" `                      //注册月份
	LoginYear  string `valid:"max(4)" `                      //注册年份
	LoginNum   int    `valid:"MaxSize(6)" orm:"default(1)" ` //该月份所注册的第几个用户
}

type Count struct {
	Time string `valid:"max(7)" orm:"pk"`
	Num  int    `orm:"default(1)"`
}

type EnterpriseCount struct {
	Flag int `orm:"pk;default(1)"`
	Num  int
}

type CardLog struct {
	CardId  string
	Date    time.Time
	Operate string //操作描述
}

//type DelCard struct {
//	CardId  string `orm:"pk;column(card_id)"`
//	UserId  string `orm:"column(user_id)"`
//	Remark  string
//	DelTime time.Time `orm:"column(del_time)"`
//}

/**
*	Below are some maps for parse
 */
type EnterpriseParseStruct struct {
	IsLocalMap map[string]string
	TypeMap    map[string]string
}

type CardParseStruct struct {
	EnterpriseMap map[string]string `orm:"-"`
	KindMap       map[string]string `orm:"-"`
	StateMap      map[string]string `orm:"-"`
	CityMap       map[string]string `orm:"-"`
}

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

//根据confluence
//zyj
//var UserParse
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
func (user *User) UserParse() {
	o := orm.NewOrm()
	curTime := time.Now().String()[:7]
	var item Count
	user.Id = curTime[0:4] + curTime[5:7]
	user.LoginYear = curTime[0:4]
	user.LoginMonth = curTime[5:7]
	item.Time = user.LoginYear + "-" + user.LoginMonth
	if err := o.Read(&item); err != nil {
		item.Num = 1
		user.LoginNum = 1
		user.Id += fmt.Sprintf("%07d", user.LoginNum)
		o.Insert(&item)
		return
	}
	item.Num += 1
	fmt.Println(item)
	user.LoginNum = item.Num
	user.Id += fmt.Sprintf("%07d", user.LoginNum)
	fmt.Println(user)
	o.Update(&item)
	return
}

func (enterprise *Enterprise) EnterpriseParse() error {
	o := orm.NewOrm()
	var item EnterpriseCount
	item.Flag = 1 //flag这个没有意义，只是用于充当主键，取出数据库中的数据
	if err := o.Read(&item); err != nil {
		return errors.New("fail to get registerNum")
	}
	item.Num += 1
	enterprise.RegisterNum = int64(item.Num)
	o.Update(&item)
	if enterprise.IsLocal == true {
		enterprise.Id = "1"
	} else {
		enterprise.Id = "2"
	}
	if enterprise.Type == "bank" {
		enterprise.Id += "1"
	} else if enterprise.Type == "supermarket" {
		enterprise.Id += "2"
	} else if enterprise.Type == "store" {
		enterprise.Id += "3"
	}
	enterprise.Id += strconv.Itoa(int(enterprise.RegisterNum))
	return nil
}

// func (delCard *DelCard) GetTime() {
// 	delCard.DelTime = time.Now();
// }
