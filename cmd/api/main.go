package main

import (
	config2 "exampleclean.com/refactor/app/config"
	"exampleclean.com/refactor/app/di"
	"github.com/spf13/viper"
	"log"
)

func main() {

	//viper.AddConfigPath("./cmd")
	//viper.SetConfigName("conf")
	//viper.SetConfigType("env")
	//
	//viper.AutomaticEnv()
	//
	//err := viper.ReadInConfig()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

	//fmt.Println(viper.Get("DB_DRIVER")

	viper.AddConfigPath("./cmd")
	viper.SetConfigName("conf")
	viper.SetConfigType("env")

	config, configErr := config2.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}

}
