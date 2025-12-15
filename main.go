package main

import (
	"github.com/u00io/gomisc/logger"
	"github.com/u00io/localports/forms/mainform"
	"github.com/u00io/localports/localstorage"
)

func main() {
	localstorage.Init("localports")
	logger.Init(localstorage.Path() + "/logs")
	mainform.Run()
}
