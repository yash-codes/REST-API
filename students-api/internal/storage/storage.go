package storage

import "github.com/yash-codes/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id string) (types.Student, error)
	GetAllStudentsList() ([]types.Student, error)
	DeleteStudentById(id string) (int64, error)
}
