package main

import (
	"github.com/vesli/ntm/config"
	mgo "gopkg.in/mgo.v2"
)

type service struct {
	Config  *config.Config
	Session *mgo.Session
}

func newService(pathToConfig string) (*service, error) {
	conf, err := config.LoadConfig(pathToConfig)
	if err != nil {
		return nil, err
	}

	session, err := mgo.Dial(conf.DbURL)
	if err != nil {
		return nil, err
	}

	service := &service{
		Config:  conf,
		Session: session,
	}

	return service, nil
}

func (s *service) Close() {
	s.Session.Close()
}
