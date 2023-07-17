package context

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/cacheMemory"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/snowflake"
	"github.com/farseer-go/redis"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// session名称
const sessionId = "SessionId"

// 自动过期时间，单位：秒
const sessionTimeout = 1200

// 存储session每一项值
type nameValue struct {
	Name  string
	Value any
}

type HttpSession struct {
	id    string
	store cache.ICacheManage[nameValue]
}

// InitSession 初始化httpSession
func InitSession(w http.ResponseWriter, r *http.Request) *HttpSession {
	httpSession := &HttpSession{}

	c, _ := r.Cookie(sessionId)
	if c != nil {
		httpSession.id = c.Value
	}

	// 第一次请求
	if httpSession.id == "" {
		httpSession.id = strconv.FormatInt(snowflake.GenerateId(), 10)
		// 写入Cookies
		http.SetCookie(w, &http.Cookie{
			Name:     sessionId,
			Value:    httpSession.id,
			Path:     "/",
			HttpOnly: true,
		})
	}

	// 设置存储方式
	cacheId := "SessionId_" + httpSession.id
	if !container.IsRegister[cache.ICacheManage[nameValue]](cacheId) {
		// 设置定义的超时时间
		customSessionTimeout := sessionTimeout
		age := configure.GetInt("Webapi.Session.Age")
		if age > 0 {
			customSessionTimeout = age
		}
		// 根据配置，设置存储方式
		switch strings.ToLower(configure.GetString("Webapi.Session.Store")) {
		case "redis":
			httpSession.store = redis.SetProfiles[nameValue](cacheId, "Name", configure.GetString("Webapi.Session.StoreConfigName"), func(op *cache.Op) {
				op.SlidingExpiration(time.Duration(customSessionTimeout) * time.Second)
			})
		default:
			httpSession.store = cacheMemory.SetProfiles[nameValue](cacheId, "Name", func(op *cache.Op) {
				op.SlidingExpiration(time.Duration(customSessionTimeout) * time.Second)
			})
		}
	} else {
		httpSession.store = container.Resolve[cache.ICacheManage[nameValue]](cacheId)
	}
	return httpSession
}

// GetValue 获取Session
func (r *HttpSession) GetValue(name string) any {
	item, _ := r.store.GetItem(name)
	return item.Value
}

// SetValue 设置Session
func (r *HttpSession) SetValue(name string, val any) {
	r.store.SaveItem(nameValue{
		Name:  name,
		Value: val,
	})
}

// Remove 删除Session
func (r *HttpSession) Remove(name string) {
	r.store.Remove(name)
}

// Clear 清空Session
func (r *HttpSession) Clear() {
	r.store.Clear()
}
