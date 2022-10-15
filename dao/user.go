package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"talkRoom/models"
)

var (
	EmailUsed = errors.New("邮箱已使用")
)

//UserAdd 创建用户
func UserAdd(user *models.User) error {
	if err := mysqlDb.Where("email=?", user.Email).First(&models.User{}).Error; err != gorm.ErrRecordNotFound {
		return EmailUsed
	} else if err = mysqlDb.Create(user).Error; err != nil {
		return err
	} else if err = AddFriend(user.ID, user.ID); err != nil { //添加自己为好友
		return err
	}
	user.Name = "用户" + strconv.Itoa(int(user.ID)) //初始名（可以通过数据库inserted实现的）
	return mysqlDb.Save(user).Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := mysqlDb.First(&user, "email=?", email).Error
	return &user, err
}

func GetUserById(id uint) (*models.User, error) {
	var user models.User
	err := mysqlDb.First(&user, "id=?", id).Error
	return &user, err
}

//UserFuzzySearch 用户模糊查找
func UserFuzzySearch(s string) (users *[]models.UserBrief, err error) {
	var us []models.UserBrief
	if id, err := strconv.Atoi(s); err != nil {
		err = mysqlDb.Debug().Model(&models.User{}).Find(&us, fmt.Sprint("name LIKE '%", s, "%'")).Error
	} else {
		err = mysqlDb.Debug().Model(&models.User{}).Find(&us, fmt.Sprint("id = ", id, " OR name LIKE '%", s, "%'")).Error
	}

	users = &us
	log.Println(users)
	return
}

//UserSave 保存修改用户信息
func UserSave(user *models.User) error {
	return mysqlDb.Save(user).Error
}
