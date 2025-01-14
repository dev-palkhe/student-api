package repository

import (
	"database/sql"
	"log"

	"github.com/dev-palkhe/student-api/internal/models"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *PostgresStudentRepository {
	return &PostgresStudentRepository{db: db}
}

func (r *PostgresStudentRepository) Create(student models.Student) (models.Student, error) {
	err := r.db.QueryRow(
		"INSERT INTO students (id, name, age, course, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, age, course, created_at, updated_at",
		student.ID, student.Name, student.Age, student.Course, student.CreatedAt, student.UpdatedAt,
	).Scan(&student.ID, &student.Name, &student.Age, &student.Course, &student.CreatedAt, &student.UpdatedAt)

	if err != nil {
		log.Printf("Error inserting student: %v", err)
		return models.Student{}, err
	}
	return student, nil
}

func (r *PostgresStudentRepository) GetAll() ([]models.Student, error) {
	rows, err := r.db.Query("SELECT id, name, age, course, created_at, updated_at FROM students")
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
	if err := rows.Err(); err != nil { // Important: check for errors during iteration
		return nil, err
	}
	return students, nil
}

func (r *PostgresStudentRepository) GetByID(id uuid.UUID) (models.Student, error) {
	var s models.Student
	err := r.db.QueryRow("SELECT id, name, age, course, created_at, updated_at FROM students WHERE id = $1", id).Scan(
		&s.ID, &s.Name, &s.Age, &s.Course, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return models.Student{}, sql.ErrNoRows
	} else if err != nil {
		return models.Student{}, err
	}
	return s, nil
}

func (r *PostgresStudentRepository) Update(student models.Student) error {
	_, err := r.db.Exec(
		"UPDATE students SET name = $1, age = $2, course = $3, updated_at = $4 WHERE id = $5",
		student.Name, student.Age, student.Course, student.UpdatedAt, student.ID,
	)
	return err
}

func (r *PostgresStudentRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM students WHERE id = $1", id)
	return err
}
