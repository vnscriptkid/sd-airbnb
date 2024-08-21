package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User model
type UserRole string

const (
	UserRoleHost  UserRole = "Host"
	UserRoleGuest UserRole = "Guest"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255"`
	Email     string `gorm:"uniqueIndex;size:255"`
	Role      string `gorm:"size:50"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Hotel model
type Hotel struct {
	ID          uint `gorm:"primaryKey"`
	HostID      uint
	Name        string `gorm:"size:255"`
	Description string `gorm:"type:text"`
	Address     string `gorm:"size:500"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Room model
type RoomType string

const (
	RoomTypeSingle RoomType = "Single"
	RoomTypeDouble RoomType = "Double"
	RoomTypeSuite  RoomType = "Suite"
)

type Room struct {
	ID            uint `gorm:"primaryKey"`
	HotelID       uint
	Name          string `gorm:"size:255"`
	Description   string `gorm:"type:text"`
	Type          string `gorm:"size:50"` // e.g., Single, Double, Suite
	Quantity      uint   // Number of such rooms available
	PricePerNight float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type BookingStatus string

const (
	BookingStatusReserved BookingStatus = "Reserved"
	BookingStatusPending  BookingStatus = "Confirm"
	BookingStatusCanceled BookingStatus = "Canceled"
)

// Booking model
type Booking struct {
	ID          uint `gorm:"primaryKey"`
	RoomID      uint
	GuestID     uint
	StartDate   time.Time
	EndDate     time.Time
	TotalAmount float64
	Status      BookingStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RoomAvailability model
type RoomAvailability struct {
	ID        uint `gorm:"primaryKey"`
	RoomID    uint
	Date      time.Time
	RoomsLeft uint `gorm:"check:rooms_left >= 0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	// Database connection
	dsn := "host=localhost user=postgres dbname=postgres password=123456 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable Logger for debugging
	})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &Hotel{}, &Room{}, &RoomAvailability{}, &Booking{})

	// Setting up Gin
	r := gin.Default()

	// Example routes (not implemented)
	r.GET("/hotels", func(c *gin.Context) {
		// Logic to get hotels and rooms
		c.JSON(200, gin.H{"message": "List of hotels"})
	})

	// Example Create Room route
	r.POST("/rooms", func(c *gin.Context) {
		// Logic to create a room
		c.JSON(201, gin.H{"message": "Room created"})
	})

	// Example Reserve route
	r.POST("/bookings/reserve", func(c *gin.Context) {
		var booking Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// If user checks in on 1 and checks out on 3, the booking should be for 2 nights (1st and 2nd)
		// Need to check if the room is available for the selected dates
		// RoomID, Date, RoomsLeft
		// room-1, 1st,  X
		// room-1, 2nd,  Y
		// Check if X > 0 and Y > 0

		// Transaction to handle room availability and booking
		err := db.Transaction(func(tx *gorm.DB) error {
			var availabilities []RoomAvailability

			noNights := booking.EndDate.Sub(booking.StartDate).Hours() / 24
			lastNight := booking.EndDate.AddDate(0, 0, -1)

			if err := tx.Set("gorm:query_option", "FOR UPDATE").
				Where("room_id = ? AND date BETWEEN ? AND ?", booking.RoomID, booking.StartDate, lastNight).
				Find(&availabilities).Error; err != nil {
				return fmt.Errorf("room not available for the selected dates")
			}

			if len(availabilities) != int(noNights) {
				return fmt.Errorf("room not available for the selected dates")
			}

			for _, availability := range availabilities {
				if availability.RoomsLeft < 1 {
					return fmt.Errorf("no rooms available for the selected dates")
				}
			}

			// Reduce available rooms
			ids := make([]uint, len(availabilities))
			for i, availability := range availabilities {
				ids[i] = availability.ID
			}

			if err := tx.Model(&RoomAvailability{}).
				Where("id IN ?", ids).
				Update("rooms_left", gorm.Expr("rooms_left - ?", 1)).Error; err != nil {
				return err
			}

			// Create booking with "Reserved" status
			booking.Status = "Reserved"
			if err := tx.Create(&booking).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		} else {
			c.JSON(201, gin.H{"message": "Room reserved", "booking": booking})
		}
	})

	// Example Confirm route
	r.POST("/bookings/:id/confirm", func(c *gin.Context) {
		bookingID := c.Param("id")
		var booking Booking

		if err := db.First(&booking, bookingID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Booking not found"})
			return
		}

		if booking.Status != "Reserved" {
			c.JSON(400, gin.H{"error": "Booking is not in reserved state"})
			return
		}

		// Simulate payment confirmation with Stripe (mocked for this demo)
		// Assume payment is successful and update booking status to "Booked"
		booking.Status = "Booked"
		db.Save(&booking)

		c.JSON(200, gin.H{"message": "Booking confirmed", "booking": booking})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
