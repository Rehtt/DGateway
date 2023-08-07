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
	flag.Parse()
	if *showVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
}

func main() {
	log.Println("run")
	if err := InitRedis(); err != nil {
		panic(err)
	}
	go dgApi()

	g := goweb.New()
	g.Any("/#", gateway)
	if err := http.ListenAndServe(":80", g); err != nil {
		panic(err)
	}
}
