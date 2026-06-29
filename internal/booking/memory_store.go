package booking

type MemoryStore struct {
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking) error {
	if _, ok := s.bookings[b.SeatID]; ok {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b
	return nil
}

func (s *MemoryStore) ListBookings(movieId string) []Booking {
	var bookings []Booking

	for _, v := range s.bookings {
		if movieId == v.MovieID {
			bookings = append(bookings, v)
		}
	}

	return bookings
}
