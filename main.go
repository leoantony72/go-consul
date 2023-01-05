package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
)

func main() {
	register()
	app := gin.Default()

	app.Run(":8090")

	app.GET("/", check)
}
func check(c *gin.Context) {
	c.JSON(200, gin.H{"message": "all good "})
}

func register() {
	config := consulapi.DefaultConfig()
	consul, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println(err)
	}
	serviceId := "go-server4321"
	port, _ := strconv.Atoi(getPort()[1:len(getPort())])
	fmt.Printf("port:%v \n", port)
	address := getHostname()
	fmt.Printf("address:%v \n", address)

	registeration := &consulapi.AgentServiceRegistration{
		ID:      serviceId,
		Name:    "testserver",
		Port:    port,
		Address: address,
		Check: &consulapi.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", address, port),
			Interval: "10s",
			Timeout:  "30s",
		},
	}

	regiErr := consul.Agent().ServiceRegister(registeration)
	if regiErr != nil {
		log.Panic(regiErr)
		log.Printf("Failed to register service: %s:%v ", address, port)
	} else {
		log.Printf("successfully register service: %s:%v", address, port)
	}
}

func getPort() (port string) {
	port = os.Getenv("PORT")
	if len(port) == 0 {
		port = "8090"
	}
	port = ":" + port
	return
}

func getHostname() (hostname string) {
	hostname, _ = os.Hostname()
	return
}
