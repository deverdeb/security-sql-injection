package tasks

import (
	"database/sql"
	"fmt"
	"strconv"
	"web-appli/src/db"
)

const RequestSelect = `SELECT id, userId, name, description, priority, status, archived FROM tasks`
const RequestCreate = `INSERT INTO tasks(userId, name, description, priority, status, archived) VALUES (?, ?, ?, ?, ?, ?)`
const RequestUpdate = `UPDATE tasks SET userId = ?, name = ?, description = ?, priority = ?, status = ?, archived = ? WHERE id = ?`
const RequestDelete = `DELETE FROM tasks WHERE id = ?`

var Dao = TasksDao{}

type TasksDao struct {
}

func (dao TasksDao) Create(task *Task) (*Task, error) {
	conn := db.GetConnection()

	statement, err := conn.Prepare(RequestCreate)
	if err != nil {
		return nil, fmt.Errorf("query preparation failed: %s\nerror: %w", RequestCreate, err)
	}
	result, err := statement.Exec(task.UserId, task.Name, task.Description, task.Priority, task.Status, task.Archived)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %s\nerror: %w", RequestCreate, err)
	}
	task.Id, _ = result.LastInsertId()
	return task, nil
}

func (dao TasksDao) Update(task *Task) (*Task, error) {
	conn := db.GetConnection()

	statement, err := conn.Prepare(RequestUpdate)
	if err != nil {
		return nil, fmt.Errorf("query preparation failed: %s\nerror: %w", RequestUpdate, err)
	}
	_, err = statement.Exec(task.UserId, task.Name, task.Description, task.Priority, task.Status, task.Archived, task.Id)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %s\nerror: %w", RequestUpdate, err)
	}
	return task, nil
}

func (dao TasksDao) Delete(task *Task) error {
	conn := db.GetConnection()

	statement, err := conn.Prepare(RequestDelete)
	if err != nil {
		return fmt.Errorf("query preparation failed: %s\nerror: %w", RequestDelete, err)
	}
	_, err = statement.Exec(task.Id)
	if err != nil {
		return fmt.Errorf("query execution failed: %s\nerror: %w", RequestDelete, err)
	}
	return nil
}

func (dao TasksDao) FindAll() ([]*Task, error) {
	conn := db.GetConnection()
	sqlRequest := RequestSelect + " ORDER BY name"
	statement, err := conn.Prepare(sqlRequest)
	if err != nil {
		return nil, fmt.Errorf("query preparation failed: %s\nerror: %w", RequestCreate, err)
	}
	row, err := statement.Query()
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %s\nerror: %w", RequestCreate, err)
	}
	defer row.Close()
	return dao.extractResults(row)
}

// FindByIdAndUserId est vulnérable.
// Nous aurions dû passer par une requête préparée pour passer l'identifiant.
// Et d'ailleurs, nous devrions passer l'identifiant sous forme de "int64".
func (dao TasksDao) FindByIdAndUserId(id string, userId int64) (*Task, error) {
	sqlRequest := RequestSelect + " WHERE id = " + id + " AND userId = " + strconv.FormatInt(userId, 10) + " ORDER BY id"
	result, err := db.ExecuteQuery(dao.extractResults, sqlRequest)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	} else {
		return result[0], nil
	}
}

func (dao TasksDao) FindByUserId(userId int64) ([]*Task, error) {
	sqlRequest := RequestSelect + " WHERE userId = ? ORDER BY id"
	return db.ExecuteQuery(dao.extractResults, sqlRequest, userId)
}

func (dao TasksDao) SearchByText(userId int64, search string) ([]*Task, error) {
	conn := db.GetConnection()
	sqlRequest := RequestSelect + " WHERE userId = ? and (name like ? or description like ?) ORDER BY name"
	statement, err := conn.Prepare(sqlRequest)
	if err != nil {
		return nil, fmt.Errorf("query preparation failed: %s\nerror: %w", RequestCreate, err)
	}
	likeString := "%" + search + "%"
	row, err := statement.Query(userId, likeString, likeString)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %s\nerror: %w", RequestCreate, err)
	}
	defer row.Close()
	return dao.extractResults(row)
}

func (dao TasksDao) extractResults(rows *sql.Rows) ([]*Task, error) {
	tasks := make([]*Task, 0, 0)
	for rows.Next() {
		currentTask := Task{}
		err := rows.Scan(&currentTask.Id, &currentTask.UserId, &currentTask.Name, &currentTask.Description,
			&currentTask.Priority, &currentTask.Status, &currentTask.Archived)
		if err != nil {
			return nil, fmt.Errorf("query result extraction failed\nerror: %w", err)
		}
		tasks = append(tasks, &currentTask)
	}
	// Retourner la liste
	return tasks, nil
}
