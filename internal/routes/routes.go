package routes

import (
	"go_movie_booking_app/internal/booking"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, bookingHandler *booking.Handler) {
	api := router.Group("api")
	v1 := api.Group("v1")

	bookings := v1.Group("/movies")
	{
		bookings.POST("/hold", bookingHandler.HoldSeat)
		bookings.PATCH("/confirm/:booking_id", bookingHandler.ConfirmSeat)
		bookings.GET("/:movie_id", bookingHandler.List)
	}
}
