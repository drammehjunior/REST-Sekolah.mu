package main

import (
	config2 "exampleclean.com/refactor/app/config"
	"exampleclean.com/refactor/app/di"
	"fmt"
	"log"
)

func main() {

	config, configErr := config2.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	fmt.Println(config)

	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}

}
