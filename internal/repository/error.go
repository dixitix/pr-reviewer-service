// Package repository содержит интерфейсы и ошибки доступа к хранилищу данных
package repository

import "errors"

// Ошибки, которые могут возвращать репозитории.
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)
