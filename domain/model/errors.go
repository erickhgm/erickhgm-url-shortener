package model

import "fmt"

type DocumentAlreadyExistsError struct {
	Id  string
	Url string
}

func (e *DocumentAlreadyExistsError) Error() string {
	return fmt.Sprintf("Document Id %v already exists on database: %v", e.Id, e.Url)
}

type DocumentNotFoundError struct {
	Id string
}

func (e *DocumentNotFoundError) Error() string {
	return fmt.Sprintf("Document Id %v not found", e.Id)
}

type InvalidUrlError struct {
	Messsage string
}

func (e *InvalidUrlError) Error() string {
	return fmt.Sprintf(e.Messsage)
}
