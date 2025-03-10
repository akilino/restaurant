package main

import (
	"bytes"
	"encoding/json"
	"github.com/akilino/restaurant/model"
	"github.com/akilino/restaurant/router"
	"github.com/akilino/restaurant/service"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCarRentalAPI(t *testing.T) {
	service := service.NewRentalService()
	r := router.SetupRouter(service)

	// Add Car Test
	car := map[string]interface{}{"id": 1, "make": "Toyota", "model": "Corolla"}
	carJSON, _ := json.Marshal(car)

	req := httptest.NewRequest("POST", "/cars", bytes.NewBuffer(carJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.Code)
	}

	// Rent Car Test
	req = httptest.NewRequest("POST", "/rent/1", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Return Car Test
	req = httptest.NewRequest("POST", "/return/1", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}
}

func TestRentalService(t *testing.T) {
	rs := service.NewRentalService()

	// Add cars
	rs.AddCar(&model.Car{ID: 1, Make: "Toyota", Model: "Corolla", IsRented: false})
	rs.AddCar(&model.Car{ID: 2, Make: "Honda", Model: "Civic", IsRented: false})

	// Test initial car availability
	availableCars := rs.ListAvailableCars()
	if len(availableCars) != 2 {
		t.Errorf("Expected 2 cars, got %d", len(availableCars))
	}

	// Test renting a car
	err := rs.RentCar(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Checks if car is marked as rented
	availableCars = rs.ListAvailableCars()
	if len(availableCars) != 1 {
		t.Errorf("Expected 1 cars, got %d", len(availableCars))
	}

	// Test renting an already rented car
	err = rs.RentCar(2)
	if err == nil {
		t.Errorf("Expected error when renting an already rented car, got nil")
	}

	err = rs.ReturnCar(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check if the car is available again
	availableCars = rs.ListAvailableCars()
	if len(availableCars) != 2 {
		t.Errorf("Expected 2 cars, got %d", len(availableCars))
	}
}
