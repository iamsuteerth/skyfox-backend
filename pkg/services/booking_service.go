package services

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/govalues/decimal"
	"github.com/iamsuteerth/skyfox-backend/pkg/constants"
	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	movieservice "github.com/iamsuteerth/skyfox-backend/pkg/movie-service"
	"github.com/iamsuteerth/skyfox-backend/pkg/repositories"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/rs/zerolog/log"

	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
)

type BookingService interface {
	GetSeatMapForShow(ctx context.Context, showID int) ([]models.SeatMapEntry, error)
	GetBookingById(ctx context.Context, bookindID int) (*models.Booking, error)
	GenerateQRCode(ctx context.Context, bookingID int) (string, error)
	GeneratePDF(ctx context.Context, bookingID int) (string, error)
}

type bookingService struct {
	showRepo                repositories.ShowRepository
	bookingRepo             repositories.BookingRepository
	bookingSeatMappingRepo  repositories.BookingSeatMappingRepository
	slotRepo                repositories.SlotRepository
	adminBookedCustomerRepo repositories.AdminBookedCustomerRepository
	skyCustomerRepo         repositories.SkyCustomerRepository
	movieService            movieservice.MovieService
}

func NewBookingService(
	showRepo repositories.ShowRepository,
	bookingRepo repositories.BookingRepository,
	bookingSeatMappingRepo repositories.BookingSeatMappingRepository,
	slotRepo repositories.SlotRepository,
	adminBookedCustomerRepo repositories.AdminBookedCustomerRepository,
	skyCustomerRepo repositories.SkyCustomerRepository,
	movieService movieservice.MovieService,
) BookingService {
	return &bookingService{
		showRepo:                showRepo,
		bookingRepo:             bookingRepo,
		bookingSeatMappingRepo:  bookingSeatMappingRepo,
		slotRepo:                slotRepo,
		adminBookedCustomerRepo: adminBookedCustomerRepo,
		movieService:            movieService,
		skyCustomerRepo:         skyCustomerRepo,
	}
}

func (s *bookingService) GetSeatMapForShow(ctx context.Context, showID int) ([]models.SeatMapEntry, error) {
	show, err := s.showRepo.FindById(ctx, showID)
	if err != nil {
		return nil, err
	}

	seatMap, err := s.showRepo.GetSeatMapForShow(ctx, showID)
	if err != nil {
		return nil, err
	}

	deluxeOffset, _ := decimal.NewFromFloat64(constants.DELUXE_OFFSET)

	for i := range seatMap {
		if seatMap[i].SeatType == "Deluxe" {
			seatMap[i].Price, _ = show.Cost.Add(deluxeOffset)
		} else {
			seatMap[i].Price = show.Cost
		}
	}

	return seatMap, nil
}

func (s *bookingService) GetBookingById(ctx context.Context, bookingID int) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetBookingById(ctx, bookingID)
	if err != nil {
		log.Error().Err(err).Msg("Booking not found for the given id")
		return nil, err
	}
	return booking, nil
}

func (s *bookingService) GenerateQRCode(ctx context.Context, bookingID int) (string, error) {
	ticketData, err := s.getTicketData(ctx, bookingID)
	if err != nil {
		return "", err
	}

	showTimeFormatted := ticketData.ShowTime
	if t, err := time.Parse("15:04:05.000000", ticketData.ShowTime); err == nil {
		showTimeFormatted = t.Format("03:04 PM")
	}

	qrContent := fmt.Sprintf("SKYFOX CINEMA BOOKING #%d\nCustomer Name: %s\nContact Number: %s\nShow: %s\nDate: %s\nTime: %s\nSeats: %v\nAmount Paid: %.2f\nPayment Type : %s\nStatus : %s",
		ticketData.BookingID,
		ticketData.CustomerName,
		ticketData.ContactNumber,
		ticketData.ShowName,
		ticketData.ShowDate,
		showTimeFormatted,
		ticketData.SeatNumbers,
		ticketData.AmountPaid,
		ticketData.PaymentType,
		ticketData.Status,
	)

	qr, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return "", utils.NewInternalServerError("QR_GENERATION_FAILED", "Failed to generate QR code", err)
	}

	base64QR := base64.StdEncoding.EncodeToString(qr)

	return base64QR, nil
}

