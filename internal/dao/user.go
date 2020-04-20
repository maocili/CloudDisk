package dao

import (
	"CloudDisk/internal/model"
	"CloudDisk/tools/mysql"
	"errors"

	"github.com/go-xorm/xorm"
)

func Signup(userinfo model.UserInfo) error {

	_, err := mysql.DB.Transaction(func(session *xorm.Session) (interface{}, error) {

		//检查账户是否已存在
		has_account, err := session.Table(&userinfo).Where("username = ?", userinfo.Username).Exist()
		// has_phoneNumber, err := session.Table(&userinfo).Where("phone_number = ?", user_info.PhoneNumber).Exist()
		switch {
		case has_account == true:
			return nil, errors.New("账户已存在")
		// case has_phoneNumber:
		// 	return nil, errors.New("手机号已注册")
		case err != nil:
			return nil, err
		}

		if _, err := session.Insert(&userinfo); err != nil {
			println(err.Error())
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func IsexistUsername(username string) (bool, error) {

	var userinfo model.UserInfo
	has_account, err := mysql.DB.Table(&userinfo).Where("username = ?", username).Exist()
	return has_account, err
}
