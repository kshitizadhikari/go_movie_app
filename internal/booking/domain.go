package booking

import (
	"errors"
	"time"
)

var (
	ErrSeatAlreadyBooked    = errors.New("Seat has already been booked")
	ErrSeatConfirm          = errors.New("Failed to confirm booking")
	ErrSeatAlreadyConfirmed = errors.New("Seat has noot been booked yet")
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
	HoldSeat(b Booking) (*Booking, error)
	ConfirmSeat(booking_id string) (*Booking, error)
	ListBookings(movieId string) []Booking
}

func (b *Booking) ToResponse() Booking {
	return Booking{
		ID:        b.ID,
		UserID:    b.UserID,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		Status:    b.Status,
		ExpiresAt: b.ExpiresAt,
	}
}
