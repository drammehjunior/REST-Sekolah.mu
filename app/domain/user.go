package domain

type Users struct {
	Id        int64  `json:"id" gorm:"primary_key; auto_increment; not_null"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
