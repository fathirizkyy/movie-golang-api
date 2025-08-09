package main

import (
	"backend/confiq"
	"backend/routes"
)

func main() {
	confiq.ConnectDB()
	
	
	r:=routes.SetupRouter()
	r.Run(":8000")

}