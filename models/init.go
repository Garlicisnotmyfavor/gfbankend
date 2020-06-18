package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var Log *logs.BeeLogger = logs.NewLogger(10000) //定义日志Log，因为需要在整个project用到，所以需要定义为全局

func init() {
	//日志暂且设置输出到控制台
	if err := Log.SetLogger("console"); err != nil {
		fmt.Println("fail to set logger!")
		os.Exit(-1)
	}
	Log.SetLevel(logs.LevelDebug)
	Log.EnableFuncCallDepth(true)
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		Log.Error("fail to register driver", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=Asia%%2FShanghai",
		beego.AppConfig.String("mysql::user"),
		beego.AppConfig.String("mysql::password"),
		beego.AppConfig.String("mysql::addr"),
		beego.AppConfig.String("mysql::database"))

	if err := orm.RegisterDataBase("default", "mysql", dsn); err != nil {
		Log.Error("fail to register database", err)
	}
	//下面要修改，根据数据库设计
	//zjn
	orm.RegisterModel(new(Card)) //登记orm
	orm.RegisterModel(new(Enterprise))
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Count))
	orm.RegisterModel(new(Manager))
	orm.RegisterModel(new(EnterpriseCount))
	orm.RegisterModel(new(CardLog))
	orm.RegisterModel(new(CardDemo))

	orm.Debug = true
	_ = orm.RunSyncdb("default", false, true)
	testData()
}

/*
*@function:进行时间加减
*@param {基准时间(time.Time)，需要加减的时间(string)}
*@return {t}time.Time
 */
func CalTime(t time.Time, timeStr string) time.Time {
	timePart, err := time.ParseDuration(timeStr)
	if err != nil {
		fmt.Println(err)
		return t
	}
	return t.Add(timePart)
}

