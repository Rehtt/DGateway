package main

import (
	"flag"
	"fmt"
	goweb "github.com/Rehtt/Kit/web"
	"log"
	"net/http"
	"os"
)

var (
	Version string
)

func init() {
	showVersion := flag.Bool("v", false, "version")
	redisAddr := flag.String("rdb_addr", "127.0.0.1:6379", "redis connect addr")
	redisPassword := flag.String("rdb_password", "", "redis password")
	redisUsername := flag.String("rdb_username", "", "redis username")
	redisDB := flag.Int("rdb_db", 0, "redis db")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if err := InitRedis(*redisAddr, *redisUsername, *redisPassword, *redisDB); err != nil {
		panic(err)
	}
}

func main() {
	log.Println("DGateway - Rehtt")
	log.Println("version:", Version)
	log.Println("running")

	go dgApi()

	g := goweb.New()
	g.Any("/#", gateway)
	if err := http.ListenAndServe(":80", g); err != nil {
		panic(err)
	}
}
