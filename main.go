package main

import (
	"flag"
	"fmt"
	"github.com/Rehtt/Kit/util"
	goweb "github.com/Rehtt/Kit/web"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	Version string
)

func init() {
	showVersion := flag.Bool("v", false, "version")
	redisAddr := flag.String("rdb_addr", util.Getenv("RDB_ADDR", "127.0.0.1:6379"), "redis connect addr")
	redisPassword := flag.String("rdb_password", util.Getenv("RDB_PASSWORD", ""), "redis password")
	redisUsername := flag.String("rdb_username", util.Getenv("RDB_USERNAME", ""), "redis username")
	redisDB := flag.String("rdb_db", util.Getenv("RDB_DB", "0"), "redis db")
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	db, err := strconv.Atoi(*redisDB)
	if err != nil {
		panic("rdb_db error")
	}
	if err := InitRedis(*redisAddr, *redisUsername, *redisPassword, db); err != nil {
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
