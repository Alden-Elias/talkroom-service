package models

type Messages struct {
	Messages    []MsgItem `bson:"messages" json:"messages"`
	UnreadCount int       `bson:"-" json:"unreadCount"`
}

type MsgItem struct {
	From uint   `bson:"from" json:"from"`
	Time int64  `bson:"time" json:"time"`
	Msg  string `bson:"msg" json:"msg"`
}

type SystemMsg struct {
	EventId uint `json:"eventId"`
	Json    any  `json:"json"`
}
