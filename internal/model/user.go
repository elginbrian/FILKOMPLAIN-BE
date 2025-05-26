package model

type User struct {
	ID              uint   `gorm:"primaryKey"`
	Username        string `gorm:"unique"`
	Email           string `gorm:"unique"`
	Password        string
	Type            string `json:"type" gorm:"type:varchar(16);default:'user'"`
	NIM             string `json:"nim" gorm:"type:varchar(20)"`
	ProgramStudi    string `json:"program_studi" gorm:"type:varchar(100)"`
	PhoneNumber     string `json:"phone_number" gorm:"type:varchar(20)"`
	ProfileImageURL string `json:"profile_image_url" gorm:"type:varchar(255);default:''"`
}
