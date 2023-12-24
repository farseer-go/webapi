package context

type IHttpSession interface {
	// GetValue 获取Session
	GetValue(name string) any
	// SetValue 设置Session
	SetValue(name string, val any)
	// Remove 删除Session
	Remove(name string)
	// Clear 清空Session
	Clear()
}
