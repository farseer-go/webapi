package context

import (
	"sync"
	"time"
)

const sessionDefaultAge = 1800

var DefaultSession = newSessionMange()

func newSessionMange() *SessionManager {
	
	sessionManager := &SessionManager{
		cookieName: sessionId,
		storage:    newFromMemory(),
		maxAge:     sessionDefaultAge,
	}

	go sessionManager.GC()
	return sessionManager
}

type SessionManager struct {
	cookieName string
	storage    Provider
	maxAge     int64
	lock       sync.Mutex
}

func (m *SessionManager) GC() {
	tick := time.NewTicker(60 * time.Second)
	for range tick.C {
		m.lock.Lock()
		m.storage.GCSession()
		m.lock.Unlock()
	}
}

type Provider interface {
	// InitSession 初始化一个session，id根据需要生成后传入
	InitSession(sid string, maxAge int64) Session
	// SetSession 根据sid，获得当前session
	SetSession(session Session) error
	// DestroySession 销毁session
	DestroySession(sid string) error
	// GCSession 回收
	GCSession()
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

func (fm *FromMemory) InitSession(sid string, maxAge int64) Session {
	fm.lock.Lock()
	defer fm.lock.Unlock()
	newSession := newSessionFromMemory()
	newSession.sid = sid
	if maxAge != 0 {
		newSession.maxAge = maxAge
	}
	newSession.lastAccessedTime = time.Now()
	fm.sessions[sid] = newSession
	return newSession
}

type SessionFromMemory struct {
	sid              string //唯一标识
	lock             sync.Mutex
	lastAccessedTime time.Time   //最后一次访问时间
	maxAge           int64       //超时时间
	data             map[any]any //主数据
}

type Session interface {
	Set(key, value any)
	Get(key any) any
	Remove(key any) error
	GetId() string
}

func newSessionFromMemory() *SessionFromMemory {
	return &SessionFromMemory{
		data:   make(map[any]any),
		maxAge: sessionDefaultAge,
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

func (si *SessionFromMemory) GetId() string {
	return si.sid
}
