package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
)

type User struct {
	Userid       int64  `json:"userid" orm:"pk;column(userid);auto"`
	UserName     string `json:"username" orm:"column(username);size(50);null"`
	Password     string `json:"_" orm:"column(possword);size(50);null"`
	Salt         string `json:"_"  orm:"column(salt);size(32)"`
	RegisterTime int64  `json:"_" orm:"column(register_time);null"`
	UpdateTime   int64  `json:"_" orm:"column(update_time);null" `
	LoginTime    int64  `json:"login_time" orm:"column(login_time);null" `
	LoginIp      string `json:"login_ip" orm:"column(login_ip);null" `
	Sessionid    string `json:"sessionid" orm:"column(sessionid);size(32);null"`
	Enable       bool   `json:"enable"  orm:"column(enable)"`
}

func (u *User) TableEngine() string {
	return "INNODB"
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) GetUser() []*User {
	var User []*User
	num, err := Orm.QueryTable("user").All(&User)
	if err != nil {
		beego.Debug("【查询User错误", err.Error(), "]")
	}
	beego.Debug("【查询user的记录数据是", num, "]条")
	return User
}

func (u *User) GetUserById(Userid int64) (User, error) {
	var User User
	//查询一个，不能使用eamil *Group
	err := Orm.QueryTable("user").Filter("userid", Userid).One(&User)
	if err != nil {
		beego.Debug("查询user错误", err.Error())
		return User, err
	}
	return User, nil
}

func (u *User) GetUserByName(username string) User {
	var user User
	//查询一个，不能使用eamil *Group
	err := Orm.QueryTable("user").Filter("username", username).One(&user)
	if err != nil {
		beego.Debug("查询user错误", err.Error())
		return user
	}
	return user
}

func (u *User) AddUser() (int64, error) {
	u.RegisterTime = time.Now().Unix()
	u.UpdateTime = u.RegisterTime
	count, err := Orm.QueryTable("user").Filter("username", u.UserName).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		err = errors.New("用户名已经存在")
		return 0, err
	}
	fmt.Println(&u)
	fmt.Println(u)
	id, err := Orm.Insert(u)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *User) UpdateUser() (err error) {
	_, err = Orm.Update(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) DelUser(id int64) error {
	if id == 0 {
		return errors.New("内置用户不允许删除")
	}
	_, err := Orm.QueryTable("user").Filter("userid", id).Delete()
	if err != nil {
		beego.Error("删除User错误", err.Error())
		return err
	}
	return nil
}
