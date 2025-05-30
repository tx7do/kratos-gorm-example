package models

import "gorm.io/gorm"

// User 用户信息表
type User struct {
	gorm.Model
	UserName string `gorm:"column:username;comment:'账号名'"`
	NickName string `gorm:"column:nickname;comment:'昵称'"`
	Password string `gorm:"column:password;comment:'登录密码'"`
}

func (u User) TableName() string {
	return "users"
}
