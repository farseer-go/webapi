package context

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"sync"
	"time"
)

const DEFEALT_TIME = 1800

func NewSessionMange() *SessionManager {
	SessionManager := &SessionManager{
		cookieName: "lz_cookie",
		storage:    newFromMemory(),
		maxAge:     1800,
	}

	go SessionManager.GC()
	return SessionManager
}

var SessionM *SessionManager = NewSessionMange()

type SessionManager struct {
	cookieName string
	storage    Provider
	maxAge     int64
	lock       sync.Mutex
}

type Provider interface {
	//初始化一个session，id根据需要生成后传入
	InitSession(sid string, maxAge int64) (Session, error)
	//根据sid，获得当前session
	SetSession(session Session) error
	//销毁session
	DestroySession(sid string) error
	//回收
	GCSession()
}

func (fm *FromMemory) DestroySession(sid string) error {
	if _, ok := fm.sessions[sid]; ok {
		delete(fm.sessions, sid)
		return nil
	}
	return nil
}

func (fm *FromMemory) SetSession(session Session) error {
	fm.sessions[session.GetId()] = session
	return nil
}

func (fm *FromMemory) GCSession() {
	sessions := fm.sessions
	if len(sessions) < 1 {
		return
	}
	for k, v := range sessions {
		t := (v.(*SessionFromMemory).lastAccessedTime.Unix()) + (v.(*SessionFromMemory).maxAge)
		if t < time.Now().Unix() {
		}
		delete(fm.sessions, k)
	}
}

func (m *SessionManager) GC() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.storage.GCSession()
	age2 := int(60 * time.Second)
	time.AfterFunc(time.Duration(age2), func() {
		m.GC()
	})
}

func newFromMemory() *FromMemory {
	return &FromMemory{
		sessions: make(map[string]Session, 0),
	}
}

type FromMemory struct {
	lock     sync.Mutex
	sessions map[string]Session
}

type CookieMemory struct {
	lock   sync.Mutex
	cookie Cookie
}

type SessionFromMemory struct {
	sid              string //唯一标识
	lock             sync.Mutex
	lastAccessedTime time.Time   //最后一次访问时间
	maxAge           int64       //超时时间
	data             map[any]any //主数据
}

type CookieFromMemory struct {
	lock sync.Mutex
	data map[any]any //主数据
}

type Session interface {
	Set(key, value any)
	Get(key any) any
	Remove(key any) error
	GetId() string
}

type Cookie interface {
	Set(key, value any)
	Get(key any) any
	Remove(key any) error
}

func newSessionFromMemory() *SessionFromMemory {
	return &SessionFromMemory{
		data:   make(map[any]any),
		maxAge: DEFEALT_TIME,
	}
}

func (si *SessionFromMemory) Set(key, value any) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.data[key] = value
}

func (si *SessionFromMemory) Get(key any) any {
	if value := si.data[key]; value != nil {
		return value
	}
	return nil
}

func (si *SessionFromMemory) Remove(key any) error {
	if value := si.data[key]; value != nil {
		delete(si.data, key)
	}
	return nil
}

func newCookieFromMemory() *CookieFromMemory {
	return &CookieFromMemory{
		data: make(map[any]any),
	}
}

func (si *CookieFromMemory) Set(key, value any) {
	si.lock.Lock()
	defer si.lock.Unlock()
	si.data[key] = value
}

func (si *CookieFromMemory) Get(key any) any {
	if value := si.data[key]; value != nil {
		return value
	}
	return nil
}

func (si *CookieFromMemory) Remove(key any) error {
	if value := si.data[key]; value != nil {
		delete(si.data, key)
	}
	return nil
}

func (si *SessionFromMemory) GetId() string {
	return si.sid
}

func (fm *FromMemory) InitSession(sid string, maxAge int64) (Session, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()
	newSession := newSessionFromMemory()
	newSession.sid = sid
	if maxAge != 0 {
		newSession.maxAge = maxAge
	}
	newSession.lastAccessedTime = time.Now()
	fm.sessions[sid] = newSession
	return newSession, nil
}

func InitCookie() (Cookie, error) {
	var cookie Cookie
	cookie = newCookieFromMemory()
	return cookie, nil
}

func getSId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
