package users

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"web-appli/src/db"
)

const RequestSelect = `SELECT id, firstname, lastname, login, isAdmin FROM users`
const RequestCreate = `INSERT INTO users(firstname, lastname, login, password, isAdmin) VALUES (?, ?, ?, ?, ?)`

var Dao = UsersDao{}

type UsersDao struct {
}

func (dao UsersDao) Create(user *User, password string) (*User, error) {
	var err error
	user.Id, err = db.ExecuteCreate(RequestCreate, user.Firstname, user.Lastname, user.Login, dao.passwordHash(password), user.IsAdmin)
	return user, err
}

func (dao UsersDao) FindAll() ([]*User, error) {
	return db.ExecuteQuery(dao.extractResults, RequestSelect)
}

func (dao UsersDao) FindByLoginAndPassword(login, password string) (*User, error) {
	conn := db.GetConnection()

	// VULNERABILITE - Injection SQL possible ici
	sqlRequest := RequestSelect + " WHERE login = '%LOGIN%' AND password = '%PASSWORD%' ORDER BY login"
	sqlRequest = strings.ReplaceAll(sqlRequest, "%LOGIN%", login)
	sqlRequest = strings.ReplaceAll(sqlRequest, "%PASSWORD%", dao.passwordHash(password))
	// log.Printf("authent request: %s", sqlRequest)
	row, err := conn.Query(sqlRequest)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %s", sqlRequest)
	}
	defer row.Close()
	users, err := dao.extractResults(row)
	if err != nil {
		return nil, err
	}
	if len(users) <= 0 {
		// Utilisateur non trouvé
		return nil, nil
	} else {
		// Utilisateur trouvé
		result := users[0]
		return result, nil
	}
}

func (dao UsersDao) SeachByFirstnameOrLastname(search string) ([]*User, error) {
	conn := db.GetConnection()
	sqlRequest := RequestSelect + " WHERE firstname like ? OR lastname like ? ORDER BY firstname"
	statement, err := conn.Prepare(sqlRequest)
	if err != nil {
		return nil, fmt.Errorf("query preparation failed: %s\nerror: %w", RequestCreate, err)
	}
	likeString := "%" + search + "%"
	row, err := statement.Query(likeString, likeString)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %s\nerror: %w", RequestCreate, err)
	}
	defer row.Close()
	return dao.extractResults(row)
}

func (dao UsersDao) extractResults(rows *sql.Rows) ([]*User, error) {
	users := make([]*User, 0, 0)
	for rows.Next() {
		currentUser := User{}
		err := rows.Scan(&currentUser.Id, &currentUser.Firstname, &currentUser.Lastname, &currentUser.Login, &currentUser.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("query result extraction failed\nerror: %w", err)
		}
		users = append(users, &currentUser)
	}
	// Retourner la liste
	return users, nil
}

func (dao UsersDao) passwordHash(password string) string {
	//hash := md5.Sum([]byte(password))
	//return base64.StdEncoding.EncodeToString(hash[:])
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}
