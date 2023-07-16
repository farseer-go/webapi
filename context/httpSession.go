package context

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/cacheMemory"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/snowflake"
	"net/http"
	"strconv"
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

// 初始化httpSession
func initSession(w http.ResponseWriter, r *http.Request) *HttpSession {
	c, _ := r.Cookie(sessionId)
	httpSession := &HttpSession{
		id: c.Value,
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
		httpSession.store = cacheMemory.SetProfiles[nameValue](cacheId, "Name", func(op *cache.Op) {
			op.SlidingExpiration(sessionTimeout * time.Second)
		})
	} else {
		httpSession.store = container.Resolve[cache.ICacheManage[nameValue]](cacheId)
	}
	return httpSession
}

// Get 获取Session
func (r *HttpSession) Get(name string) any {
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
