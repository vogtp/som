package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"github.com/vogtp/go-hcl"
	"github.com/vogtp/som/pkg/core/cfg"
	"github.com/vogtp/som/pkg/core/msg"
	"gorm.io/gorm"
)

const dbName = "som.sqlite"

// Access wraps the DB
type Access struct {
	hcl hcl.Logger
	db  *gorm.DB
}

func (a *Access) getDb() *gorm.DB {
	if a.db == nil {
		a.hcl = hcl.New() //FIXME remove
		db, err := a._initDB()
		if err != nil {
			panic(fmt.Errorf("cannot init DB: %w", err))
		}
		a.db = db
	}
	return a.db
}

func (a *Access) _initDB() (*gorm.DB, error) {
	path := fmt.Sprintf("%s/%s", viper.GetString(cfg.DataDir), dbName)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot open DB %q: %w", path, err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&IncidentModel{},
		&statiModel{},
		&ErrorModel{},
		&counterModel{},
		&msg.FileMsgItem{},
	)
	a.hcl.Infof("Loaded DB %s", path)
	return db, err
}
