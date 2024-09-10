package webapi

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/webapi/context"
	"github.com/farseer-go/webapi/middleware"
	"github.com/farseer-go/webapi/websocket"
	"net"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"sync"
)

type serveMux struct {
	mu          sync.RWMutex                  // lock
	m           map[string]*context.HttpRoute // 完全地址匹配
	es          []*context.HttpRoute          // 前缀匹配
	hosts       bool                          // whether any patterns contain hostnames
	regexpRoute []*context.HttpRoute          // 正则匹配的路由
}

func (mux *serveMux) checkHandle(pattern string, handler http.Handler) {
	if pattern == "" {
		panic("webapi: invalid pattern")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("webapi: multiple registrations for " + pattern)
	}

	if handler == nil {
		panic("webapi: nil handler")
	}

	if mux.m == nil {
		mux.m = make(map[string]*context.HttpRoute)
	}
}

// HandleRoute 注册路由
func (mux *serveMux) HandleRoute(route *context.HttpRoute) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if route.Method.Count() == 0 {
		route.Method = collections.NewList("GET")
	}

	// 正则路径匹配（如果有）
	route.RouteRegexp = context.NewRouteRegexp(route.RouteUrl, context.RegexpTypePath, context.RouteRegexpOptions{
		StrictSlash:    false,
		UseEncodedPath: false,
	})

	// websocket
	if route.Schema == "ws" {
		route.Handler = websocket.SocketHandler(route)
	} else { // http
		route.Handler = HttpHandler(route)
	}

	// 检查规则
	mux.checkHandle(route.RouteUrl, route.Handler)

	// 完全地址匹配的路由
	mux.m[route.RouteUrl] = route

	if route.RouteUrl[len(route.RouteUrl)-1] == '/' {
		mux.es = appendSorted(mux.es, route)
	}

	if route.RouteUrl[0] != '/' {
		mux.hosts = true
	}

	// 如果使用了正则匹配
	if route.RouteRegexp.UseRegex {
		mux.regexpRoute = append(mux.regexpRoute, route)
		// 如果参数名称没有显示指定时，则自动按顺序指定
		if route.ParamNames.Count() == 0 {
			varNames := route.RouteRegexp.GetVarNames()
			for _, varName := range varNames {
				route.ParamNames.Add(varName)
			}
		}
	}
}

// HandleFunc registers the handler function for the given pattern.
func (mux *serveMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("webapi: nil handler")
	}
	mux.Handle(pattern, http.HandlerFunc(handler))
}

// Handle 添加一个请求的处理函数
// If a handler already exists for pattern, Handle panics.
func (mux *serveMux) Handle(pattern string, handler http.Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	// 检查规则
	mux.checkHandle(pattern, handler)

	route := &context.HttpRoute{RouteUrl: pattern, Handler: handler}
	mux.m[pattern] = route

	if pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, route)
	}

	if pattern[0] != '/' {
		mux.hosts = true
	}
}

// 当tcp收到发送请求时，会调用此方法
func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_ = r.ParseForm()
	h := mux.serveHTTPHandler(r)
	h.Handler.ServeHTTP(w, r)
}

func (mux *serveMux) serveHTTPHandler(r *http.Request) (route *context.HttpRoute) {
	// CONNECT requests are not canonicalized.
	if r.Method == "CONNECT" {
		// If r.URL.Path is /tree and its handler is not registered,
		// the /tree -> /tree/ redirect applies to CONNECT requests
		// but the path canonicalization does not.
		if u, ok := mux.redirectToPathSlash(r.URL.Host, r.URL.Path, r.URL); ok {
			return &context.HttpRoute{Handler: http.RedirectHandler(u.String(), http.StatusMovedPermanently)}
		}
		return mux.handler(r.Host, r.URL.Path, r)
	}

	// All other requests have any port stripped and path cleaned
	// before passing to mux.handler.
	host := stripHostPort(r.Host)
	path := cleanPath(r.URL.Path)

	// If the given path is /tree and its handler is not registered,
	// redirect for /tree/.
	if u, ok := mux.redirectToPathSlash(host, path, r.URL); ok {
		return &context.HttpRoute{Handler: http.RedirectHandler(u.String(), http.StatusMovedPermanently)}
	}

	if path != r.URL.Path {
		u := &url.URL{Path: path, RawQuery: r.URL.RawQuery}
		return &context.HttpRoute{Handler: http.RedirectHandler(u.String(), http.StatusMovedPermanently)}
	}

	return mux.handler(host, r.URL.Path, r)
}

// handler 找到pattern对应的Handler
func (mux *serveMux) handler(host, path string, r *http.Request) (route *context.HttpRoute) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		route = mux.match(host+path, r)
	}
	if route == nil {
		route = mux.match(path, r)
	}
	if route == nil {
		route = &context.HttpRoute{Handler: http.NotFoundHandler()}
		flog.Debugf("%s %s%s 404", r.Method, r.Host, path)
	}
	return
}

// 找到匹配的路由
func (mux *serveMux) match(path string, req *http.Request) *context.HttpRoute {
	// 完全匹配
	v, ok := mux.m[path]
	if ok {
		return v
	}

	// 正则匹配
	for _, r := range mux.regexpRoute {
		match, isMatch := r.RouteRegexp.Match(path)
		// 匹配到了
		if isMatch {
			for n, val := range match {
				req.Form.Add(n, val)
			}
			return r
		}
	}

	// 前缀匹配
	for _, r := range mux.es {
		if strings.HasPrefix(path, r.RouteUrl) {
			return r
		}
	}
	return nil
}

// 初始化中间件，统一初始化可以保证应用的路由设计跟中间件设置顺序没有要求
func (mux *serveMux) initMiddleware(lstMiddleware collections.List[context.IMiddleware]) {
	middleware.InitMiddleware(mux.m, lstMiddleware)
}

func (mux *serveMux) redirectToPathSlash(host, path string, u *url.URL) (*url.URL, bool) {
	mux.mu.RLock()
	shouldRedirect := mux.shouldRedirectRLocked(host, path)
	mux.mu.RUnlock()
	if !shouldRedirect {
		return u, false
	}
	path = path + "/"
	u = &url.URL{Path: path, RawQuery: u.RawQuery}
	return u, true
}

func (mux *serveMux) shouldRedirectRLocked(host, path string) bool {
	p := []string{path, host + path}

	for _, c := range p {
		if _, exist := mux.m[c]; exist {
			return false
		}
	}

	n := len(path)
	if n == 0 {
		return false
	}
	for _, c := range p {
		if _, exist := mux.m[c+"/"]; exist {
			return path[n-1] != '/'
		}
	}

	return false
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}

// stripHostPort returns h without any trailing ":<port>".
func stripHostPort(h string) string {
	// If no port on host, return unchanged
	if !strings.Contains(h, ":") {
		return h
	}
	host, _, err := net.SplitHostPort(h)
	if err != nil {
		return h // on error, return unchanged
	}
	return host
}

func appendSorted(es []*context.HttpRoute, e *context.HttpRoute) []*context.HttpRoute {
	n := len(es)
	i := sort.Search(n, func(i int) bool {
		return len(es[i].RouteUrl) < len(e.RouteUrl)
	})
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, &context.HttpRoute{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])                // Move shorter entries down
	es[i] = e
	return es
}

// GetHttpContext 在minimalApi模式下也可以获取到上下文
func GetHttpContext() *context.HttpContext {
	return context.RoutineHttpContext.Get()
}
