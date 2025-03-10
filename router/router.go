package router

import (
	"github.com/akilino/restaurant/handler"
	"github.com/akilino/restaurant/service"
	"github.com/gorilla/mux"
)

func SetupRouter(service *service.RentalService) *mux.Router {
	carHandler := handler.NewCarHandler(service)
	router := mux.NewRouter()

	router.HandleFunc("/cars", carHandler.AddCar).Methods("POST")
	router.HandleFunc("/cars", carHandler.ListCars).Methods("GET")
	router.HandleFunc("/rent/{id:[0-9]+}", carHandler.RentCar).Methods("POST")
	router.HandleFunc("/return/{id:[0-9]+}", carHandler.ReturnCar).Methods("POST")

	return router
}
