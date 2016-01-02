package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fiorix/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configFile := flag.String("c", "server.conf", "")
	flag.Usage = func() {
		fmt.Println("Usage: server [-c server.conf] ")
		os.Exit(1)
	}
	flag.Parse()

	var err error
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open("sqlite3", config.DB.DBConn)
	if err != nil {
		log.Fatal(err)
	}
	rc := redis.New(config.DB.Redis)

	_ = NewApiRouter(config, rc, &db)

	log.Fatal(http.ListenAndServe(config.HTTP.Addr, nil))
	log.Printf("Server ready")
}
