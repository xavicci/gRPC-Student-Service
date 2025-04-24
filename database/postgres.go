package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/xavicci/gRPC-Student-Service/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (repo *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)", student.Id, student.Name, student.Age)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var student = models.Student{}
	for rows.Next() {
		err := rows.Scan(&student.Id, &student.Name, &student.Age)
		if err != nil {
			return nil, err
		}
	}
	return &student, nil
}
