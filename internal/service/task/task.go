package task

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/TimeTracker-Effective-Mobile/internal/model"
	"github.com/sirupsen/logrus"
)

type taskService struct {
	storage storage
}

type storage interface {
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

var (
	externalApi string
)

func New(storage storage) *taskService {
	externalApi = os.Getenv("EXTERNAL_USER_API")
	return &taskService{
		storage: storage,
	}
}

func (t *taskService) AddUser(passport string) (model.User, error) {
	fields := strings.Fields(passport)
	query := fmt.Sprintf("%s/info?passportSerie=%s&passportNumber=%s", externalApi, fields[0], fields[1])
	var user model.User
	resp, err := http.Get(query)
	if err != nil {
		logrus.Debug(err)
		return user, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		logrus.Debug(err)
		return user, err
	}
	if err := json.Unmarshal(body, &user); err != nil {
		logrus.Debugf("Can not unmarshal JSON, %s", err.Error())
		return user, err
	}
	return user, t.storage.SaveUser(&user)
}

func (t *taskService) GetUsersInfo(query map[string][]string) ([]model.User, error) {
	return t.storage.GetUsersInfo(query)
}

func (t *taskService) GetSortedTaskByUser(userId int, query map[string][]string) ([]model.Task, error) {
	return t.storage.GetSortedTaskByUser(userId, query)
}

func (t *taskService) StartNewTask(userId int, name string) (model.Task, error) {
	return t.storage.StartNewTask(userId, name)
}

func (t *taskService) StartExistingTask(taskId int) error {
	return t.storage.StartExistingTask(taskId)
}

func (t *taskService) TaskExists(taskId int) bool {
	return t.storage.TaskExists(taskId)
}

func (t *taskService) UserExists(userId int) bool {
	return t.storage.UserExists(userId)
}

func (t *taskService) IsActiveTask(taskId int) bool {
	return t.storage.IsActiveTask(taskId)
}

func (t *taskService) StopTask(taskId int) (model.Task, error) {
	return t.storage.StopTask(taskId)
}

func (t *taskService) DeleteUser(userId int) error {
	return t.storage.DeleteUser(userId)
}

func (t *taskService) UpdateUser(user model.User) error {
	return t.storage.UpdateUser(user)
}
