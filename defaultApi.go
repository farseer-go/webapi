package webapi

var defaultApi = applicationBuilder{
	area: "/",
}

type applicationBuilder struct {
	area string
}
