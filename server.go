package main

import (
	"log"
	"net/http"
	"strconv"
)

func startServer(config *Config) {
	err := initMongoDB(config)
	defer closeMongoDB()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	log.Printf("Listening on port %d\n", config.HTTP.Port)
	err = http.ListenAndServeTLS(":"+strconv.Itoa(config.HTTP.Port), config.HTTP.HttpsCertFile, config.HTTP.HttpsKeyFile, &Handler{})
	log.Fatal(err)
}
