package models

type AuthData struct {
	IDUser int    `json:"id_user"`
	Token  string `json:"token"`
}

type RegisterData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
