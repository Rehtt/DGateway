package main

import (
	"fmt"
	"github.com/Rehtt/DGateway/model"
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

func uriKey(method, path string) string {
	return fmt.Sprintf("URI|%s|%s", strings.ToTitle(method), path)
}

func dgApi() {
	g := goweb.New()
	api := g.Grep("/api")
	api.POST("/reg", register)
	api.GET("/list", listUri)
	if err := http.ListenAndServe(":8001", g); err != nil {
		panic(err)
	}
}
func listUri(ctx *goweb.Context) {
	keys := rdb.Keys(ctx, "URI|*").Val()
	if len(keys) == 0 {
		return
	}
	var out = make([]model.Register, len(keys))
	for i, k := range keys {
		jsoniter.UnmarshalFromString(rdb.Get(ctx, k).Val(), &out[i])
	}
	ctx.WriteJSON(out)
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
	body.RemoteBase = body.Scheme + "://" + strings.Split(ctx.Request.RemoteAddr, ":")[0] + ":" + strconv.Itoa(body.Port)
	newb, _ := jsoniter.MarshalToString(body)

	var errs []string
	for _, route := range body.Routes {
		key := uriKey(route.Method, route.Uri)
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
