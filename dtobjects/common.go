package dtobjects

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

type BaseDto struct {
	Success  bool     `json:"success"`
	Messages []string `json:"messages"`
}

type ErrorDto struct {
	BaseDto
	Errors map[string]interface{} `json:"errors"`
}

func BadRequestDto(err error) ErrorDto {
	res := ErrorDto{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	fmt.Println("error on validation --->", errs)
	res.Messages = make([]string, len(errs))
	count := 0
	for _, v := range errs {
		log.Println("----=====>", v.ActualTag())
		if v.ActualTag() == "required" {
			var message = fmt.Sprintf("%v is required", v.Field())
			res.Errors[v.Field()] = message
			res.Messages[count] = message
		} else {
			var message = fmt.Sprintf("%v has to be %v", v.Field(), v.ActualTag())
			res.Errors[v.Field()] = message
			res.Messages = append(res.Messages, message)
		}
		count++
	}
	return res
}

func DetailedErrors(key string, err error) map[string]interface{} {
	return map[string]interface{}{
		"success":  false,
		"messages": []string{fmt.Sprintf("s -> %v", err.Error())},
		"errors":   err,
	}
}

func CreateSuccessDto(result map[string]interface{}) map[string]interface{} {
	result["success"] = true
	return result
}
