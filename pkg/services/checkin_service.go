package services

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
)

type CheckInService interface {
	FindConfirmedBookings(ctx context.Context) ([]*models.Booking, error)
	MarkBookingsCheckedIn(ctx context.Context, bookingIDs []int) (checkedIn []int, alreadyDone []int, invalid []int, err error)
}

type checkInService struct {
	bookingRepo repositories.BookingRepository
	showRepo    repositories.ShowRepository
}

func NewCheckInService(
	bookingRepo repositories.BookingRepository,
	showRepo repositories.ShowRepository,
) CheckInService {
	return &checkInService{
		bookingRepo: bookingRepo,
		showRepo:    showRepo,
	}
}

func (s *checkInService) FindConfirmedBookings(ctx context.Context) ([]*models.Booking, error) {
	bookings, err := s.bookingRepo.FindConfirmedBookings(ctx)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (s *checkInService) MarkBookingsCheckedIn(ctx context.Context, bookingIDs []int) ([]int, []int, []int, error) {
	if len(bookingIDs) == 0 {
		return nil, nil, nil, nil
	}

	bookings, err := s.bookingRepo.FindBookingsByIds(ctx, bookingIDs)
	if err != nil {
		return nil, nil, nil, err
	}

	now := time.Now()
	checkedIn := make([]int, 0)
	alreadyDone := make([]int, 0)
	invalid := make([]int, 0)

	bookingMap := make(map[int]*models.Booking, len(bookings))
	for _, b := range bookings {
		bookingMap[b.Id] = b
	}

	for _, id := range bookingIDs {
		b := bookingMap[id]
		if b == nil {
			invalid = append(invalid, id)
			continue
		}
		if b.Status == "CheckedIn" {
			alreadyDone = append(alreadyDone, id)
			continue
		}
		if b.Status != "Confirmed" {
			invalid = append(invalid, id)
			continue
		}
		show, err := s.showRepo.FindById(ctx, b.ShowId)
		if err != nil || show == nil {
			invalid = append(invalid, id)
			continue
		}
		if !isWithinCheckInWindow(now, show.Date, show.Slot.StartTime) {
			invalid = append(invalid, id)
			continue
		}
		endTimeParsed, parseErr := parseShowEndTime(show.Date, show.Slot.EndTime)
		if parseErr != nil || now.After(endTimeParsed) {
			invalid = append(invalid, id)
			continue
		}
		ok, err := s.bookingRepo.MarkBookingCheckedIn(ctx, id)
		if err != nil {
			invalid = append(invalid, id)
			continue
		}
		if ok {
			checkedIn = append(checkedIn, id)
		} else {
			alreadyDone = append(alreadyDone, id)
		}
	}
	return checkedIn, alreadyDone, invalid, nil
}

func isWithinCheckInWindow(now time.Time, showDate time.Time, startTime string) bool {
	startTimeParsed, err := parseShowStartTime(showDate, startTime)
	if err != nil {
		return false
	}
	checkInStartTime := startTimeParsed.Add(-1 * time.Hour)
	return now.After(checkInStartTime) || now.Equal(checkInStartTime)
}

func parseShowStartTime(showDate time.Time, timeStr string) (time.Time, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return time.Time{}, utils.NewInternalServerError("INVALID_SLOT_TIME", "Slot time format invalid", nil)
	}
	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return time.Time{}, utils.NewInternalServerError("INVALID_SLOT_TIME", "Slot time values invalid", nil)
	}
	return time.Date(showDate.Year(), showDate.Month(), showDate.Day(), hour, minute, 0, 0, showDate.Location()), nil
}

func parseShowEndTime(showDate time.Time, slotEndTime string) (time.Time, error) {
	parts := strings.Split(slotEndTime, ":")
	if len(parts) < 2 {
		return time.Time{}, utils.NewInternalServerError("INVALID_SLOT_TIME", "Slot end time format invalid", nil)
	}
	hour, err1 := strconv.Atoi(parts[0])
	minute, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return time.Time{}, utils.NewInternalServerError("INVALID_SLOT_TIME", "Slot end time values invalid", nil)
	}
	nextDay := hour >= 0 && hour < 4
	date := showDate
	if nextDay {
		date = date.AddDate(0, 0, 1)
	}
	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, 0, 0, date.Location()), nil
}
