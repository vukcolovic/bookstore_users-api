package users

import (
	"UdemyApp/bookstore_users-api/datasources/postgresql/users_db"
	"UdemyApp/bookstore_users-api/logger"
	"UdemyApp/bookstore_users-api/utils/errors"
	"database/sql"
	"fmt"
)

const (
	queryInsertUser = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES($1, $2, $3, $4, $5, $6);"
	queryGetUserById = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE id = $1;"
	queryUpdateUser = "UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE id = $4;"
	queryDeleteUser = "DELETE FROM users WHERE id = $1;"
	queryFindByStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status = $1;"
	queryFindByEmailPasswordAndStatus = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email = $1 AND Password = $2 AND Status = $3;"
	errorNoRows = "no rows in result set"

)

func (user *User) Get() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUserById)
	if err != nil {
		logger.Error("error when trying to prepare get user statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Id)
	if err := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
		logger.Error("error when trying to get user by id", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Save() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if err != nil {
		logger.Error("error when trying to save user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Update() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()


	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		logger.Error("error when trying to update user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) Delete() *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	if err != nil {
		logger.Error("error when trying to delete user", err)
		return errors.NewInternalServerError("database error")
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find user statment", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find user", err)
		return nil, errors.NewInternalServerError("database error")
	}
	defer rows.Close()

	result := make([]User, 0)

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when scan user row trying to find user", err)
			return nil, errors.NewInternalServerError("database error")
		}
		result = append(result, user)
	}

	if len(result) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return result, nil
}

func (user *User) FindByEmailAndPassword()  *errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailPasswordAndStatus)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statment", err)
		return errors.NewInternalServerError("database error")
	}
	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)
	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); getErr != nil {
		if getErr == sql.ErrNoRows {
			return errors.NewNotFoundError("invalid user credentials")
		}
		logger.Error("error when trying to get user by email and password", getErr)
		return errors.NewInternalServerError("database error")
	}

	return nil
}