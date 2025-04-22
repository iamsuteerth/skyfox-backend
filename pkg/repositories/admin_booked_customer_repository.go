package repositories

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type AdminBookedCustomerRepository interface {
	Create(ctx context.Context, customer *models.AdminBookedCustomer) error
	FindById(ctx context.Context, id int) (*models.AdminBookedCustomer, error)
}

type adminBookedCustomerRepository struct {
	db *pgxpool.Pool
}

func NewAdminBookedCustomerRepository(db *pgxpool.Pool) AdminBookedCustomerRepository {
	return &adminBookedCustomerRepository{db: db}
}

func (repo *adminBookedCustomerRepository) Create(ctx context.Context, customer *models.AdminBookedCustomer) error {
	query := `
		INSERT INTO admin_booked_customer (name, number)
		VALUES ($1, $2)
		RETURNING id
	`

	err := repo.db.QueryRow(ctx, query,
		customer.Name,
		customer.Number,
	).Scan(&customer.Id)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create admin booked customer")
		return utils.NewInternalServerError("DATABASE_ERROR", "Failed to create customer record", err)
	}

	return nil
}

func (repo *adminBookedCustomerRepository) FindById(ctx context.Context, id int) (*models.AdminBookedCustomer, error) {
	query := `
		SELECT id, name, number
		FROM admin_booked_customer
		WHERE id = $1
	`

	var customer models.AdminBookedCustomer
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&customer.Id,
		&customer.Name,
		&customer.Number,
	)

	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("Failed to find admin booked customer by ID")
		return nil, utils.NewNotFoundError("CUSTOMER_NOT_FOUND", "Customer not found", err)
	}

	return &customer, nil
}
