package minimal

import "github.com/farseer-go/collections"

func Init() {
	lstRouteTable = collections.NewList[routeTable]()
}
