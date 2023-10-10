package jobtype

import (
	"time"
)

type PayloadOrderNotice struct {
	Order Order
	Goods []CartGoods
}

type Order struct {
	Oid         int64     `gorm:"column:oid;primaryKey;autoIncrement:true" json:"oid"`
	UID         int64     `gorm:"column:uid;not null" json:"uid"`
	Member      int32     `gorm:"column:member;not null;default:1;comment:用餐人数" json:"member"`
	GoodsInfo   string    `gorm:"column:goods_info;not null" json:"goods_info"`
	State       int32     `gorm:"column:state;not null;comment:状态 0未确认 1已确认 2未买单 3已买单 4签单 5 已取消" json:"state"`
	Money       int32     `gorm:"column:money;not null" json:"money"`
	PreMoney    int32     `gorm:"column:pre_money;not null;comment:订金" json:"pre_money"`
	PreDatetime time.Time `gorm:"column:pre_datetime;default:1970-01-01 00:00:00" json:"pre_datetime"`
	Phone       string    `gorm:"column:phone;not null" json:"phone"`
	Extra       string    `gorm:"column:extra;not null;comment:备注" json:"extra"`
	Desktop     int32     `gorm:"column:desktop;not null;comment:桌号" json:"desktop"`
	DesktopName string    `gorm:"column:desktop_name;not null;comment:桌号" json:"desktop_name"`
	Operation   int64     `gorm:"column:operation;not null;default:1;comment:管理员id" json:"operation"`
	CreateTime  time.Time `gorm:"column:create_time;not null;default:1970-01-01 00:00:00" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;not null;default:1970-01-01 00:00:00" json:"update_time"`
}

type CartGoods struct {
	Group    int64  `json:"group"`
	Gid      int64  `json:"gid"`
	Desktop  int32  `json:"desktop"`
	Num      int32  `json:"num"`
	Name     string `json:"name"`
	SubTitle string `json:"subTitle"`
	Img      string `json:"img"`
	Price    int32  `json:"price"`
}
