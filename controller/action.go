package controller

type Action struct {
	Name   string // Action名称
	Method string // POST/GET/PUT/DELETE
	Params string // 函数的入参名称
}
