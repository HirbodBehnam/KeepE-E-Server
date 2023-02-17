package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func getDatabaseConnectionURL() string {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatalln("please set DATABASE_URL environmental variable")
	}
	return dbUrl
}

func getListenAddress() string {
	listen := os.Getenv("HTTP_LISTEN")
	if listen == "" {
		log.Fatalln("please set HTTP_LISTEN environmental variable")
	}
	return listen
}
