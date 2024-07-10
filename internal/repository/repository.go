package repository

import (
	"github.com/TimeTracker-Effective-Mobile/internal/model"
	"github.com/TimeTracker-Effective-Mobile/internal/repository/postgres"
)

type dbStorage interface {
	GetUsersInfo(query map[string][]string) ([]model.User, error)
	GetSortedTaskByUser(userId int, query map[string][]string) ([]model.Task, error)
	StartNewTask(userId int, name string) (model.Task, error)
	StartExistingTask(taskId int) error
	TaskExists(taskId int) bool
	UserExists(userId int) bool
	IsActiveTask(taskId int) bool
	StopTask(taskId int) (model.Task, error)
	DeleteUser(userId int) error
	UpdateUser(user model.User) error
	SaveUser(user *model.User) error
}

type repository struct {
	db dbStorage
}

func New() *repository {
	return &repository{db: postgres.New()}
}

func (r *repository) GetUsersInfo(query map[string][]string) ([]model.User, error) {
	return r.db.GetUsersInfo(query)
}

func (r *repository) GetSortedTaskByUser(userId int, query map[string][]string) ([]model.Task, error) {
	return r.db.GetSortedTaskByUser(userId, query)
}

func (r *repository) StartNewTask(userId int, name string) (model.Task, error) {
	return r.db.StartNewTask(userId, name)
}

func (r *repository) StartExistingTask(taskId int) error {
	return r.db.StartExistingTask(taskId)
}

func (r *repository) TaskExists(taskId int) bool {
	return r.db.TaskExists(taskId)
}

func (r *repository) UserExists(userId int) bool {
	return r.db.UserExists(userId)
}

func (r *repository) IsActiveTask(taskId int) bool {
	return r.db.IsActiveTask(taskId)
}

func (r *repository) StopTask(taskId int) (model.Task, error) {
	return r.db.StopTask(taskId)
}

func (r *repository) DeleteUser(userId int) error {
	return r.db.DeleteUser(userId)
}

func (r *repository) UpdateUser(user model.User) error {
	return r.db.UpdateUser(user)
}

func (r *repository) SaveUser(user *model.User) error {
	return r.db.SaveUser(user)
}
