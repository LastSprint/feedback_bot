package Entry

import "time"

// ConfusingMessage is models which is stored in DB
type ConfusingMessage struct {
	// AuthorId is UserID of message author (who wrote it)
	AuthorId string `bson:"author_id,omitempty"`
	// ReporterUserId is UserID of reporter (person why marked this message as confusing)
	ReporterUserId string `bson:"reporter_user_id,omitempty"`
	// Text is just a text of the message
	Text string `bson:"text,omitempty"`
	// MessageId is a uniq id of the message
	MessageId  string    `bson:"messageId,omitempty"`
	ReportDate time.Time `bson:"reportDate,omitempty"`
}
