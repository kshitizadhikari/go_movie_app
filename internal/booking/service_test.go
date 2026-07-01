package booking

import (
	"sync"
	"sync/atomic"
	"testing"

	"go_movie_booking_app/internal/adapters/redis"

	"github.com/google/uuid"
)

func TestBookingConcurrency(t *testing.T) {
	store := NewRedisStore(redis.NewClient("localhost:6379"))
	svc := NewBookingService(store)

	const numGoRoutines = 100_000

	var (
		successes atomic.Int64
		failures  atomic.Int64
		wg        sync.WaitGroup
	)

	wg.Add(numGoRoutines)

	for i := range numGoRoutines {
		go func(userNum int) {
			defer wg.Done()
			_, err := svc.HoldSeat(Booking{
				MovieID: "movie1",
				SeatID:  "A1",
				UserID:  uuid.New().String(),
			})

			if err != nil {
				failures.Add(1)
			} else {
				successes.Add(1)
			}

		}(i)
	}

	wg.Wait()
	if got := successes.Load(); got != 1 {
		t.Errorf("expected exactly 1 success, got %d", got)
	}

	if got := failures.Load(); got != int64(numGoRoutines-1) {
		t.Errorf("expected %d failures, got %d", numGoRoutines-1, got)
	}
}
