package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const defaultHoldTTL = 2 * time.Minute

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}

func (s *RedisStore) ConfirmSeat(booking_id string) (*Booking, error) {
	ctx := context.Background()
	booking_key := fmt.Sprintf("booking:%s", booking_id)

	data, err := s.rdb.Get(ctx, booking_key).Bytes()
	if err != nil {
		return nil, ErrSeatConfirm
	}

	var booking Booking
	if err := json.Unmarshal(data, &booking); err != nil {
		return nil, err
	}

	if booking.Status == "confirmed" {
		return nil, ErrSeatAlreadyConfirmed
	}

	booking.Status = "confirmed"

	data, err = json.Marshal(booking)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, booking_key, data, 0).Err(); err != nil {
		return nil, err
	}

	if err := s.rdb.Persist(ctx, booking_key).Err(); err != nil {
		return nil, err
	}

	seat_key := fmt.Sprintf("seat:%s:%s", booking.MovieID, booking.SeatID)
	if err := s.rdb.Persist(ctx, seat_key).Err(); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *RedisStore) ListBookings(movieID string) []Booking {
	ctx := context.Background()

	movieBookingsKey := fmt.Sprintf("movie:%s:bookings", movieID)

	ids, err := s.rdb.SMembers(ctx, movieBookingsKey).Result()
	if err != nil {
		return nil
	}

	bookings := make([]Booking, 0, len(ids))

	for _, id := range ids {
		data, err := s.rdb.Get(ctx, fmt.Sprintf("booking:%s", id)).Result()
		if err != nil {
			continue
		}

		var booking Booking
		if err := json.Unmarshal([]byte(data), &booking); err != nil {
			continue
		}

		bookings = append(bookings, booking)
	}

	return bookings
}

func (s *RedisStore) HoldSeat(b Booking) (*Booking, error) {
	ctx := context.Background()

	id := uuid.NewString()
	now := time.Now()

	booking := &Booking{
		ID:        id,
		UserID:    b.UserID,
		MovieID:   b.MovieID,
		SeatID:    b.SeatID,
		Status:    "held",
		ExpiresAt: now.Add(defaultHoldTTL),
	}

	seatKey := fmt.Sprintf("seat:%s:%s", booking.MovieID, booking.SeatID)

	res := s.rdb.SetArgs(ctx, seatKey, booking.ID, redis.SetArgs{
		Mode: "NX",
		TTL:  defaultHoldTTL,
	})

	if err := res.Err(); err != nil {
		return nil, ErrSeatAlreadyBooked
	}

	if res.Val() != "OK" {
		return nil, ErrSeatAlreadyBooked
	}

	if err := s.rdb.Set(ctx, sessionKey(id), seatKey, defaultHoldTTL).Err(); err != nil {
		return nil, err
	}

	bookingKey := fmt.Sprintf("booking:%s", id)

	data, err := json.Marshal(booking)
	if err != nil {
		return nil, err
	}

	if err := s.rdb.Set(ctx, bookingKey, data, defaultHoldTTL).Err(); err != nil {
		return nil, err
	}

	movieBookingsKey := fmt.Sprintf("movie:%s:bookings", booking.MovieID)

	if err := s.rdb.SAdd(ctx, movieBookingsKey, booking.ID).Err(); err != nil {
		return nil, err
	}

	if err := s.rdb.Expire(ctx, movieBookingsKey, defaultHoldTTL).Err(); err != nil {
		return nil, err
	}

	return booking, nil
}
