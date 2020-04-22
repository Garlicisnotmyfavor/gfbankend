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
	//testData()
}

func testData(){
	//testUserData1 := User{Id:"2018091620000",Tel:"13925678240",Mail:"123456@qq.com",Password:"123456789",LoginNum:0}
	//testCardData1 := Card{CardId:"123456790123456",UserId:"2018091620000",State:"California",City:"San Jose", CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()}
	//testUserData2 := User{Id:"2018091620001",Tel:"13665372240",Mail:"654321@qq.com",Password:"123908789",LoginNum:1}
	//testCardData2 := Card{CardId:"123456790000000",UserId:"2018091620001",State:"California",City:"San Jose",CardType:"Discount",Enterprise:"subway",ExpireTime:time.Now()}
	//testUserData3 := User{Id:"2018091620002",Tel:"13778372240",Mail:"666666@qq.com",Password:"123009889",LoginNum:2}
	//testCardData3 := Card{CardId:"123456790333000",UserId:"2018091620002",State:"California",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()}
	//testUserData4 := User{Id:"2018091620020",Tel:"13778788240",Mail:"78902166@qq.com",Password:"33235323",LoginNum:20}
	container := []Card {
		Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333002",State:"Cambera",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333003",State:"Beijing",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333004",State:"Chengdu",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333005",State:"Guangzhou",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333006",State:"Foshan",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333007",State:"Hangzhou",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333008",State:"Shanghai",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333009",State:"Shandong",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333010",State:"Xian",City:"San Jose",CardType:"RechargeDiscount",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333011",State:"Mianyang",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333012",State:"Changchun",City:"San Jose",CardType:"IntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333013",State:"Nanjing",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now()},
		Card{CardId:"123456790333014",State:"Tokyo",City:"San Jose",CardType:"RechargeIntegralDiscount",Enterprise:"starbucks",ExpireTime:time.Now()},
	}
	// testCardData4 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()}
	// testCardData5 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()}
	// testCardData6 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now()}
	// testCardData7 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Discount",Enterprise:"starbucks",ExpireTime:time.Now()}
	// testCardData8 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now()}
	// testCardData9 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"Recharge",Enterprise:"starbucks",ExpireTime:time.Now()}
	// testCardData10 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now()}
	// testCardData11 := Card{CardId:"123456790333001",State:"California",City:"San Jose",CardType:"RechargeIntegral",Enterprise:"starbucks",ExpireTime:time.Now()}
	testEnterpriseData1 := Enterprise{Id:"001",Name:"starbucks"}
	testEnterpriseData2 := Enterprise{Id:"002",Name:"subway"}
	o := orm.NewOrm()
	_, _ = o.Insert(&testEnterpriseData1)
	_, _ = o.Insert(&testEnterpriseData2)
	_, _ = o.InsertMulti(len(container), container)
}
