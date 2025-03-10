package service

import (
	"errors"
	"github.com/akilino/restaurant/model"
	"sync"
)

type RentalService struct {
	cars      map[int]*model.Car
	mu        sync.Mutex
	rentReqCh chan int
}

func NewRentalService() *RentalService {
	return &RentalService{
		cars:      make(map[int]*model.Car),
		rentReqCh: make(chan int, 10),
	}
}

func (rs *RentalService) AddCar(car *model.Car) {
	rs.mu.Lock()
	rs.cars[car.ID] = car
	defer rs.mu.Unlock()
}

func (rs *RentalService) RentCar(carID int) error {
	rs.rentReqCh <- carID

	rs.mu.Lock()
	defer rs.mu.Unlock()

	car, exist := rs.cars[carID]
	if !exist {
		return errors.New("car not found")
	}
	if car.IsRented {
		return errors.New("car already rented")
	}
	car.IsRented = true
	return nil
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
