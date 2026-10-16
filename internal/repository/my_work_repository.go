package repository

import (
	"context"
	"fmt"
	"my_work/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateMyWork(pool *pgxpool.Pool, title string, completed bool, memberID string) (*models.MyWork, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	INSERT INTO my_works(title,completed,member_id)
	VALUES($1,$2,$3)
	RETURNING id ,title,completed,created_at,updated_at,member_id
	`
	var myWork models.MyWork
	var err error = pool.QueryRow(ctx, query, title, completed, memberID).Scan(
		&myWork.ID,
		&myWork.Title,
		&myWork.Completed,
		&myWork.CreatedAt,
		&myWork.UpdatedAt,
		&myWork.MemberID,
	)
	if err != nil {
		return nil, err
	}
	return &myWork, nil
}

func GetAllMyWorks(pool *pgxpool.Pool, memberID string) ([]models.MyWork, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT * FROM my_works
	WHERE member_id=$1
	ORDER BY created_at DESC
	`
	var rows, err = pool.Query(ctx, query, memberID)

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
			&myWork.MemberID,
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

func GetMyWorkByID(pool *pgxpool.Pool, id int, memberID string) (*models.MyWork, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT * FROM my_works
	WHERE id=$1 and member_id=$2
	`
	var myWork models.MyWork
	var err error = pool.QueryRow(ctx, query, id, memberID).Scan(
		&myWork.ID,
		&myWork.Title,
		&myWork.Completed,
		&myWork.CreatedAt,
		&myWork.UpdatedAt,
		&myWork.MemberID,
	)

	if err != nil {
		return nil, err
	}
	return &myWork, nil
}

func UpdateMyWork(pool *pgxpool.Pool, id int, title string, completed bool, memberID string) (*models.MyWork, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	UPDATE my_works
	SET title=$1,completed=$2,updated_at=CURRENT_TIMESTAMP
	WHERE id=$3 and member_id=$4
	RETURNING id,title,completed,created_at,updated_at,member_id
	`

	var myWork models.MyWork
	var err error = pool.QueryRow(ctx, query, title, completed, id, memberID).Scan(
		&myWork.ID,
		&myWork.Title,
		&myWork.Completed,
		&myWork.CreatedAt,
		&myWork.UpdatedAt,
		&myWork.MemberID,
	)

	if err != nil {
		return nil, err
	}
	return &myWork, nil
}
func DeleteMyWork(pool *pgxpool.Pool, id int, memberID string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	DELETE FROM my_works
	WHERE id=$1 and member_id=$2
    `

	var commandTag, err = pool.Exec(ctx, query, id, memberID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("myWork with id %d  not found", id)
	}
	return nil
}
