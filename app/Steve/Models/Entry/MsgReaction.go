package Entry

type MsgReaction struct {
	ChannelId string   `bson:"channel_id"`
	MessageId string   `bson:"message_id"`
	Reactions []string `bson:"reactions,omitempty"`
}
