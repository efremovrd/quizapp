package models

type User struct {
	Id, Login, Password string
}

func (u *User) EqPasswords(password string) (res bool) {
	if u.Password == password {
		res = true
	}
	return
}
