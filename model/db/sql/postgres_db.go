package sql

import (
	"CRUD/model"
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type postgresDB struct {
	db *gorm.DB
}

func NewPostgresDb() (*postgresDB, error) {
	dbUrl := os.Getenv("POSTGRES_URL")
	if dbUrl == "" {
		return nil, errors.New("POSTGRES_URL environment variable not set")
	}
	db, err := gorm.Open(postgres.Open(dbUrl))
	db.Config.PrepareStmt = true
	if err != nil {
		return nil, err
	}
	migrateErr := db.AutoMigrate(&model.Task{})
	if migrateErr != nil {
		return nil, migrateErr
	}
	return &postgresDB{db: db}, nil
}

func (postgres *postgresDB) CreateTask(task *model.Task) (uint64, error) {
	result := postgres.db.Create(task)
	if result.Error != nil {
		return 0, result.Error
	}
	return (*task).Id, nil
}

func (postgres *postgresDB) GetAllTasks() (*[]model.Task, error) {
	tasks := new([]model.Task)
	result := postgres.db.Find(tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (postgres *postgresDB) GetTaskById(taskId uint64) (*model.Task, error) {
	task := new(model.Task)
	result := postgres.db.First(task, taskId)
	if result.Error != nil {
		return nil, result.Error
	}
	return task, nil
}

func (postgres *postgresDB) UpdateTask(task *model.Task) (int64, error) {
	result := postgres.db.Model(task).Updates(task)
	return result.RowsAffected, result.Error
}

func (postgres *postgresDB) DeleteTask(taskId uint64) (int64, error) {
	result := postgres.db.Delete(&model.Task{}, taskId)
	return result.RowsAffected, result.Error
}
