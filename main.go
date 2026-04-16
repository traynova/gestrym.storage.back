package main

import (
	"flag"
	app "gestrym-storage/src"
)

//	@title			Gestrym API Storage
//	@version		1.0
//	@description	API para el manejo de entrenamientos.
//	@BasePath		/gestrym-storage

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key

//	@securityDefinitions.basic	BasicAuth

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	isLocalDeploy := flag.Bool("local", false, "'--local=true' para desplegar en ambiente local")
	flag.Parse()
	app.Run(*isLocalDeploy)
}
