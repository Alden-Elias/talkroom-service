package models

type FriendRelationship struct {
	UserID   uint `gorm:"primaryKey;ForeignKey:UserID;"` //表示关系的主体
	FriendID uint `gorm:"primaryKey;ForeignKey:UserID;"`
}
