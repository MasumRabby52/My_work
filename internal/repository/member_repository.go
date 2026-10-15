package repository

import (
	"context"
	"my_work/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateMember(pool *pgxpool.Pool, member *models.Member) (*models.Member, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	INSERT INTO members(email,password)
	VALUES ($1,$2)
	RETURNING id,email,created_at,updated_at
	`

	err := pool.QueryRow(ctx, query, member.Email, member.Password).Scan(
		&member.ID,
		&member.Email,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return member, nil
}

func GetMemberByEmail(pool *pgxpool.Pool, email string) (*models.Member, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT *
	FROM members
	WHERE email=$1
	`
	var member models.Member
	err := pool.QueryRow(ctx, query, email).Scan(
		&member.ID,
		&member.Email,
		&member.Password,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &member, nil
}

func GetMemberByID(pool *pgxpool.Pool, id string) (*models.Member, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT *
	FROM members
	WHERE id=$1
	`
	var member models.Member
	err := pool.QueryRow(ctx, query, id).Scan(
		&member.ID,
		&member.Email,
		&member.Password,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &member, nil
}
