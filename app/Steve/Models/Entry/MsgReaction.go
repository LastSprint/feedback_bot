package Entry

type MsgReaction struct {
	ChannelId string         `bson:"channel_id,omitempty"`
	MessageId string         `bson:"message_id,omitempty"`
	Reactions map[string]int `bson:"reactions,omitempty"`
	Year      int            `bson:"year"`
	Week      int            `bson:"week"`
}
