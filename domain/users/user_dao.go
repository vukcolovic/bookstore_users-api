package users

import (
	"UdemyApp/bookstore_users-api/datasources/postgresql/users_db"
	"UdemyApp/bookstore_users-api/utils/date_utils"
	"UdemyApp/bookstore_users-api/utils/errors"
	"fmt"
	"strings"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created) VALUES($1, $2, $3, $4);"
	queryGetUserById = "SELECT * FROM users WHERE id = $1;"
	queryUpdateUser = "UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4;"
	queryDeleteUser = "DELETE FROM users WHERE id = $1;"
	errorNoRows = "no rows in result set"

)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUserById)
	if err != nil {
		return errors.NewInternalServerError("prepare statment error : " + err.Error())
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); err != nil {
		if strings.Contains(err.Error(), errorNoRows) {
			return errors.NewNotFoundError(fmt.Sprintf("User %d not found", user.Id))
		}
		return errors.NewInternalServerError(fmt.Sprintf("error when trying to scan user with id %d: %s", user.Id, err.Error()))
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError("prepare statment error : " + err.Error())
	}
	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated)
	if err != nil {
		return errors.NewInternalServerError("insert user error : " + err.Error())
	}

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError("prepare statment error : " + err.Error())
	}
	defer stmt.Close()


	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return errors.NewInternalServerError("update user error : " + err.Error())
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		return errors.NewInternalServerError("prepare statment error : " + err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	if err != nil {
		return errors.NewInternalServerError("delete user error : " + err.Error())
	}

	return nil
}