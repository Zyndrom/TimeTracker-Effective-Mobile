package router

import (
	"net/http"
	"strconv"

	_ "github.com/TimeTracker-Effective-Mobile/docs"
	"github.com/TimeTracker-Effective-Mobile/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type router struct {
	ginRouter   *gin.Engine
	timeService timeTrackerService
}

type timeTrackerService interface {
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
	AddUser(passport string) (model.User, error)
}

func New(timeService timeTrackerService) router {
	router := router{
		ginRouter:   gin.New(),
		timeService: timeService,
	}
	router.ginRouter.Use(CORSMiddleware())
	router.ginRouter.GET("/users", router.getUsers())
	router.ginRouter.GET("/users/:user/workhours", router.getWorkHoursByUser())
	router.ginRouter.POST("/tasks/start-new", router.startNewTask())
	router.ginRouter.POST("/tasks/start-existed", router.startExistedTask())
	router.ginRouter.POST("/tasks/stop", router.stopTask())
	router.ginRouter.DELETE("/users/:user", router.deleteUser())
	router.ginRouter.PUT("/users/:user", router.updateUser())
	router.ginRouter.POST("/users", router.addUser())
	router.initSwagger()

	return router
}

func (r *router) initSwagger() {
	r.ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func (r *router) StartServer() {
	r.ginRouter.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// @Summary Get users
// @Description Retrieve a list of users based on query parameters
// @Accept json
// @Produce json
// @Param id query string false "ID"
// @Param name query string false "name"
// @Param passportNumber query string false "passportNumber"
// @Param surname query string false "surname"
// @Param patronymic query string false "patronymic"
// @Param address query string false "address"
// @Param page query string false "page"
// @Param limit query string false "limit"
// @Param offset query string false "offset"
// @Success 200 {array} model.User "List of users"
// @Failure 400 {string} string "Bad request"
// @Router /users [get]
func (r *router) getUsers() func(c *gin.Context) {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()
		user, err := r.timeService.GetUsersInfo(query)
		if err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

// @Summary Get work hours by user
// @Description Retrieves sorted tasks and work hours for a specific user
// @Accept json
// @Produce json
// @Param user path int true "User ID"
// @Param dateFrom query string false "Date From"
// @Param dateTo query string false "Date To"
// @Success 200 {array} model.Task "List of sorted tasks"
// @Failure 400 {string} string "Bad request"
// @Router /users/{user}/workhours [get]
func (r *router) getWorkHoursByUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("user"))
		if err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		if !r.timeService.UserExists(userId) {
			c.JSON(http.StatusBadRequest, "user not exist")
			return
		}
		query := c.Request.URL.Query()
		tasks, err := r.timeService.GetSortedTaskByUser(userId, query)
		if err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		c.JSON(http.StatusOK, tasks)
	}
}

type startNewTaskBody struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
}

// @Summary Start New Task
// @Description Starts a new task
// @Accept json
// @Produce json
// @Param task body startNewTaskBody true "Task details"
// @Success 200 {string} string "Task Started"
// @Success 201 {object} model.Task
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks/start-new [post]
func (r *router) startNewTask() func(c *gin.Context) {
	return func(c *gin.Context) {
		body := startNewTaskBody{}
		err := c.ShouldBindJSON(&body)
		if err != nil || body.Name == "" && body.UserId == 0 {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		var Task model.Task
		if !r.timeService.UserExists(body.UserId) {
			c.JSON(http.StatusBadRequest, "user not exist")
			return
		}
		Task, err = r.timeService.StartNewTask(body.UserId, body.Name)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusCreated, Task)

	}
}

type startExistedTaskBody struct {
	TaskId int `json:"task_id"`
}

// @Summary Resumes Existed Task
// @Description Resumes an existing
// @Accept json
// @Produce json
// @Param task body startExistedTaskBody true "Task details"
// @Success 200 {string} string "Task Started"
// @Success 201 {object} model.Task
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks/start-existed [post]
func (r *router) startExistedTask() func(c *gin.Context) {
	return func(c *gin.Context) {
		body := startExistedTaskBody{}
		err := c.ShouldBindJSON(&body)
		if err != nil || body.TaskId == 0 {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		if !r.timeService.TaskExists(body.TaskId) {
			c.JSON(http.StatusBadRequest, "task not exist")
			return
		}
		err = r.timeService.StartExistingTask(body.TaskId)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusOK, "Task Started")

	}
}

type stopTaskBody struct {
	TaskId int `json:"task_id"`
}

// @Summary Stop a task
// @Description Stop an active task
// @Accept json
// @Produce json
// @Param request body stopTaskBody true "Task stop request"
// @Success 200 {object} model.Task "Stopped task"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /tasks/stop [post]
func (r *router) stopTask() func(c *gin.Context) {
	return func(c *gin.Context) {
		body := stopTaskBody{}
		err := c.ShouldBindJSON(&body)
		if err != nil || body.TaskId == 0 {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		if !r.timeService.TaskExists(body.TaskId) {
			c.JSON(http.StatusBadRequest, "task not exist")
			return
		}
		if !r.timeService.IsActiveTask(body.TaskId) {
			c.JSON(http.StatusBadRequest, "task not active")
			return
		}
		task, err := r.timeService.StopTask(body.TaskId)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusOK, task)
	}
}

// @Summary Delete a user
// @Description Delete a user by their ID
// @Accept json
// @Produce json
// @Param user path int true "User ID"
// @Success 200 {string} string "User Deleted"
// @Failure 400 {string} string "Bad request"
// @Failure 400 {string} string "user not exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/{user} [delete]
func (r *router) deleteUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("user"))
		if err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		if !r.timeService.UserExists(userId) {
			c.JSON(http.StatusBadRequest, "user not exist")
			return
		}
		err = r.timeService.DeleteUser(userId)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusOK, "User Deleted")

	}
}

// @Summary Update a user
// @Description Update user information by their ID
// @Accept json
// @Produce json
// @Param user path int true "User ID"
// @Param request body model.User true "User update information"
// @Success 200 {string} string "User Updated"
// @Failure 400 {string} string "Bad request"
// @Failure 400 {string} string "user not exist"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users/{user} [put]
func (r *router) updateUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Param("user"))
		if err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		var user model.User

		if err := c.ShouldBindJSON(&user); err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		user.Id = userId
		if !r.timeService.UserExists(userId) {
			c.JSON(http.StatusBadRequest, "user not exist")
			return
		}
		err = r.timeService.UpdateUser(user)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusOK, "User Updated")
	}
}

type addNewUserBody struct {
	PassportNumber string `json:"passportNumber" example:"1234 567890"`
}

// @Summary Add a new user
// @Description Create a new user with the given passport number
// @Accept json
// @Produce json
// @Param request body addNewUserBody true "User creation request"
// @Success 201 {object} model.User "Created user"
// @Failure 400 {string} string "Bad request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /users [post]
func (r *router) addUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		data := addNewUserBody{}
		if err := c.ShouldBindJSON(&data); err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		user, err := r.timeService.AddUser(data.PassportNumber)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusCreated, user)
	}

}
