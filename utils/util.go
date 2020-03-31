package util

import (
	"errors"
	"fmt"
	"github.com/gfbankend/models"
	"github.com/go-gomail/gomail"
	"log"
	"math/rand"
	"strings"
	"time"
)

/*
*@function:得到6位长的验证码
*@return {[]byte} 验证码
 */
func GetRandCode() []byte {
	var code []byte
	number := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Seed(time.Now().Unix())
	var sb strings.Builder
	size := len(number)
	for i := 0; i < 6; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", number[rand.Intn(size)])
	}
	code = []byte(sb.String())
	return code
}

/*
*@function:发送验证码给target邮箱
*@param {string} 目标邮箱
*@return {[]byte}vcode，{error}err
 */
func SendEmail(target string) ([]byte, error) {
	//产生验证码
	vcode := GetRandCode()
	if len(vcode) != 6 {
		models.Log.Error("Error generating verify code")
		return nil, errors.New("Fail to generate verify code!")
	}
	//邮箱内容
	content := fmt.Sprintf("[ANZ]尊敬的客户' %s '，您本次登录所需的验证码为:%s,请勿向任何人提供您收到的验证码!", target, vcode)
	m := gomail.NewMessage()
	//设置邮件信息
	m.SetAddressHeader("From", "gfbankend@163.com", "ANZ-WORKSHOP") //设置发件人
	m.SetHeader("Subject", "Verify your device")                    //设置主题
	m.SetBody("text/html", content)                                 //设置主体内容
	m.SetHeader("To", m.FormatAddress(target, "收件人"))               //设置收件人
	//连接邮箱服务器并发送邮件
	d := gomail.NewPlainDialer("smtp.163.com", 465, "gfbankend@163.com", "ahz12345")

	if err := d.DialAndSend(m); err != nil {
		log.Println("Fail to send: ", err)
		return nil, err
	}
	return vcode, nil
}