func (s *bookingService) GeneratePDF(ctx context.Context, bookingID int) (string, error) {
	ticketData, err := s.getTicketData(ctx, bookingID)
	if err != nil {
		return "", err
	}

	showTimeFormatted := ticketData.ShowTime
	if t, err := time.Parse("15:04:05.000000", ticketData.ShowTime); err == nil {
		showTimeFormatted = t.Format("03:04 PM")
	}

	qrCodeContent := fmt.Sprintf("SKYFOX BOOKING #%d\nCustomer: %s\nContact Number: %s\nShow: %s\nDate: %s\nTime: %s\nSeats: %v\nAmount Paid: %.2f\nPayment Type: %s\nStatus: %s",
		ticketData.BookingID,
		ticketData.CustomerName,
		ticketData.ContactNumber,
		ticketData.ShowName,
		ticketData.ShowDate,
		showTimeFormatted,
		ticketData.SeatNumbers,
		ticketData.AmountPaid,
		ticketData.PaymentType,
		ticketData.Status,
	)

	qrBytes, err := qrcode.Encode(qrCodeContent, qrcode.Medium, 256)
	if err != nil {
		return "", utils.NewInternalServerError("QR_GENERATION_FAILED", "Failed to generate QR code for PDF", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetAutoPageBreak(false, 0)

	addFontsWithFallback(pdf)

	qrReader := bytes.NewReader(qrBytes)
	qrImageOptions := gofpdf.ImageOptions{
		ImageType: "PNG",
	}

	pdf.RegisterImageOptionsReader("qrcode", qrImageOptions, qrReader)
	pdf.AddPage()

	// Header with brand color
	pdf.SetFillColor(224, 75, 0) // Primary color E04B00 from theme
	pdf.Rect(0, 0, 210, 40, "F")

	addLogoWithFallback(pdf, 10, 5, 30)

	pdf.SetFont("Poppins", "B", 24)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(45, 15)
	pdf.Cell(0, 10, "SKYFOX CINEMAS")

	pdf.SetTextColor(22, 26, 30) // text.primary from theme

	pdf.SetY(50)
	pdf.SetFont("Poppins", "B", 18)
	pdf.Cell(0, 10, "MOVIE TICKET")
	pdf.Ln(15)

	pdf.SetTextColor(224, 75, 0) // primary color from theme
	pdf.SetFont("Poppins", "B", 16)

	movieName := ticketData.ShowName
	if pdf.GetStringWidth(movieName) > 180 {
		if len(movieName) > 60 {
			movieName = movieName[:57] + "..."
		}
	}
	pdf.Cell(0, 10, movieName)
	pdf.Ln(15)

	pdf.SetTextColor(64, 67, 72) // text.secondary from theme

	seatText := strings.Join(ticketData.SeatNumbers, ", ")

	maxSeatWidth := 60.0
	if pdf.GetStringWidth(seatText) > maxSeatWidth && len(ticketData.SeatNumbers) > 5 {
		seatGroups := []string{}
		currentGroup := ""
		currentWidth := 0.0

		for i, seat := range ticketData.SeatNumbers {
			seatWithComma := seat
			if i > 0 {
				seatWithComma = ", " + seat
			}

			if currentWidth+pdf.GetStringWidth(seatWithComma) > maxSeatWidth && currentGroup != "" {
				seatGroups = append(seatGroups, currentGroup)
				currentGroup = seat
				currentWidth = pdf.GetStringWidth(seat)
			} else {
				if currentGroup != "" {
					currentGroup += ", " + seat
					currentWidth += pdf.GetStringWidth(seatWithComma)
				} else {
					currentGroup = seat
					currentWidth = pdf.GetStringWidth(seat)
				}
			}
		}

		if currentGroup != "" {
			seatGroups = append(seatGroups, currentGroup)
		}
		seatText = strings.Join(seatGroups, "\n")
	}

	// Layout columns
	pdf.SetFont("Poppins", "", 11)
	leftCol := 10.0
	rightCol := 110.0
	lineHeight := 8.0

	// Calculate extra height needed for seats
	extraHeight := 0.0
	if strings.Contains(seatText, "\n") {
		extraHeight = float64(strings.Count(seatText, "\n")) * lineHeight
	}

	// Left column with light background - adjusted height
	columnTop := pdf.GetY()
	columnHeight := 5 * lineHeight
	pdf.SetFillColor(240, 240, 245) // background.secondary from theme
	pdf.Rect(leftCol-5, columnTop-3, 95, columnHeight+6+extraHeight, "F")

	// Left column details
	pdf.SetXY(leftCol, columnTop)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30) // text.primary
	pdf.Cell(30, lineHeight, "Date:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72) // text.secondary
	pdf.Cell(60, lineHeight, ticketData.ShowDate)
	pdf.Ln(lineHeight)

	pdf.SetX(leftCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Time:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, showTimeFormatted)
	pdf.Ln(lineHeight)

	pdf.SetX(leftCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Seats:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)

	// Handle multiline seat text
	if strings.Contains(seatText, "\n") {
		initialY := pdf.GetY()
		pdf.MultiCell(60, lineHeight, seatText, "", "", false)
		pdf.SetY(initialY + float64(strings.Count(seatText, "\n")+1)*lineHeight)
	} else {
		pdf.Cell(60, lineHeight, seatText)
		pdf.Ln(lineHeight)
	}

	pdf.SetX(leftCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Amount:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, fmt.Sprintf("₹%.2f", ticketData.AmountPaid))
	pdf.Ln(lineHeight)

	pdf.SetX(leftCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Payment:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, ticketData.PaymentType)
	pdf.Ln(lineHeight)

	// Right column with brand secondary color - adjusted height to match left column
	pdf.SetFillColor(255, 177, 153) // secondary color
	pdf.Rect(rightCol-5, columnTop-3, 95, columnHeight+6+extraHeight, "F")

	// Right column details - Make customer name responsive
	pdf.SetXY(rightCol, columnTop)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Customer:")

	customerNameDisplay := ticketData.CustomerName
	if pdf.GetStringWidth(customerNameDisplay) > 60 {
		pdf.SetFont("Poppins", "", 9)
		if pdf.GetStringWidth(customerNameDisplay) > 65 {
			if len(customerNameDisplay) > 30 {
				customerNameDisplay = customerNameDisplay[:27] + "..."
			}
		}
	} else {
		pdf.SetFont("Poppins", "", 11)
	}
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, customerNameDisplay)
	pdf.Ln(lineHeight)

	pdf.SetX(rightCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Phone:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, ticketData.ContactNumber)
	pdf.Ln(lineHeight)

	pdf.SetX(rightCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Booking ID:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, fmt.Sprintf("#%d", ticketData.BookingID))
	pdf.Ln(lineHeight)

	pdf.SetX(rightCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Booking Time:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, ticketData.BookingTime)
	pdf.Ln(lineHeight)

	pdf.SetX(rightCol)
	pdf.SetFont("Poppins", "B", 11)
	pdf.SetTextColor(22, 26, 30)
	pdf.Cell(30, lineHeight, "Status:")
	pdf.SetFont("Poppins", "", 11)
	pdf.SetTextColor(64, 67, 72)
	pdf.Cell(60, lineHeight, ticketData.Status)
	pdf.Ln(20)

	qrTop := pdf.GetY() + extraHeight  // Adjust QR position based on extra height
	pdf.SetFillColor(240, 240, 245)    // background.secondary
	pdf.Rect(75, qrTop-5, 60, 70, "F") // Increased height for text
	pdf.Image("qrcode", 85, qrTop, 40, 0, false, "", 0, "")

	pdf.SetY(qrTop + 45)
	pdf.SetFont("Poppins", "", 8)   // Smaller font
	pdf.SetTextColor(142, 144, 145) // text.quaternary
	pdf.SetX(75)
	pdf.MultiCell(60, 4, "Scan this QR code at the theater entrance", "", "C", false)

	// Footer with brand color
	footerTop := 267.0
	pdf.SetFillColor(224, 75, 0) // Primary color
	pdf.Rect(0, footerTop, 210, 30, "F")

	// Thank you message in footer
	pdf.SetY(footerTop + 5)
	pdf.SetFont("Poppins", "B", 12)
	pdf.SetTextColor(255, 255, 255) // White text
	pdf.CellFormat(0, 10, "Thank you for choosing SKYFOX Cinemas!", "", 0, "C", false, 0, "")

	// Terms text in footer
	pdf.SetY(footerTop + 15)
	pdf.SetFont("Poppins", "", 8)
	pdf.SetTextColor(255, 255, 255) // White text for footer
	pdf.MultiCell(0, 5, "Please arrive 15 minutes before showtime. Tickets are non-refundable.", "", "C", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return "", utils.NewInternalServerError("PDF_GENERATION_FAILED", "Failed to generate PDF", err)
	}

	base64PDF := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64PDF, nil
}

func (s *bookingService) getTicketData(ctx context.Context, bookingID int) (*models.TicketData, error) {
	booking, err := s.bookingRepo.GetBookingById(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, utils.NewNotFoundError("BOOKING_NOT_FOUND", "Booking not found", nil)
	}

	seatNumbers, err := s.bookingSeatMappingRepo.GetSeatsByBookingId(ctx, booking.Id)
	if err != nil {
		return nil, err
	}

	show, err := s.showRepo.FindById(ctx, booking.ShowId)
	if err != nil {
		return nil, err
	}

	slot, err := s.slotRepo.GetSlotById(ctx, show.SlotId)
	if err != nil {
		return nil, err
	}

	movie, err := s.movieService.GetMovieById(ctx, show.MovieId)
	if err != nil {
		return nil, err
	}

	var customerName, contactNumber string

	if booking.CustomerId != nil {
		adminCustomer, err := s.adminBookedCustomerRepo.FindById(ctx, *booking.CustomerId)
		if err != nil {
			return nil, err
		}
		customerName = adminCustomer.Name
		contactNumber = adminCustomer.Number
	} else if booking.CustomerUsername != nil {
		customer, err := s.skyCustomerRepo.FindByUsername(ctx, *booking.CustomerUsername)
		if err != nil {
			return nil, err
		}
		customerName = customer.Name
		contactNumber = customer.Number
	}

	ticketData := &models.TicketData{
		BookingID:     booking.Id,
		ShowName:      movie.Name,
		ShowDate:      show.Date.Format("2006-01-02"),
		ShowTime:      slot.StartTime,
		CustomerName:  customerName,
		ContactNumber: contactNumber,
		AmountPaid:    booking.AmountPaid,
		NumberOfSeats: booking.NoOfSeats,
		SeatNumbers:   seatNumbers,
		Status:        booking.Status,
		BookingTime:   booking.BookingTime.Format("2006-01-02 15:04:05"),
		PaymentType:   string(booking.PaymentType),
	}

	return ticketData, nil
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Failed to read font file")
		return nil
	}
	return data
}

