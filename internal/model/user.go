package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
	Type     string `json:"type" gorm:"type:varchar(16);default:'user'"`
}
