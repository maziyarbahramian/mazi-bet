package main

import (
	_ "mazi-bet/docs"
	"mazi-bet/server"
)

//	@Title			MAZI-BET
//	@version		1.0
//	@description	Interview Challenge

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						localhost:8080
// @BasePath					/
// @query.collection.format	multi
func main() {
	server.StartServer()
}
