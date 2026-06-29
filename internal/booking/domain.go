package booking

import (
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked = errors.New("Seat has already been booked")
)

type Booking struct {
	ID        string
	UserID    string
	MovieID   string
	SeatID    string
	Status    string
	ExpiresAt time.Time
}

type BookingStore interface {
	Book(b Booking) error
	ListBookings(movieId string) []Booking
}
