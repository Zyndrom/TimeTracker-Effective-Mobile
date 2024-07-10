package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/TimeTracker-Effective-Mobile/internal/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type postgresql struct {
	db *sql.DB
}

func New() *postgresql {
	usr := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DATABASE")
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		usr, pass, host, port, dbName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf(err.Error())
	}
	psg := &postgresql{db: db}
	return psg
}

func (p *postgresql) SaveUser(user *model.User) error {
	query := `INSERT INTO users (passport_number, name, surname, patronymic, address) VALUES ($1, $2, $3, $4, $5) returning id;`
	err := p.db.QueryRow(query, user.PassportNumber, user.Name, user.Surname, user.Patronymic, user.Address).Scan(&user.Id)
	return err
}

func (p *postgresql) UpdateUser(user model.User) error {
	query := `UPDATE users SET passport_number = $1, name = $2, surname= $3, patronymic= $4, address= $5 WHERE id = $6;`

	_, err := p.db.Exec(query, user.PassportNumber, user.Name, user.Surname, user.Patronymic, user.Address, user.Id)
	return err
}

func (p *postgresql) DeleteUser(userId int) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := p.db.Exec(query, userId)
	return err
}

func (p *postgresql) UserExists(userId int) bool {
	query := `SELECT COUNT(*) FROM users WHERE id = $1;`
	row := p.db.QueryRow(query, userId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		logrus.Fatalln(err)
		return true
	}
	return count > 0

}

func (p *postgresql) GetUsersInfo(query map[string][]string) ([]model.User, error) {
	page := 1
	if val, ok := query["page"]; ok {
		num, err := strconv.Atoi(val[0])
		if err == nil {
			page = num
		}

	}
	limit := 10
	if val, ok := query["limit"]; ok {
		num, err := strconv.Atoi(val[0])
		if err == nil {
			limit = num
		}

	}
	offset := limit * (page - 1)
	if val, ok := query["offset"]; ok {
		num, err := strconv.Atoi(val[0])
		if err == nil {
			offset = num
		}
	}
	sqlQuery, params := generateFilter(query, limit, offset)
	var rows *sql.Rows
	var err error
	users := []model.User{}
	if len(params) <= 0 {
		rows, err = p.db.Query(`SELECT * FROM users LIMIT $1 OFFSET $2`, limit, offset)
	} else {
		rows, err = p.db.Query(sqlQuery, params...)
	}
	if err != nil {
		logrus.Fatalln(err)
		return users, err
	}

	for rows.Next() {
		user := model.User{}
		patronymic := sql.NullString{}
		err := rows.Scan(&user.Id, &user.PassportNumber, &user.Name, &user.Surname, &patronymic, &user.Address)
		user.Patronymic = patronymic.String
		if err != nil {
			logrus.Debug(err)
			continue
		}
		users = append(users, user)
	}
	return users, nil
}

func generateFilter(query map[string][]string, limit, offset int) (string, []any) {
	sqlQuery := `SELECT * FROM users WHERE`
	arr := []any{}

	if val, ok := query["id"]; ok {
		sqlQuery += fmt.Sprintf(" id = $%d,", len(arr)+1)
		arr = append(arr, val[0])
	}
	if val, ok := query["passportNumber"]; ok {
		sqlQuery += fmt.Sprintf(" passport_number = $%d,", len(arr)+1)
		arr = append(arr, val[0])
	}
	if val, ok := query["name"]; ok {
		sqlQuery += fmt.Sprintf(" name = $%d,", len(arr)+1)
		arr = append(arr, val[0])
	}
	if val, ok := query["surname"]; ok {
		sqlQuery += fmt.Sprintf(" surname = $%d,", len(arr)+1)
		arr = append(arr, val[0])
	}
	if val, ok := query["patronymic"]; ok {
		sqlQuery += fmt.Sprintf(" patronymic = $%d,", len(arr)+1)
		arr = append(arr, val[0])
	}
	if val, ok := query["address"]; ok {
		sqlQuery += fmt.Sprintf(" address = $%d,", len(arr)+1)
		arr = append(arr, val[0])
	}
	sqlQuery = sqlQuery[:len(sqlQuery)-1]
	sqlQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d;", len(arr)+1, len(arr)+2)
	arr = append(arr, limit, offset)
	return sqlQuery, arr
}

