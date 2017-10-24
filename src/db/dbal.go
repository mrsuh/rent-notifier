package dbal

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Recipient struct {
	Id             bson.ObjectId `bson:"_id"`
	TelegramChatId string        `bson:"telegram_chat_id"`
	City           int           `bson:"city"`
	Subways        []int         `bson:"subways"`
	Type           int           `bson:"type"`
}

type Note struct {
	Subways     []int    `json:"subways"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	Type        int      `json:"type"`
	Link        string   `json:"link"`
	Photos      []string `json:"photos"`
	City        int      `json:"city"`
	Contact     string   `json:"contact"`
}

type Notification struct {
	Note       Note
	Recipients []Recipient
}

type DBAL struct {
	session *mgo.Session
	db      *mgo.Database
}

func (dbal *DBAL) AddRecipient(recipient Recipient) {
	recipient.Id = bson.NewObjectId()
	dbal.db.C("recipients").Insert(&recipient)
}

func (dbal *DBAL) RemoveRecipient(recipient Recipient) {
	dbal.db.C("recipients").Remove(bson.M{"_id": recipient.Id})
}

func (dbal *DBAL) FindRecipientsByNote(note Note) []Recipient {
	result := []Recipient{}
	dbal.db.C("recipients").Find(bson.M{"City": note.City, "Subways": 1, "Type": note.Type}).All(&result) //todo subways

	return result
}

func Connect(dsn string) *DBAL {
	session, err := mgo.Dial(dsn)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return &DBAL{session, session.DB("go-test")}
}
