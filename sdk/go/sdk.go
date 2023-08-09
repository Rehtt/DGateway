package dgateway

import (
	"bytes"
	"errors"
	"github.com/Rehtt/DGateway/model"
	"github.com/Rehtt/Kit/util"
	goweb "github.com/Rehtt/Kit/web"
	jsoniter "github.com/json-iterator/go"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewHttpServer() *goweb.GOweb {
	return goweb.New()
}

func ListenAndServe(addr string, g *goweb.GOweb) (err error) {
	port := 80
	if s := strings.Split(addr, ":"); len(s) == 2 {
		port, err = strconv.Atoi(s[1])
		if err != nil {
			return err
		}
	}
	go func(port int, g *goweb.GOweb) {
		t := time.NewTicker(10 * time.Minute)
		method, list := g.List()
		for {
			reg := model.Register{
				Uid:    "",
				Port:   port,
				Scheme: "http",
				Routes: make([]*model.Route, 0, len(list)),
			}
			for i := range list {
				reg.Routes = append(reg.Routes, &model.Route{
					Method: method[i],
					Uri:    list[i],
				})
			}
			if err := Register(model.Register{
				Uid:    "",
				Port:   port,
				Scheme: "http",
				Routes: nil,
			}, util.Getenv("DG_BASE", "http://dgateway:8080")); err != nil {
				log.Println("register dgateway error:", err)
			}
			<-t.C
		}
	}(port, g)

	return http.ListenAndServe(addr, g)
}

func Register(reg model.Register, dgatewayBasePath string) error {
	var tmp bytes.Buffer
	if err := jsoniter.NewEncoder(&tmp).Encode(reg); err != nil {
		return err
	}
	resp, err := http.Post(dgatewayBasePath+"/api/reg", "application/json", &tmp)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var out model.Response
	if err := jsoniter.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}
	if out.Error != "" {
		return errors.New(out.Error)
	}
	return nil
}
