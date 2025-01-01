package db

import (
	models "CRUD/model"
)

type TaskStorage interface {
	CreateTask(task *models.Task) (uint64, error)
	GetAllTasks() (*[]models.Task, error)
	GetTaskById(taskId uint64) (*models.Task, error)
	UpdateTask(task *models.Task) (int64, error)
	DeleteTask(taskId uint64) (int64, error)
}
