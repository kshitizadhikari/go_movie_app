package booking

import "sync"

type ConcurrentStore struct {
	bookings map[string]Booking
	sync.RWMutex
}

func NewConcurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
}

func (s *ConcurrentStore) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.bookings[b.SeatID]; ok {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentStore) ListBookings(movieId string) []Booking {
	s.RLock()
	defer s.RUnlock()
	var bookings []Booking

	for _, v := range s.bookings {
		if movieId == v.MovieID {
			bookings = append(bookings, v)
		}
	}

	return bookings
}
