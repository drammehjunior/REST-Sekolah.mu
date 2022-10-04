package domain

type Users struct {
	Id        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
