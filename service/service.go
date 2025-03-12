package service

import (
	"errors"
	"fmt"
	"github.com/akilino/restaurant/model"
	"sync"
)

type RentalService struct {
	cars      map[int]*model.Car
	mu        sync.Mutex
	rentReqCh chan RentRequest
}

type RentRequest struct {
	CarID      int
	ResponseCh chan error
}

func NewRentalService() *RentalService {
	rs := &RentalService{
		cars:      make(map[int]*model.Car),
		rentReqCh: make(chan RentRequest, 10),
	}

	// Start a goroutine to process rental requests
	go rs.processRentals()

	return rs
}

func (rs *RentalService) Mutex() *sync.Mutex {
	return &rs.mu
}

func (rs *RentalService) processRentals() {
	for rentRequest := range rs.rentReqCh {
		rs.mu.Lock()
		car, exists := rs.cars[rentRequest.CarID]
		if !exists {
			rentRequest.ResponseCh <- fmt.Errorf("car %d does not exist", rentRequest.CarID)
		} else if car.IsRented {
			rentRequest.ResponseCh <- fmt.Errorf("car %d already rented", rentRequest.CarID)
		} else {
			car.IsRented = true
			rentRequest.ResponseCh <- nil
		}

		rs.mu.Unlock()
	}
}

func (rs *RentalService) GetCar(carID int) (*model.Car, bool) {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	car, exists := rs.cars[carID]
	return car, exists
}

func (rs *RentalService) AddCar(car *model.Car) {
	rs.mu.Lock()
	rs.cars[car.ID] = car
	defer rs.mu.Unlock()
}

func (rs *RentalService) RentCar(carID int) error {
	responseCh := make(chan error)
	rs.rentReqCh <- RentRequest{carID, responseCh}

	return <-responseCh
}

func (rs *RentalService) ReturnCar(carID int) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	car, exist := rs.cars[carID]
	if !exist {
		return errors.New("car not found")
	}

	if !car.IsRented {
		return errors.New("car is not rented")
	}
	car.IsRented = false
	return nil
}

func (rs *RentalService) ListAvailableCars() []model.Car {
	rs.mu.Lock()
	defer rs.mu.Unlock()

	var availableCars []model.Car
	for _, car := range rs.cars {
		if !car.IsRented {
			availableCars = append(availableCars, *car)
		}
	}
	return availableCars
}
