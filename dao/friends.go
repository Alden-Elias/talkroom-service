package dao

import (
	"gorm.io/gorm"
	"talkRoom/models"
)

func AddFriend(uid, fid uint) error {
	err := mysqlDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&models.FriendRelationship{
			UserID:   uid,
			FriendID: fid,
		}).Error; err != nil {
			return err
		}
		return tx.Save(&models.FriendRelationship{
			UserID:   fid,
			FriendID: uid,
		}).Error
	})
	if err != nil {
		return err
	}
	return CreateMassagesStorage(uid, fid)
}

func ListFriendsBrief(uid uint) (friends *[]models.UserBrief, err error) {
	var fIds []uint
	var f []models.UserBrief
	if err = mysqlDb.Model(&models.FriendRelationship{}).Select("friend_id").Find(&fIds, "user_id=?", uid).Error; err != nil {
		return
	}
	err = mysqlDb.Model(&models.User{}).Find(&f, "id IN ?", fIds).Error
	friends = &f
	return
}

func IsFriend(uid, fid uint) bool {
	if err := mysqlDb.First(&models.FriendRelationship{}, "user_id=? AND friend_id=?", uid, fid).Error; err != nil {
		return false
	}
	return true
}
