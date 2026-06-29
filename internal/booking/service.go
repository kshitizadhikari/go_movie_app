package booking

type BookingService struct {
	store BookingStore
}

func NewBookingService(store BookingStore) *BookingService {
	return &BookingService{
		store: store,
	}
}

func (s *BookingService) Book(b Booking) error {
	return s.store.Book(b)
}

func (s *BookingService) ListBookings(movieId string) []Booking {
	return s.store.ListBookings(movieId)
}
