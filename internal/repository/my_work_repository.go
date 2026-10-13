package repository

import (
	"context"
	"my_work/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateMyWork(pool *pgxpool.Pool, title string, completed bool) (*models.MyWork, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	INSERT INTO my_works(title,completed)
	VALUES($1,$2)
	RETURNING id ,title,completed,created_at,updated_at
	`
	var myWork models.MyWork
	var err error = pool.QueryRow(ctx, query, title, completed).Scan(
		&myWork.ID,
		&myWork.Title,
		&myWork.Completed,
		&myWork.CreatedAt,
		&myWork.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &myWork, nil
}

func GetAllTodos(pool *pgxpool.Pool) ([]models.MyWork, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT * FROM my_works
	`
	var rows, err = pool.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var myWorks []models.MyWork = []models.MyWork{}

	for rows.Next() {
		var myWork models.MyWork
		err = rows.Scan(
			&myWork.ID,
			&myWork.Title,
			&myWork.Completed,
			&myWork.CreatedAt,
			&myWork.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		myWorks = append(myWorks, myWork)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return myWorks, nil

}
