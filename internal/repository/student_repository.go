package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dev-palkhe/student-api/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type StudentRepository interface {
	Create(ctx context.Context, student models.Student) (models.Student, error)
	GetAll(ctx context.Context) ([]models.Student, error)
	GetByID(ctx context.Context, id uuid.UUID) (models.Student, error)
	Update(ctx context.Context, student models.Student) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type postgresStudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &postgresStudentRepository{db: db}
}

func NewPostgresDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database")
	return db, nil
}

func (r *postgresStudentRepository) Create(ctx context.Context, student models.Student) (models.Student, error) {
	err := r.db.QueryRowContext(ctx, "INSERT INTO students (id, name, age, course) VALUES ($1, $2, $3, $4) RETURNING id, name, age, course, created_at, updated_at",
		student.ID, student.Name, student.Age, student.Course).Scan(&student.ID, &student.Name, &student.Age, &student.Course, &student.CreatedAt, &student.UpdatedAt)
	if err != nil {
		log.Printf("Error inserting student: %v", err)
		return models.Student{}, err
	}
	return student, nil
}

func (r *postgresStudentRepository) GetAll(ctx context.Context) ([]models.Student, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, age, course, created_at, updated_at FROM students")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Course, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}

func (r *postgresStudentRepository) GetByID(ctx context.Context, id uuid.UUID) (models.Student, error) {
	var s models.Student
	err := r.db.QueryRowContext(ctx, "SELECT id, name, age, course, created_at, updated_at FROM students WHERE id = $1", id).Scan(
		&s.ID, &s.Name, &s.Age, &s.Course, &s.CreatedAt, &s.UpdatedAt)
	if err == sql.ErrNoRows {
		return models.Student{}, sql.ErrNoRows // Return sql.ErrNoRows directly
	} else if err != nil {
		return models.Student{}, err
	}
	return s, nil
}

func (r *postgresStudentRepository) Update(ctx context.Context, student models.Student) error {
	_, err := r.db.ExecContext(ctx, "UPDATE students SET name = $1, age = $2, course = $3, updated_at = $4 WHERE id = $5",
		student.Name, student.Age, student.Course, student.UpdatedAt, student.ID)
	return err
}

func (r *postgresStudentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM students WHERE id = $1", id)
	return err
}
