package dbal

import (
	"gopkg.in/mgo.v2"
)

type Connection struct {
	Session *mgo.Session
	Database     string
}

func NewConnection(dsn string) *Connection {
	session, err := mgo.Dial(dsn)
	if err != nil {
		panic(err)
	}

	info, err := mgo.ParseURL(dsn)

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)
	session.SetSafe(&mgo.Safe{})

	return &Connection{Session: session, Database: info.Database}
}