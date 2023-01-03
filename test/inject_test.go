package test

type ITestInject interface {
	Call() string
}

type TestInject struct {
}

func (receiver TestInject) Call() string {
	return "ok"
}