package main

import (
	"fmt"
)

type Car struct {
	speed   int
	battery int
}

func NewCar(speed, battery int) *Car {
	return &Car{
		speed:   speed,
		battery: battery,
	}
}

func GetSpeed(car *Car) int {
	return car.speed
}

func GetBattery(car *Car) int {
	return car.battery
}

func ChargeCar(car *Car, minutes int) {
	chargeAdded := minutes / 2
	car.battery += chargeAdded
	if car.battery > 100 {
		car.battery = 100
	}
}

func TryFinish(car *Car, distance int) string {
	batteryNeeded := (distance) / 2
	if distance%2 != 0 {
		batteryNeeded = (distance / 2)
	} else {
		batteryNeeded = distance / 2
	}

	if car.battery < batteryNeeded {
		car.battery = 0
		return ""
	}

	car.battery -= batteryNeeded
	timeNeeded := float64(distance) / float64(car.speed)
	return fmt.Sprintf("%.2f", timeNeeded)
}
