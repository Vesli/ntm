package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"github.com/ntm/models"
)

// Config API
type Config struct {
	Port string

	FBURL    string
	FBFields string

	DBURL      string
	DBUser     string
	DBName     string
	DBPassword string
}

// Service api
type Service struct {
	Conf   *Config
	Router *chi.Mux
	DB     *gorm.DB
}

func loadEnvironment() *Config {
	c := &Config{
		Port:       ":" + os.Getenv("PORT"),
		DBURL:      os.Getenv("DBURL"),
		DBUser:     os.Getenv("DBUSER"),
		DBName:     os.Getenv("DBNAME"),
		DBPassword: os.Getenv("DBPWD"),
	}

	return c
}

func prepareDB(conf *Config) *gorm.DB {
	db, err := gorm.Open(
		"postgres",
		"host="+conf.DBURL+" user="+conf.DBUser+" dbname="+conf.DBName+" sslmode=disable password="+conf.DBPassword)
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&models.User{}, &models.Event{})

	return db
}

func (s *Service) run() {
	log.Println("Running on port:", s.Conf.Port)
	log.Fatal(http.ListenAndServe(s.Conf.Port, s.Router))
}

func (s *Service) close() {
	s.DB.Close()
}

func newService() *Service {
	s := &Service{
		Conf:   loadEnvironment(),
		Router: chi.NewRouter(),
	}

	s.DB = prepareDB(s.Conf)
	s.Router.Mount("/api", routes(s))

	return s
}