func testData() {
	testUserData1 := User{Id: "2018091620000", Tel: "13925678240", Mail: "123456@qq.com", Password: "123456789", LoginNum: 0, LoginYear: "2020", LoginMonth: "1"}
	testCardData1 := Card{CardId: "123456790123456", UserId: "2018091620000", State: "California", City: "San Jose", CardType: "Integrate", Enterprise: "starbucks", ExpireTime: time.Now(), Coupons: "empty"}
	testUserData2 := User{Id: "2018091620001", Tel: "13665372240", Mail: "654321@qq.com", Password: "123908789", LoginNum: 1, LoginYear: "2020", LoginMonth: "1"}
	testCardData2 := Card{CardId: "123456790000000", UserId: "2018091620001", State: "California", City: "San Jose", CardType: "Discount", Enterprise: "subway", ExpireTime: time.Now(), Coupons: "empty"}
	testUserData3 := User{Id: "2018091620002", Tel: "13778372240", Mail: "666666@qq.com", Password: "123009889", LoginNum: 2, LoginYear: "2020", LoginMonth: "1"}
	testCardData3 := Card{CardId: "123456790333000", UserId: "2018091620002", State: "California", City: "San Jose", CardType: "Integrate", Enterprise: "starbucks", ExpireTime: time.Now(), Coupons: "empty"}
	testUserData4 := User{Id: "2018091620020", Tel: "13778788240", Mail: "78902166@qq.com", Password: "33235323", LoginNum: 20, LoginYear: "2020", LoginMonth: "1"}
	container := []Card{
		{CardId: "123456790333001", State: "A", City: "S", CardType: "Integrate", Enterprise: "starbucks", ExpireTime: time.Now(), StartTime: time.Now(), Coupons: "per 100 score to get a free cup of coffee", Score: 12, UserId: "2018091620000", TypeId: 1},
		{CardId: "123456790333002", State: "B", City: "S", CardType: "Integrate", Enterprise: "starbucks", ExpireTime: time.Now(), StartTime: time.Now(), Coupons: "per 100 score to get a free cup of coffee", Score: 13, UserId: "2018091620001", TypeId: 1},
		{CardId: "123456790333003", State: "C", City: "S", CardType: "Integrate", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "per 100 score to get a free cup of coffee",Score:14, UserId: "2018091620002",TypeId:1},
		{CardId: "123456790333004", State: "D", City: "S", CardType: "Integrate", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "per 100 score to get a free cup of coffee",Score:15, UserId: "2018091620020",TypeId:1}, //积分卡
		{CardId: "123456790333005", State: "E", City: "S", CardType: "Discount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:1, UserId: "empty",TypeId:2,DiscountTimes: 5},
		{CardId: "123456790333006", State: "F", City: "S", CardType: "Discount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:2, UserId: "empty",TypeId:2,DiscountTimes: 5},
		{CardId: "123456790333007", State: "G", City: "S", CardType: "Discount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:2, UserId: "empty",TypeId:2,DiscountTimes: 5},
		{CardId: "123456790333008", State: "H", City: "S", CardType: "Discount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:2, UserId: "empty",TypeId:2,DiscountTimes: 5},
		{CardId: "123456790333009", State: "I", City: "S", CardType: "Recharge", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Money:33, UserId: "empty",TypeId:3},
		{CardId: "123456790333010", State: "J", City: "S", CardType: "Recharge", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Money:34, UserId: "empty",TypeId:3},
		{CardId: "123456790333011", State: "K", City: "S", CardType: "Recharge", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Money:35, UserId: "empty",TypeId:3},
		{CardId: "123456790333012", State: "L", City: "S", CardType: "Recharge", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Money:36, UserId: "empty",TypeId:3},
		{CardId: "123456790333013", State: "M", City: "S", CardType: "RechargeIntegral", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Score:12,Money:36, UserId: "empty",TypeId:4},
		{CardId: "123456790333014", State: "N", City: "S", CardType: "RechargeIntegral", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Score:14,Money:37, UserId: "empty",TypeId:4},
		{CardId: "123456790333015", State: "O", City: "S", CardType: "RechargeIntegral", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Score:15,Money:39, UserId: "empty",TypeId:4},
		{CardId: "123456790333016", State: "P", City: "S", CardType: "RechargeIntegral", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "empty",Score:16,Money:38, UserId: "empty",TypeId:4},
		{CardId: "123456790333017", State: "Q", City: "S", CardType: "RechargeDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:1,Money:36, UserId: "empty",TypeId:5,DiscountTimes: 5},
		{CardId: "123456790333018", State: "R", City: "S", CardType: "RechargeDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:2,Money:37, UserId: "empty",TypeId:5,DiscountTimes: 5},
		{CardId: "123456790333019", State: "S", City: "S", CardType: "RechargeDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:3,Money:38, UserId: "empty",TypeId:5,DiscountTimes: 5},
		{CardId: "123456790333020", State: "T", City: "S", CardType: "RechargeDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",CouponsNum:4,Money:39, UserId: "empty",TypeId:5,DiscountTimes: 5},
		{CardId: "123456790333021", State: "U", City: "S", CardType: "IntegralDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",Score:12,CouponsNum:1, UserId: "empty",TypeId:6,DiscountTimes: 5},
		{CardId: "123456790333022", State: "V", City: "S", CardType: "IntegralDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",Score:12,CouponsNum:1, UserId: "empty",TypeId:6,DiscountTimes: 5},
		{CardId: "123456790333023", State: "W", City: "S", CardType: "IntegralDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",Score:12,CouponsNum:1, UserId: "empty",TypeId:6,DiscountTimes: 5},
		{CardId: "123456790333024", State: "X", City: "S", CardType: "IntegralDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars",Score:12,CouponsNum:1, UserId: "empty",TypeId:6,DiscountTimes: 5},
		{CardId: "123456790333026", State: "Z", City: "S", CardType: "RechargeIntegralDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars;per 100 score to get a free cup of coffee",Score:13,CouponsNum:2, UserId: "empty",TypeId:7,DiscountTimes: 5},
		{CardId: "123456790333027", State: "AB", City: "S", CardType: "RechargeIntegralDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars;per 100 score to get a free cup of coffee",Score:14,CouponsNum:3, UserId: "empty",TypeId:7,DiscountTimes: 5},
		{CardId: "123456790333028", State: "AC", City: "S", CardType: "RechargeIntegralDiscount", Enterprise: "starbucks", ExpireTime: time.Now(),StartTime: time.Now(), Coupons: "90% discount while having 5 stars;per 100 score to get a free cup of coffee",Score:15,CouponsNum:4, UserId: "empty",TypeId:7,DiscountTimes: 5},
	}
	demo := []CardDemo{
		{Id: 1,CardType:"Integrate",Enterprise:"starbucks",State:"empty",City: "empty",Coupons:"per 100 score to get a free cup of coffee",Describe: "go to get a starbucks integrate card!",ExpireTime: CalTime(time.Now(),"720h")},
		{Id: 2,CardType:"Discount",Enterprise:"starbucks",State:"empty",City: "empty",Coupons:"90% discount while having 5 stars",Describe: "go to get a starbucks discount card!",ExpireTime: CalTime(time.Now(),"720h")},
		{Id: 3,CardType:"Recharge",Enterprise:"starbucks",State:"empty",City: "empty",Coupons:"empty",Describe: "go to get a starbucks Recharge card!",ExpireTime: CalTime(time.Now(),"720h")},
		{Id: 4,CardType:"RechargeIntegral",Enterprise:"starbucks",State:"empty",City: "empty",Coupons:"per 100 score to get a free cup of coffee",Describe: "go to get a starbucks RechargeIntegral card!",ExpireTime: CalTime(time.Now(),"720h")},
		{Id: 5,CardType:"RechargeDiscount",Enterprise:"starbucks",State:"empty",City: "empty",Coupons:"90% discount while having 5 stars",Describe: "go to get a starbucks RechargeDiscount card!",ExpireTime: CalTime(time.Now(),"720h")},
		{Id: 6,CardType:"IntegralDiscount",Enterprise:"starbucks",State:"empty",City: "empty",Coupons:"90% discount while having 5 stars;per 100 score to get a free cup of coffee",Describe: "go to get a starbucks IntegralDiscount card!",ExpireTime: CalTime(time.Now(),"720h")},
		{Id: 7,CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",State:"empty",City: "empty",Coupons:"90% discount while having 5 stars;per 100 score to get a free cup of coffee",Describe: "go to get a starbucks RechargeIntegralDiscount card!",ExpireTime: CalTime(time.Now(),"720h")},
	}
	testEnterpriseData1 := Enterprise{Id: "001", Name: "starbucks", Addr: "empty", IsLocal: false, Type: "empty", HelpMsg: "empty", Website: "empty", LicenseId: "empty"}
	testEnterpriseData2 := Enterprise{Id: "002", Name: "subway", Addr: "empty", IsLocal: false, Type: "empty", HelpMsg: "empty", Website: "empty", LicenseId: "empty"}
	enterprise := Enterprise{Id: "13002", Name: "starbucks"}
	ecount := EnterpriseCount{Flag: 1, Num: 0}
	o := orm.NewOrm()
	o.Insert(&testUserData1)
	o.Insert(&testUserData2)
	o.Insert(&testUserData3)
	o.Insert(&testUserData4)
	o.Insert(&testCardData1)
	o.Insert(&testCardData2)
	o.Insert(&testCardData3)
	o.Insert(&enterprise)

	_, _ = o.Insert(&testEnterpriseData1)
	_, _ = o.Insert(&testEnterpriseData2)

	_, _ = o.Insert(&ecount)
	for i := 0; i < len(container); i++ {
		_, _ = o.Insert(&container[i])
	}
	for i := 0; i < len(demo); i++ {
		_, _ = o.Insert(&demo[i])
	}
}
