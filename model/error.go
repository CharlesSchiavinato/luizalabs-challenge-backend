package model

import "fmt"

type Error struct {
	// Código do Erro
	Code float32 `json:"code" validate:"required" example:"400.1" format:"float"`
	// Descrição do Erro
	Message string `json:"message" validate:"required" example:"error message"`
}

func BadRequestDeserialize(controllerTitle string) *Error {
	return &Error{
		Code:    400.1,
		Message: fmt.Sprintf("Error deserializing %s from body", controllerTitle),
	}
}

func BadRequestModelValidate(controllerTitle, message string) *Error {
	return &Error{
		Code:    400.2,
		Message: fmt.Sprintf("Error validating %s: %s", controllerTitle, message),
	}
}

func BadRequestParamValidate(message string) *Error {
	return &Error{
		Code:    400.3,
		Message: fmt.Sprintf("Error validating parameters: %s", message),
	}
}

func BadRequestRepositoryPersist(controllerTitle, message string) *Error {
	return &Error{
		Code:    400.4,
		Message: fmt.Sprintf("Error persisting %s in repository: %s", controllerTitle, message),
	}
}

func BadRequestFormParsing() *Error {
	return &Error{
		Code:    400.5,
		Message: "Error parsing form",
	}
}

func BadRequestRetrievingFile() *Error {
	return &Error{
		Code:    400.6,
		Message: "Error retrieving file",
	}
}

func BadRequestFileType() *Error {
	return &Error{
		Code:    400.7,
		Message: "Error invalid file type",
	}
}

func BadRequestFileRecordValidate(message string) *Error {
	return &Error{
		Code:    400.8,
		Message: message,
	}
}

func NotFound(controllerTitle string) *Error {
	return &Error{
		Code:    404.1,
		Message: fmt.Sprintf("%s not found", controllerTitle),
	}
}

func InternalServerErrorGeneral(message string) *Error {
	return &Error{
		Code:    500.1,
		Message: message,
	}
}

func InternalServerErrorRepositoryPersist(controllerTitle string) *Error {
	return &Error{
		Code:    500.2,
		Message: fmt.Sprintf("Error persisting %s in repository", controllerTitle),
	}
}

func InternalServerErrorRepositoryLoad(controllerTitle string) *Error {
	return &Error{
		Code:    500.3,
		Message: fmt.Sprintf("Error loading %s from repository", controllerTitle),
	}
}
