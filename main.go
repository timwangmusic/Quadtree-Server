package main

import "github.com/weihesdlegend/quadtree-server/server"

func main() {
	svr := server.Server{}
	svr.Init()

	svr.Run()
}
