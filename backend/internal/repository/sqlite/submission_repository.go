package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"spondtest/backend/internal/domain"
	"spondtest/backend/internal/repository"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(databasePath string) (*SubmissionRepository, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	repo := &SubmissionRepository{db: db}
	if err := repo.applySchema(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}

	return repo, nil
}

func (r *SubmissionRepository) Save(ctx context.Context, submission domain.Submission) error {
	createdAt := submission.CreatedAt.UTC().Format(time.RFC3339Nano)
	duplicateScope := repository.DuplicateScopeFromSubmission(submission)

	_, err := r.db.ExecContext(ctx, `
INSERT INTO submissions (
	submission_id,
	name,
	email,
	phone_number,
	birth_date,
	member_type_id,
	club_id,
	form_id,
	created_at,
	email_normalized,
	phone_normalized
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`,
		submission.SubmissionID,
		submission.Name,
		submission.Email,
		submission.PhoneNumber,
		duplicateScope.BirthDate,
		submission.MemberTypeID,
		duplicateScope.ClubID,
		duplicateScope.FormID,
		createdAt,
		duplicateScope.Email,
		duplicateScope.Phone,
	)
	if err != nil {
		if isUniqueConstraintError(err) {
			return repository.ErrDuplicateSubmission
		}

		return err
	}

	return nil
}

func (r *SubmissionRepository) List(ctx context.Context) ([]domain.Submission, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT
	submission_id,
	name,
	email,
	phone_number,
	birth_date,
	member_type_id,
	club_id,
	form_id,
	created_at
FROM submissions
ORDER BY created_at ASC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	submissions := make([]domain.Submission, 0)
	for rows.Next() {
		var submission domain.Submission
		var createdAt string

		if err := rows.Scan(
			&submission.SubmissionID,
			&submission.Name,
			&submission.Email,
			&submission.PhoneNumber,
			&submission.BirthDate,
			&submission.MemberTypeID,
			&submission.ClubID,
			&submission.FormID,
			&createdAt,
		); err != nil {
			return nil, err
		}

		parsedCreatedAt, err := parseTimestamp(createdAt)
		if err != nil {
			return nil, err
		}
		submission.CreatedAt = parsedCreatedAt

		submissions = append(submissions, submission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return submissions, nil
}

func (r *SubmissionRepository) Close() error {
	if r.db == nil {
		return nil
	}
	return r.db.Close()
}

func (r *SubmissionRepository) applySchema(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, `PRAGMA journal_mode = WAL;`); err != nil {
		return fmt.Errorf("set sqlite journal mode: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, `PRAGMA busy_timeout = 5000;`); err != nil {
		return fmt.Errorf("set sqlite busy timeout: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, schemaSQL); err != nil {
		return fmt.Errorf("apply sqlite schema: %w", err)
	}

	return nil
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint failed")
}

func parseTimestamp(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("missing created_at")
	}

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err == nil {
		return parsed, nil
	}

	parsed, err = time.Parse(time.RFC3339, value)
	if err == nil {
		return parsed, nil
	}

	return time.Time{}, fmt.Errorf("parse created_at timestamp: %w", err)
}
