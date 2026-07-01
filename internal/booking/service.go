package booking

type BookingService struct {
	store BookingStore
}

func NewBookingService(store BookingStore) *BookingService {
	return &BookingService{
		store: store,
	}
}

func (s *BookingService) HoldSeat(b Booking) (*Booking, error) {
	res, err := s.store.HoldSeat(b)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *BookingService) ConfirmSeat(booking_id string) (*Booking, error) {
	res, err := s.store.ConfirmSeat(booking_id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *BookingService) ListBookings(movieId string) []Booking {
	return s.store.ListBookings(movieId)
}
