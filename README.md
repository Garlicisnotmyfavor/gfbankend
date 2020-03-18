## GF需求的后端分析

#### 功能分析

- 用户登陆
- 用户对卡片的信息调整
  - 添加卡片
    - 扫描卡片二维码添加卡片（卡号）
    - 填写卡片卡号和商店名字，添加成功（是否成功要判断卡号解析出来的公司名和输入是否匹配）
  - 修改卡片
    - 修改卡片卡号和商店名字 （同上规则）
  - 删除卡片
  - 恢复卡片
- 用户对卡片的使用
  - nfc扫描——积分卡累计星星（只是对其中的一种卡么？）
  - 返回卡的信息——点击查看时（公司信息，积分里还分积分类型，券的类型）
  - 优惠券的使用
  - 卡片超过时间就注销

#### 数据细节

- 用户信息内容

  - 用户id——8位

- 卡片信息

  - 卡号：16位——

    ```
    var CardParseMaps = CardParseStruct{
    	map[string]string{
    		"001": "ANZ",
    		"002": "Calvin Klein",
    		"003": "Starbucks",
    		"004": "Subway",
    	},//店铺
    	map[string]string{
    		"1": "Recharge",
    		"2": "Integrate",
    		"3": "Discount",
    		"4": "RechargeIntegral",
    		"5": "RechargeDiscount",
    		"6": "IntegralDiscount",
    		"7": "RechargeIntegralDiscount",
    	},//卡片类型
    	map[string]string{
    		"1": "New South Wales",
    		"2": "Queensland",
    		"3": "South Australia",
    		"4": "Tasmania",
    		"5": "Victoria",
    		"6": "Western Australia",
    		"7": "Australia Capital Territory",
    		"8": "Northern Territory",
    	},//州
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
    	},//城市
    }
    ```

    

  - 店铺名称（在卡号里）

  - 卡片类型

    - 类型1
      - 积分数（类型，数量）
      - 券（类型，数量）
    - 类型2
      - 二维码
    - 类型3
      - 条形码

- 卡片使用详情

  - 类型1：数量的变化

  

#### 数据库设计

- user

  | user_id(key) |      |      |      |      |
  | ------------ | ---- | ---- | ---- | ---- |
  |              |      |      |      |      |
  |              |      |      |      |      |
  |              |      |      |      |      |

  

- card

  | card_id(key) | user_id | type            | type_id | enterprise | state | city | money | score_num | coupons_num | expire_time |
  | ------------ | ------- | --------------- | ------- | ---------- | ----- | ---- | ----- | --------- | ----------- | ----------- |
  | 123456       |         | membership_card | 1       | uestc      |       |      | 100   | 3         | 0           | 12/12/2020  |
  | 123456       |         | type1           | 2       | uestc      |       |      | 100   | 0         | 3           | 12/12/2020  |
  |              |         |                 |         |            |       |      |       |           |             |             |

- enterprise

  | enterprise(key) | tel         | location | time     | type            | type_id |
  | --------------- | ----------- | -------- | -------- | --------------- | ------- |
  | uestc           | 13880059462 | china    | 2/2/2020 | membership_card | 1       |
  | uestc           |             |          |          | membership_card | 2       |
  | uestc           |             |          |          | type3           | 1       |

- membership_card

  | type_id | coupons_discripe | point_rule | point_discripe   | time |
  | ------- | ---------------- | ---------- | ---------------- | ---- |
  | 1       | null             | 5          | buy 5 get 1 free | 30   |
  | 2       | free coffee      | null       | null             | 30   |
  | 3       | null             | 4          | buy 4 get 1 free | 30   |

- type2

  | type_id |      |      |      |      |
  | ------- | ---- | ---- | ---- | ---- |
  |         |      |      |      |      |
  |         |      |      |      |      |
  |         |      |      |      |      |

- type3

  | type_id |      |      |      |      |
  | ------- | ---- | ---- | ---- | ---- |
  |         |      |      |      |      |
  |         |      |      |      |      |
  |         |      |      |      |      |

  