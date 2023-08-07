package main

import (
	"dgateway/model"
	"fmt"
	goweb "github.com/Rehtt/Kit/web"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func gateway(ctx *goweb.Context) {
	key := fmt.Sprintf("%s|%s", strings.ToTitle(ctx.Request.Method), ctx.Request.URL.Path)
	value := rdb.Get(ctx, key).Val()
	if value == "" {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		ctx.Writer.Write([]byte("404 page not found"))
		return
	}
	var tmp model.Register
	if err := jsoniter.UnmarshalFromString(value, &tmp); err != nil {
		ctx.WriteJSON(&model.Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	p, err := NewProxy(tmp.Remote)
	if err != nil {
		ctx.WriteJSON(&model.Response{Error: err.Error()}, http.StatusBadGateway)
		return
	}
	p.ServeHTTP(ctx.Writer, ctx.Request)
}

func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}
	return httputil.NewSingleHostReverseProxy(url), nil
}
