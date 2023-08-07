package main

import (
	"dgateway/model"
	"fmt"
	goweb "github.com/Rehtt/Kit/web"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	registerExpiration = 11 * time.Minute
)

func dgApi() {
	g := goweb.New()
	g.POST("/api/reg", register)
	if err := http.ListenAndServe(":8001", g); err != nil {
		panic(err)
	}
}
func register(ctx *goweb.Context) {
	var body model.Register
	if err := ctx.ReadJSON(&body); err != nil {
		ctx.WriteJSON(&model.Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	if body.Uid == "" {
		ctx.WriteJSON(&model.Response{Error: "uid is null"}, http.StatusBadRequest)
		return
	}
	if body.Scheme == "" {
		body.Scheme = "http"
	}
	body.Remote = body.Scheme + "://" + strings.Split(ctx.Request.RemoteAddr, ":")[0] + ":" + strconv.Itoa(body.Port)
	newb, _ := jsoniter.MarshalToString(body)

	var errs []string
	for _, route := range body.Routes {
		key := fmt.Sprintf("%s|%s", strings.ToTitle(route.Method), route.Uri)
		if value := rdb.Get(ctx, key).Val(); value == "" || jsoniter.Get([]byte(value), "uid").ToString() == body.Uid {
			rdb.Set(ctx, key, newb, registerExpiration)
			continue
		}
		errs = append(errs, key)
	}
	if len(errs) > 0 {
		ctx.WriteJSON(&model.Response{Error: fmt.Sprintf("路由冲突：%s", strings.Join(errs, ","))}, http.StatusOK)
		return
	}
	ctx.WriteJSON(&model.Response{}, http.StatusOK)
	return
}
