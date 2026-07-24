package dbutils

import (
	"errors"
	"strings"
)

// errorFilters is a list of validator functions to identify and convert raw database errors
var errorFilters = []func(err error) (bool, error){
	filterDuplicationUsername,
	filterDuplicationEmail,
	filterRecordNotFound,
}

// Common database errors returned by the repository layer.
var (
	ErrDuplicationUsername = errors.New("username already exists")
	ErrDuplicationEmail    = errors.New("email already exists")
	ErrRecordNotFound      = errors.New("record not found")
)

// CatchDBError inspects the raw database error and maps it to a clean application error if a match is found
func CatchDBError(err error) error {
	if err == nil {
		return nil
	}
	for _, filter := range errorFilters {
		match, filteredErr := filter(err)
		if match {
			return filteredErr
		}
	}
	return err
}

// filterDuplicationUsername checks if the error is caused by a duplicate username constraint
func filterDuplicationUsername(err error) (bool, error) {
	errStr := strings.ToLower(err.Error())
	//check both Postgres and Sqlite duplicate username error
	return strings.Contains(errStr, `duplicate key value violates unique constraint "uni_users_username"`) || strings.Contains(errStr, `unique constraint failed: users.username`), ErrDuplicationUsername
}

// filterDuplicationEmail checks if the error is caused by a duplicate email constraint
func filterDuplicationEmail(err error) (bool, error) {
	errStr := strings.ToLower(err.Error())
	//check both Postgres and Sqlite duplicate email error
	return strings.Contains(errStr, `duplicate key value violates unique constraint "uni_users_email"`) || strings.Contains(errStr, `unique constraint failed: users.email`), ErrDuplicationEmail
}

// filterRecordNotFound checks if the error is caused by a missing database record
func filterRecordNotFound(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "record not found"), ErrRecordNotFound
}
