package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"unique"`
	CID    string `json:"cid" gorm:"unique"`
	Phone     string `json:"phone" gorm:"unique"`
	Password []byte `json:"-"`

}

type Product struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"unique"`
	PID     string `json:"pid" gorm:"unique"`
	Price    string `json:"price"`
	Imagesrc []byte `json:"image"`

}

type Order struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	TID     string `json:"tid" gorm:"unique"`
	CID     string `json:"name"`
	PID    string `json:"price"`
	Imagesrc []byte `json:"image"`	

}

type Wallet struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	WID     string `json:"wid" gorm:"unique"`
	CID    string `json:"cid"`

}

