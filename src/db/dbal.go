package dbal

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type City struct {
	Id        int    `bson:"_id"`
	Name      string `bson:"name"`
	Regexp    string `bson:"regexp"`
	HasSubway bool   `bson:"has_subway"`
}

type Subway struct {
	Id     int    `bson:"_id"`
	Name   string `bson:"name"`
	Regexp string `bson:"regexp"`
	City   int    `bson:"city"`
}

type Type struct {
	Id     int    `bson:"_id"`
	Regexp string `bson:"regexp"`
}

const (
	RECIPIENT_TELEGRAM = `telegram`
	RECIPIENT_VK       = `vk`
)

type Recipient struct {
	Id       bson.ObjectId `bson:"_id"`
	ChatId   int           `bson:"chat_id"`
	ChatType string        `bson:"chat_type"`
	City     int           `bson:"city"`
	Subways  []int         `bson:"subways"`
	Types    []int         `bson:"types"`
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

type DBAL struct {
	session *mgo.Session
	db      *mgo.Database
	cities []City
}

func (dbal *DBAL) AddRecipient(recipient Recipient) {
	recipient.Id = bson.NewObjectId()
	dbal.db.C("recipients").Insert(&recipient)
}

func (dbal *DBAL) RemoveRecipient(recipient Recipient) {
	dbal.db.C("recipients").Remove(bson.M{"chat_id": recipient.ChatId, "chat_type": recipient.ChatType})
}

func (dbal *DBAL) FindRecipientsByChatIdAndChatType(chatId int, chatType string) []Recipient {
	result := []Recipient{}
	dbal.db.C("recipients").Find(bson.M{"chat_id": chatId, "chat_type": chatType}).All(&result)

	return result
}

func (dbal *DBAL) FindRecipientsByNote(note Note) []Recipient {
	result := []Recipient{}

	conditions := make([]bson.M, 0)
	conditions = append(conditions, bson.M{"city": note.City, "subways": bson.M{"$in": note.Subways}, "types": note.Type})
	conditions = append(conditions, bson.M{"city": note.City, "subways": bson.M{"$size": 0}, "types": note.Type})

	dbal.db.C("recipients").Find(bson.M{"$or": conditions}).All(&result)

	return result
}

func (dbal *DBAL) AddCity(city City) {
	dbal.db.C("cities").Insert(&city)
}

func (dbal *DBAL) FindCities() []City {
	result := []City{}

	if len(dbal.cities) > 0 {
		return dbal.cities
	}

	dbal.db.C("cities").Find(bson.M{}).All(&result)

	dbal.cities = result

	return result
}

func (dbal *DBAL) AddSubway(subway Subway) {
	dbal.db.C("subways").Insert(&subway)
}

func (dbal *DBAL) FindSubwaysByCity(city City) []Subway {
	result := []Subway{}
	dbal.db.C("subways").Find(bson.M{"city": city.Id}).All(&result)

	return result
}

func (dbal *DBAL) FindSubwaysByIds(ids []int) []Subway {
	result := []Subway{}
	dbal.db.C("subways").Find(bson.M{"_id": bson.M{"$in": ids}}).All(&result)

	return result
}

func (dbal *DBAL) FindTypes() []Type {
	result := []Type{}

	result = append(result, Type{Id: 0, Regexp: "комнат"})
	result = append(result, Type{Id: 1, Regexp: "однушк|квартир"})
	result = append(result, Type{Id: 2, Regexp: "двушк|квартир"})
	result = append(result, Type{Id: 3, Regexp: "тр(е|ё)шк|квартир"})
	result = append(result, Type{Id: 4, Regexp: "четыр|квартир"})
	result = append(result, Type{Id: 5, Regexp: "студи|квартир"})

	return result
}

func Connect(dsn string) *DBAL {
	session, err := mgo.Dial(dsn)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	return &DBAL{session: session, db: session.DB("rent-notifier")}
}
