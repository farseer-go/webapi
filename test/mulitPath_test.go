package test

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"testing"
)

// http请求的路由，多路复用器
var serveMux = new(http.ServeMux)

func MulitPath(t *testing.T) {
	l, _ := net.Listen("tcp", "127.0.0.1:8080")
	_ = http.Serve(l, route())
}

type routeInfo struct {
	path    string
	handler http.HandlerFunc
}

var routePath = []routeInfo{
	{path: "^/index/\\d+$", handler: index}, // \d: 匹配数字
	{path: "^/home/\\w+$", handler: home},   // \w：匹配字母、数字、下划线
}

func route() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, route := range routePath {
			ok, err := regexp.Match(route.path, []byte(r.URL.Path))
			if err != nil {
				fmt.Println(err.Error())
			}
			if ok {
				route.handler(w, r)
				return
			}
		}
		_, _ = w.Write([]byte("404 not found"))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("index"))
}

func home(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("home"))
}
