package context

import "github.com/farseer-go/collections"

type HttpData struct {
	value collections.Dictionary[string, any]
}

// Get 获取值
func (receiver *HttpData) Get(key string) any {
	return receiver.value.GetValue(key)
}

// Set 设置值
func (receiver *HttpData) Set(key string, val any) {
	receiver.value.Add(key, val)
}
