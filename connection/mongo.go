package connection

import (
	"fmt"
	"net/url"

	"github.com/fankserver/fankserver-cli/config"

	mgo "gopkg.in/mgo.v2"
)

var (
	mainSession  *mgo.Session
	mainDatabase *mgo.Database
)

type MongoDB struct {
	Session    *mgo.Session
	Database   *mgo.Database
	Collection *mgo.Collection
}

func (m *MongoDB) Init() *mgo.Session {
	conf := config.GetConfig()

	if mainSession == nil {
		u := &url.URL{
			Scheme: "mongodb",
			Host:   fmt.Sprintf("%s:%d", conf.DB["mongo"].Hostname, conf.DB["mongo"].Port),
		}
		if conf.DB["mongo"].Username != "" {
			u.User = url.UserPassword(conf.DB["mongo"].Username, conf.DB["mongo"].Password)
		}

		var err error
		mainSession, err = mgo.Dial(u.String())
		if err != nil {
			panic(err)
		}

		mainSession.SetMode(mgo.Monotonic, true)
		mainDatabase = mainSession.DB(conf.DB["mongo"].Database)
	}

	m.Session = mainSession.Copy()
	m.Database = m.Session.DB(conf.DB["mongo"].Database)

	return m.Session
}

func (m *MongoDB) C(collection string) *mgo.Collection {
	m.Collection = m.Session.DB(config.GetConfig().DB["mongo"].Database).C(collection)
	return m.Collection
}

func (m *MongoDB) Close() {
	m.Session.Close()
}

func (m *MongoDB) DropDb() {
	err := m.Session.DB(config.GetConfig().DB["mongo"].Database).DropDatabase()
	if err != nil {
		panic(err)
	}
}

func (m *MongoDB) RemoveAll(collection string) {
	conf := config.GetConfig()
	m.Session.DB(conf.DB["mongo"].Database).C(collection).RemoveAll(nil)
	m.Collection = m.Session.DB(conf.DB["mongo"].Database).C(collection)
}

func (m *MongoDB) IsDup(err error) bool {
	return mgo.IsDup(err)
}