func (p *postgresql) StartNewTask(userId int, name string) (model.Task, error) {
	query := `INSERT INTO tasks (owner, name) VALUES ($1, $2) returning id, created_at, updated_at, active, duration;`
	task := model.Task{}
	err := p.db.QueryRow(query, userId, name).Scan(&task.Id, &task.CreatedAt, &task.UpdatedAt, &task.IsActive, &task.Duration)

	if err != nil {
		logrus.Fatalln(err)
		return task, err
	}
	return task, nil
}

func (p *postgresql) StartExistingTask(taskId int) error {
	query := `UPDATE tasks SET active = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1;`

	_, err := p.db.Exec(query, taskId)

	return err
}

func (p *postgresql) TaskExists(taskId int) bool {
	query := `SELECT COUNT(*) FROM tasks WHERE id = $1;`
	row := p.db.QueryRow(query, taskId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		logrus.Fatalln(err)
		return true
	}
	return count > 0

}
func (p *postgresql) GetTask(taskId int) (model.Task, error) {
	query := `SELECT * FROM tasks WHERE id = $1;`
	task := model.Task{}
	row := p.db.QueryRow(query, taskId)
	err := row.Scan(&task.Id, &task.Owner.Id, &task.Name, &task.CreatedAt, &task.UpdatedAt, &task.IsActive, &task.Duration)

	if err != nil {
		logrus.Fatalln(err)
		return task, err
	}
	return task, nil
}

func (p *postgresql) IsActiveTask(taskId int) bool {
	query := `SELECT active FROM tasks WHERE id = $1;`

	row := p.db.QueryRow(query, taskId)
	var active bool
	err := row.Scan(&active)
	if err != nil {
		logrus.Fatalln(err)
		return false
	}
	return active

}

func (p *postgresql) StopTask(taskId int) (model.Task, error) {
	query := `UPDATE tasks SET duration = duration + EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - updated_at)), active = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id = $1;`
	_, err := p.db.Exec(query, taskId)
	task := model.Task{}
	if err != nil {
		return task, err
	}
	return p.GetTask(taskId)
}

func (p *postgresql) GetSortedTaskByUser(userId int, query map[string][]string) ([]model.Task, error) {
	SQLQuery := `SELECT * FROM tasks WHERE owner = $1`
	args := []any{userId}
	if val, ok := query["dateFrom"]; ok {
		_, err := time.Parse(time.RFC3339, val[0])
		if err == nil {
			SQLQuery += fmt.Sprintf(" AND updated_at > $%d", len(args)+1)
			args = append(args, val[0])
		}
	}
	if val, ok := query["dateTo"]; ok {
		_, err := time.Parse(time.RFC3339, val[0])
		if err == nil {
			SQLQuery += fmt.Sprintf(" AND updated_at < $%d", len(args)+1)
			args = append(args, val[0])
		}
	}
	tasks := []model.Task{}

	updateDurationQuery := `UPDATE tasks SET duration = duration + EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - updated_at)), updated_at = CURRENT_TIMESTAMP WHERE owner = $1;`
	_, err := p.db.Exec(updateDurationQuery, userId)
	if err != nil {
		logrus.Fatalln(err)
		return tasks, err
	}

	SQLQuery += " ORDER BY duration DESC"
	rows, err := p.db.Query(SQLQuery, args...)
	if err != nil {
		logrus.Fatalln(err)
		return tasks, err
	}

	for rows.Next() {
		task := model.Task{}
		err := rows.Scan(&task.Id, &task.Owner.Id, &task.Name, &task.CreatedAt, &task.UpdatedAt, &task.IsActive, &task.Duration)
		if err != nil {
			logrus.Debug(err)
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}
