package main

import (
	"github.com/Rehtt/DGateway/model"
	goweb "github.com/Rehtt/Kit/web"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

func gateway(ctx *goweb.Context) {
	key := uriKey(ctx.Request.Method, ctx.Request.URL.Path, EQ)
	value := rdb.Get(ctx, key).Val()
	if value == "" {
		key = uriKey(ctx.Request.Method, "*", Match)
		var ok bool
		for _, v := range rdb.Keys(ctx, key).Val() {
			if s := strings.Split(v, "|"); len(s) == 4 {
				if ok, _ = regexp.MatchString(s[3], ctx.Request.URL.Path); ok {
					value = rdb.Get(ctx, v).Val()
					break
				}
			}
		}
		if !ok || value == "" {
			ctx.Writer.WriteHeader(http.StatusNotFound)
			ctx.Writer.Write([]byte("404 page not found"))
			return
		}
	}
	var tmp model.Register
	if err := jsoniter.UnmarshalFromString(value, &tmp); err != nil {
		ctx.WriteJSON(&model.Response{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	p, err := NewProxy(tmp.RemoteBase)
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
