package main

import (
	"fmt"
	"github.com/akilino/restaurant/router"
	"github.com/akilino/restaurant/service"
	"net/http"
)

func main() {
	service := service.NewRentalService()
	r := router.SetupRouter(service)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
