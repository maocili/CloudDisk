package model

type UserInfo struct {
	Uid      string `form : "uid"`
	Username string `form : " username"`
	Password string `form:" password"`
}
