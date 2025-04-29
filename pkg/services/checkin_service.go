package services

import (
    "context"
    "strings"
    "strconv"
    "time"

    "github.com/iamsuteerth/skyfox-backend/pkg/models"
    "github.com/iamsuteerth/skyfox-backend/pkg/repositories"
    "github.com/iamsuteerth/skyfox-backend/pkg/utils"
    "github.com/rs/zerolog/log"
)

type CheckInService interface {
	FindConfirmedBookings(ctx context.Context) ([]*models.Booking, error)
    MarkBookingsCheckedIn(ctx context.Context, bookingIDs []int) (checkedIn []int, alreadyDone []int, invalid []int, err error)
    MarkBookingCheckedIn(ctx context.Context, bookingID int) (bool, error)
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
            alreadyDone = append(alreadyDone, id) // e.g. race
        }
    }
    return checkedIn, alreadyDone, invalid, nil
}

func (s *checkInService) MarkBookingCheckedIn(ctx context.Context, bookingID int) (bool, error) {
    booking, err := s.bookingRepo.GetBookingById(ctx, bookingID)
    if err != nil {
        return false, err
    }
    if booking == nil || booking.Status != "Confirmed" {
        return false, nil
    }
    show, err := s.showRepo.FindById(ctx, booking.ShowId)
    if err != nil || show == nil {
        return false, nil
    }
    endTimeParsed, parseErr := parseShowEndTime(show.Date, show.Slot.EndTime)
    if parseErr != nil {
        log.Warn().Int("show_id", show.Id).Msg("Could not parse show end time")
        return false, nil
    }
    if time.Now().After(endTimeParsed) {
        return false, nil 
    }
    return s.bookingRepo.MarkBookingCheckedIn(ctx, bookingID)
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
    return time.Date(showDate.Year(), showDate.Month(), showDate.Day(), hour, minute, 0, 0, showDate.Location()), nil
}