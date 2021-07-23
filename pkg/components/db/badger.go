package db

import (
	badger "github.com/dgraph-io/badger/v3"
	log "github.com/sirupsen/logrus"

	"github.com/geomyidia/go-svc-conventions/pkg/components/config"
	"github.com/geomyidia/go-svc-conventions/pkg/components/msgbus"
)

type DB struct {
	Bus    *msgbus.MsgBus
	Config *config.DBConfig
	Conn   *badger.DB
}

func NewDB(cfg *config.Config, bus *msgbus.MsgBus) *DB {
	log.Debug("Setting up database connection ...")
	return &DB{
		Bus:    bus,
		Config: cfg.DB,
	}
}

// Connect ...
func (db *DB) Connect() {
	conn, err := badger.Open(badger.DefaultOptions(db.Config.Directory))
	if err != nil {
		log.Fatal(err)
	}
	db.Conn = conn
	log.Infof("Connected to database: %s", db.Config.Directory)
}

// Shutdown ...
func (db *DB) Shutdown() {
	db.Conn.Close()
	log.Debugf("DB has been shutdown.")
}
