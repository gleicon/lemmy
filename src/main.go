package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fiorix/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

type Backend struct {
	Host   string
	Port   int
	Weight int
	Active bool
}

type VHosts struct {
	LastUpdated time.Time
	Backends    map[string][]Backend
}

func monitorVhosts(config *configFile, rc *redis.Client) {

	for {
		vh, err := rc.SMembers(config.LoadBalancer.Prefix + ":vhosts")
		if err != nil {
			log.Printf("Error loading vhost list: %s", err)
			time.Sleep(time.Duration(config.LoadBalancer.VHostRefreshTime) * time.Second)
			continue
		}
		for vhost := range vh {
			vd, err := rc.ZRangeByScore(config.LoadBalancer.Prefix+":vhosts", 0, config.LoadBalancer.MaxBackends, true, false, 0, 0)
			if err != nil {
				log.Printf("Error fetching backends for %s: %#v", vhost, err)
				continue
			}
			log.Println(vd)
			log.Println(vhost)
		}
		time.Sleep(time.Duration(config.LoadBalancer.VHostRefreshTime) * time.Second)
		log.Println("VHost refreshed")
	}
}

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

	go monitorVhosts(config, rc)

	_ = NewApiRouter(config, rc, &db)

	log.Fatal(http.ListenAndServe(config.HTTP.Addr, nil))
	log.Printf("Server ready")
}
