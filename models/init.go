package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
	//"github.com/gfbankend/models"
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
	orm.Debug = true
	_ = orm.RunSyncdb("default", false, true)
	testData()
}

func testData(){
	testUserData1 := User{Id:"2018091620000",Tel:"13925678240",Mail:"123456@qq.com",Password:"123456789",LoginNum:0,LoginYear:"2020",LoginMonth:"1"}
	testCardData1 := Card{CardId:"123456790123456",UserId:"2018091620000",State:"California",City:"San Jose", CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty"}
	testUserData2 := User{Id:"2018091620001",Tel:"13665372240",Mail:"654321@qq.com",Password:"123908789",LoginNum:1,LoginYear:"2020",LoginMonth:"1"}
	testCardData2 := Card{CardId:"123456790000000",UserId:"2018091620001",State:"California",City:"San Jose",CardType:"Discount",Enterprise:"subway",ExpireTime:time.Now(),Coupons:"empty"}
	testUserData3 := User{Id:"2018091620002",Tel:"13778372240",Mail:"666666@qq.com",Password:"123009889",LoginNum:2,LoginYear:"2020",LoginMonth:"1"}
	testCardData3 := Card{CardId:"123456790333000",UserId:"2018091620002",State:"California",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty"}
	testUserData4 := User{Id:"2018091620020",Tel:"13778788240",Mail:"78902166@qq.com",Password:"33235323",LoginNum:20,LoginYear:"2020",LoginMonth:"1"}
	container := []Card {
		Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333002",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"1234567903330015",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"1234567903330016",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"1234567903330017",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"1234567903330018",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"1234567903330019",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"1234567903330020",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333003",State:"Beijing",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333004",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333021",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333022",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333023",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333024",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333025",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333026",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333005",State:"Guangzhou",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333006",State:"Foshan",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333027",State:"Foshan",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333028",State:"Foshan",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333029",State:"Foshan",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333030",State:"Foshan",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333031",State:"Foshan",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333032",State:"Hangzhou",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333008",State:"Shanghai",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333033",State:"Shanghai",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333034",State:"Shanghai",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333035",State:"Shanghai",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333036",State:"Shanghai",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333037",State:"Shanghai",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333038",State:"Shandong",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333010",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333039",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333040",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333041",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333042",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333043",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333044",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333045",State:"Mianyang",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333012",State:"Changchun",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333046",State:"Changchun",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333047",State:"Changchun",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333048",State:"Changchun",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333049",State:"Changchun",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333050",State:"Changchun",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333013",State:"Nanjing",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333014",State:"Tokyo",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333051",State:"Nanjing",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333052",State:"Nanjing",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333053",State:"Nanjing",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333054",State:"Nanjing",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
		Card{CardId:"123456790333055",State:"Nanjing",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now(),Coupons:"empty",UserId:"empty"},
	}
	testEnterpriseData1 := Enterprise{Id:"001",Name:"starbucks",Addr:"empty",IsLocal:false,Type:"empty",HelpMsg:"empty",Website:"empty",LicenseId:"empty"}
	testEnterpriseData2 := Enterprise{Id:"002",Name:"subway",Addr:"empty",IsLocal:false,Type:"empty",HelpMsg:"empty",Website:"empty",LicenseId:"empty"}
	o := orm.NewOrm()
	o.Insert(&testUserData1)
	o.Insert(&testUserData2)
	o.Insert(&testUserData3)
	o.Insert(&testUserData4)
	o.Insert(&testCardData1)
	o.Insert(&testCardData2)
	o.Insert(&testCardData3)
	_, _ = o.Insert(&testEnterpriseData1)
	_, _ = o.Insert(&testEnterpriseData2)
	for i:=0 ;i < len(container); i++ {
		_,_ = o.Insert(&container[i])
	}
}
