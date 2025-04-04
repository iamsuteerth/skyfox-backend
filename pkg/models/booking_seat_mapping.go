	package models

	type BookingSeatMapping struct {
		Id         int    `json:"id" gorm:"primaryKey"`
		BookingId  int    `gorm:"foreignKey:id"`
		SeatNumber string `json:"SeatNumber"`
	}

	func NewBookingSeatMapping(bookingId int, seatnumber string, seatType string) BookingSeatMapping {
		return BookingSeatMapping{
			BookingId:  bookingId,
			SeatNumber: seatnumber,
		}
	}

	func (BookingSeatMapping) TableName() string {
		return "booking_seat_mapping"
	}
