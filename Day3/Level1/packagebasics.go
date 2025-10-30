package main

import (
	"fmt" // prewritten struct interface function methods
	// precoded goroutine exist here
	"log"
	"os"

	"go-micro.dev/v4/logger"
	"go.uber.org/zap"
)

func main() {

	fmt.Println("smiling day")
	log.Println("main started")
	name := "venkatesh"
	log.Printf("applicatio  is running %s", name)
	fmt.Println("os args", os.Args)
	fmt.Println(os.Getwd())

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("error in code", err)
		return
	}
	fmt.Println(dir)

	logger, _ = zap.NewProduction()

	logger.Info("server started")

}
