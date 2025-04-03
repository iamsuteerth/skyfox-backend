package repositories

import (
	"context"

	"github.com/iamsuteerth/skyfox-backend/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SecurityQuestion struct {
	ID       int    `json:"id"`
	Question string `json:"question"`
}

type SecurityQuestionRepository interface {
	FindAll(ctx context.Context) ([]SecurityQuestion, error)
	QuestionExists(ctx context.Context, id int) (bool, error)
	FindByID(ctx context.Context, id int) (*SecurityQuestion, error)
}

type securityQuestionRepository struct {
	db *pgxpool.Pool
}

func NewSecurityQuestionRepository(db *pgxpool.Pool) SecurityQuestionRepository {
	return &securityQuestionRepository{db: db}
}

func (repo *securityQuestionRepository) FindAll(ctx context.Context) ([]SecurityQuestion, error) {
	query := `SELECT id, question FROM security_questions ORDER BY id`

	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching security questions", err)
	}
	defer rows.Close()

	var questions []SecurityQuestion
	for rows.Next() {
		var question SecurityQuestion
		if err := rows.Scan(&question.ID, &question.Question); err != nil {
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error scanning security question", err)
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (repo *securityQuestionRepository) QuestionExists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM security_questions WHERE id = $1)`

	var exists bool
	err := repo.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, utils.NewInternalServerError("DATABASE_ERROR", "Error checking security question existence", err)
	}

	return exists, nil
}

func (repo *securityQuestionRepository) FindByID(ctx context.Context, id int) (*SecurityQuestion, error) {
	query := `SELECT id, question FROM security_questions WHERE id = $1`

	var question SecurityQuestion
	err := repo.db.QueryRow(ctx, query, id).Scan(&question.ID, &question.Question)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, utils.NewNotFoundError("SECURITY_QUESTION_NOT_FOUND", "Security question not found", nil)
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error fetching security question", err)
	}

	return &question, nil
}
