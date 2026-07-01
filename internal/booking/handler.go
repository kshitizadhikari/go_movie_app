package booking

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *BookingService
}

func NewHandler(svc *BookingService) *Handler {
	return &Handler{svc: svc}
}

type CreateHoldRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	MovieID string `json:"movie_id" binding:"required"`
	SeatID  string `json:"seat_id" binding:"required"`
}

func (h *Handler) HoldSeat(c *gin.Context) {
	var req CreateHoldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b := Booking{
		UserID:    req.UserID,
		MovieID:   req.MovieID,
		SeatID:    req.SeatID,
		Status:    "hold",
		ExpiresAt: time.Now().Add(defaultHoldTTL),
	}

	res, err := h.svc.HoldSeat(b)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res.ToResponse()})
}

func (h *Handler) ConfirmSeat(c *gin.Context) {
	booking_id := c.Param("booking_id")

	res, err := h.svc.ConfirmSeat(booking_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": res.ToResponse()})

}

func (h *Handler) List(c *gin.Context) {
	id := c.Param("movie_id")
	res := h.svc.ListBookings(id)
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}
