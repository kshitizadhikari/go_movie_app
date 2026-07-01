package main

import (
	"go_movie_booking_app/internal/adapters/redis"
	"go_movie_booking_app/internal/booking"
	"go_movie_booking_app/internal/routes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	store := booking.NewRedisStore(redis.NewClient("redis:6379"))
	svc := booking.NewBookingService(store)
	bookingHandler := booking.NewHandler(svc)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	routes.RegisterRoutes(router, bookingHandler)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
		return
	}

	log.Printf("🚀 Server running on http://localhost:8080")
}
