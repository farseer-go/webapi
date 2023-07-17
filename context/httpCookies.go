package context

import (
	"net/http"
	"time"
)

// 初始化Cookies
func initCookies(w http.ResponseWriter, r *http.Request) *HttpCookies {
	return &HttpCookies{
		w: w,
		r: r,
	}
}

type HttpCookies struct {
	w http.ResponseWriter
	r *http.Request
}

// Get 获取Cookie
func (r *HttpCookies) Get(name string) *http.Cookie {
	cookie, _ := r.r.Cookie(name)
	return cookie
}

// GetValue 获取Cookie
func (r *HttpCookies) GetValue(name string) string {
	cookie, _ := r.r.Cookie(name)
	if cookie == nil {
		return ""
	}
	return cookie.Value
}

// SetValue 设置Cookie
func (r *HttpCookies) SetValue(name string, val string) {
	http.SetCookie(r.w, &http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		HttpOnly: false,
	})
}

// SetSuretyValue 设置Cookie安全值，将不允许脚本读取该值（HttpOnly）
func (r *HttpCookies) SetSuretyValue(name string, val string) {
	http.SetCookie(r.w, &http.Cookie{
		Name:     name,
		Value:    val,
		Path:     "/",
		HttpOnly: true,
	})
}

// SetCookie 设置Cookie
func (r *HttpCookies) SetCookie(cookie *http.Cookie) {
	http.SetCookie(r.w, cookie)
}

// Remove 删除Cookie
func (r *HttpCookies) Remove(name string) {
	http.SetCookie(r.w, &http.Cookie{
		Name:     name,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
	})
}
