package repositories

import (
	"context"
	"fmt"

	"github.com/iamsuteerth/skyfox-backend/pkg/models"
	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SkyCustomerRepository interface {
	FindByUsername(ctx context.Context, username string) (*models.SkyCustomer, error)
	FindByEmail(ctx context.Context, email string) (*models.SkyCustomer, error)
	ExistsByEmailOrMobile(ctx context.Context, email, mobileNumber string) (bool, string, error)
	Create(ctx context.Context, customer *models.SkyCustomer) error
	UpdateCustomerDetails(ctx context.Context, username string, updates map[string]interface{}) error
	GetCustomerProfileImg(ctx context.Context, username string) (string, error)
	UpdateProfileImageURL(ctx context.Context, username string, profileImgURL string) error
}

type skyCustomerRepository struct {
	db *pgxpool.Pool
}

func NewSkyCustomerRepository(db *pgxpool.Pool) SkyCustomerRepository {
	return &skyCustomerRepository{db: db}
}

func (repo *skyCustomerRepository) FindByUsername(ctx context.Context, username string) (*models.SkyCustomer, error) {
	query := `SELECT id, name, username, number, email, profile_img FROM customertable WHERE username = $1`

	var customer models.SkyCustomer
	err := repo.db.QueryRow(ctx, query, username).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Username,
		&customer.Number,
		&customer.Email,
		&customer.ProfileImg,
		&customer.SecurityAnswerHash,
		&customer.SecurityQuestionID,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching customer", err)
	}

	return &customer, nil
}

func (repo *skyCustomerRepository) FindByEmail(ctx context.Context, email string) (*models.SkyCustomer, error) {
	query := `SELECT id, name, username, number, email, profile_img, security_question_id, security_answer_hash FROM customertable WHERE email = $1`

	var customer models.SkyCustomer
	err := repo.db.QueryRow(ctx, query, email).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Username,
		&customer.Number,
		&customer.Email,
		&customer.ProfileImg,
		&customer.SecurityQuestionID,
		&customer.SecurityAnswerHash,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching customer by email", err)
	}

	return &customer, nil
}

func (repo *skyCustomerRepository) ExistsByEmailOrMobile(ctx context.Context, email, mobileNumber string) (bool, string, error) {
	query := `
		SELECT 
			CASE 
				WHEN email = $1 THEN 'email'
				WHEN number = $2 THEN 'mobilenumber'
				ELSE ''
			END
		FROM customertable 
		WHERE email = $1 OR number = $2
		LIMIT 1
	`

	var field string
	err := repo.db.QueryRow(ctx, query, email, mobileNumber).Scan(&field)

	if err != nil {
		if err == pgx.ErrNoRows {
			return false, "", nil
		}
		return false, "", utils.NewInternalServerError("DATABASE_ERROR", "Error checking customer existence", err)
	}

	return field != "", field, nil
}

func (repo *skyCustomerRepository) Create(ctx context.Context, customer *models.SkyCustomer) error {
	query := `
    INSERT INTO customertable (name, username, number, email, profile_img, security_question_id, security_answer_hash)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING id
    `

	err := repo.db.QueryRow(ctx, query,
		customer.Name,
		customer.Username,
		customer.Number,
		customer.Email,
		customer.ProfileImg,
		customer.SecurityQuestionID,
		customer.SecurityAnswerHash,
	).Scan(&customer.ID)

	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error creating customer", err)
	}

	return nil
}

func (repo *skyCustomerRepository) UpdateCustomerDetails(ctx context.Context, username string, updates map[string]interface{}) error {
	query := "UPDATE customertable SET "
	args := []interface{}{username}
	argIndex := 2

	for key, value := range updates {
		if argIndex > 2 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", key, argIndex)
		args = append(args, value)
		argIndex++
	}

	query += " WHERE username = $1"

	_, err := repo.db.Exec(ctx, query, args...)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error updating customer details", err)
	}

	return nil
}

func (repo *skyCustomerRepository) GetCustomerProfileImg(ctx context.Context, username string) (string, error) {
	query := "SELECT profile_img FROM customertable WHERE username = $1"

	var profileImg string
	err := repo.db.QueryRow(ctx, query, username).Scan(&profileImg)

	if err != nil {
		if err == pgx.ErrNoRows {
			return "", utils.NewNotFoundError("USER_NOT_FOUND", fmt.Sprintf("No user found with username: %s", username), nil)
		}
		return "", utils.NewInternalServerError("DATABASE_ERROR", "Error fetching profile image URL", err)
	}

	return profileImg, nil
}

func (repo *skyCustomerRepository) UpdateProfileImageURL(ctx context.Context, username string, profileImgURL string) error {
	query := "UPDATE customertable SET profile_img = $1 WHERE username = $2"

	_, err := repo.db.Exec(ctx, query, profileImgURL, username)
	if err != nil {
		return utils.NewInternalServerError("DATABASE_ERROR", "Error updating profile image URL", err)
	}

	return nil
}
