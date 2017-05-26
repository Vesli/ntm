package main

/*
	Service structure that manage the configuration and the session.
	It initialise the config struture and the session
*/

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/vesli/ntm/config"
)

type service struct {
	Config *config.Config
	DB     *gorm.DB
}

func newService(pathToConfig string) (*service, error) {
	conf, err := config.LoadConfig(pathToConfig)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(
		"postgres",
		"host="+conf.DBURL+" user="+conf.DBUser+" dbname="+conf.DBName+" sslmode=disable password="+conf.DBPassword)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{}, &Event{})

	service := &service{
		Config: conf,
		DB:     db,
	}

	return service, nil
}

func (s *service) Close() {
	s.DB.Close()
}
