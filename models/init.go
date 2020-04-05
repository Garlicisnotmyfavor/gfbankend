package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"os"
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
	//testData()
	orm.Debug = true
	_ = orm.RunSyncdb("default", false, true)
	//testData()
}

func testData(){
	testUserData1 := User{Id:"2018091620000",Tel:"13925678240",Mail:"123456@qq.com",Password:"123456789"}
	testCardData1 := Card{CardId:"123456790123456",UserId:"2018091620000",State:"California",City:"San Jose", CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()}
	testUserData2 := User{Id:"2018091620001",Tel:"13665372240",Mail:"654321@qq.com",Password:"123908789"}
	testCardData2 := Card{CardId:"123456790000000",UserId:"2018091620001",State:"California",City:"San Jose",CardType:"Discount",Enterprise:"subway",ExpireTime:time.Now()}
	testUserData3 := User{Id:"2018091620002",Tel:"13778372240",Mail:"666666@qq.com",Password:"123009889"}
	testCardData3 := Card{CardId:"123456790333000",UserId:"2018091620002",State:"California",City:"San Jose",CardType:"Integrate",Enterprise:"starbucks",ExpireTime:time.Now()}
	o := orm.NewOrm()
	o.Insert(&testUserData1)
	o.Insert(&testUserData2)
	o.Insert(&testUserData3)
	o.Insert(&testCardData1)
	o.Insert(&testCardData2)
	o.Insert(&testCardData3)
}
