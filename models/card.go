package models

type Card struct {
	CardId string `json:"card_id" orm:"pk"`
	Kind   string `json:"kind"` // 修改了”类型“这个变量的名字由type改为了kind，因为type是关键词
	Remark string `json:"remark"`
}
