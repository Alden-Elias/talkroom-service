package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"talkRoom/models"
	"time"
)

var (
	ctx = context.TODO()
)

func CreateMassagesStorage(uid1, uid2 uint) error {
	_, err := messagesCollection.InsertOne(ctx, bson.D{{"_id", unionUint2Uint64(uid1, uid2)}})
	return err
}

func StoryMessage(from, to uint, msg *string) error {
	_, err := messagesCollection.UpdateOne(
		ctx,
		bson.D{{"_id", unionUint2Uint64(from, to)}},
		bson.D{{"$push", bson.D{
			{"messages", models.MsgItem{From: from, Time: time.Now().UnixMilli(), Msg: *msg}}}}})
	if err != nil {
		return err
	}
	if from == to {
		return nil
	}
	return unreadCount(from, to)
}

func GetMessages(from, to uint) (*models.Messages, error) {
	var messages models.Messages

	err := messagesCollection.FindOne(
		ctx,
		bson.D{{"_id", unionUint2Uint64(from, to)}},
		options.FindOne().SetProjection(bson.M{"_id": 0}),
	).Decode(&messages)

	//return &messages, err
	//更新修改部分
	if err != nil {
		return nil, err
	}

	if from == to {
		messages.UnreadCount = 0
	} else if isExists, err := redisDb.Exists("Message:UnreadCount:" + unionUint2Uint64(from, to)).Result(); err != nil {
		return nil, err
	} else if isExists == 0 {
		redisDb.Set("Message:UnreadCount:"+unionUint2Uint64(from, to), 0, 0)
		messages.UnreadCount = 0
	} else if count, err := redisDb.Get("Message:UnreadCount:" + unionUint2Uint64(from, to)).Result(); err != nil {
		return nil, err
	} else if unreadCount, err := strconv.Atoi(count); err != nil {
		return nil, err
	} else {
		messages.UnreadCount = unreadCount
	}

	return &messages, err
}

func unreadCount(from, to uint) error {
	return redisDb.Incr("Message:UnreadCount:" + unionUint2Uint64(from, to)).Err()
}

func CleanUnreadCount(from, to uint) error {
	return redisDb.Set("Message:UnreadCount:"+unionUint2Uint64(from, to), 0, 0).Err()
}

func unionUint2Uint64(a, b uint) string {
	return strconv.FormatUint((uint64(uintMin(a, b)))<<32|uint64(uintMax(a, b)), 16)
}

func uintMax(a, b uint) uint {
	if a < b {
		return b
	}
	return a
}

func uintMin(a, b uint) uint {
	if a > b {
		return b
	}
	return a
}
