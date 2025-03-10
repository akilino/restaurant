package handler

import (
	"encoding/json"
	"github.com/akilino/restaurant/model"
	"github.com/akilino/restaurant/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CarHandler struct {
	Service *service.RentalService
}

func NewCarHandler(service *service.RentalService) *CarHandler {
	return &CarHandler{Service: service}
}

func (c *CarHandler) AddCar(w http.ResponseWriter, r *http.Request) {
	var car model.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c.Service.AddCar(&car)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(car)
}

func (c *CarHandler) ListCars(w http.ResponseWriter, r *http.Request) {
	availableCars := c.Service.ListAvailableCars()
	json.NewEncoder(w).Encode(availableCars)
}

func (c *CarHandler) RentCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.Service.RentCar(carID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Car rented successfully"))
}

func (c *CarHandler) ReturnCar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	carID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = c.Service.ReturnCar(carID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Car returned successfully"))
}