func addLogoWithFallback(pdf *gofpdf.Fpdf, x, y, width float64) {
	logoPath := "assets/images/logo.png"
	if _, err := os.Stat(logoPath); os.IsNotExist(err) {
		log.Warn().Msg("Logo file not found, skipping logo")
		return
	}

	pdf.Image(logoPath, x, y, width, 0, false, "", 0, "")
}

func addFontsWithFallback(pdf *gofpdf.Fpdf) {
	fontDir := "assets/fonts/poppins/"

	poppinsRegular := readFile(fontDir + "Poppins-Regular.ttf")
	poppinsBold := readFile(fontDir + "Poppins-Bold.ttf")
	poppinsItalic := readFile(fontDir + "Poppins-Italic.ttf")
	poppinsBoldItalic := readFile(fontDir + "Poppins-BoldItalic.ttf")

	if poppinsRegular != nil && poppinsBold != nil && poppinsItalic != nil && poppinsBoldItalic != nil {
		pdf.AddUTF8FontFromBytes("Poppins", "", poppinsRegular)
		pdf.AddUTF8FontFromBytes("Poppins", "B", poppinsBold)
		pdf.AddUTF8FontFromBytes("Poppins", "I", poppinsItalic)
		pdf.AddUTF8FontFromBytes("Poppins", "BI", poppinsBoldItalic)
	} else {
		log.Warn().Msg("Poppins font files not found, using Helvetica instead")
		pdf.SetFont("Helvetica", "", 12)
	}
}
