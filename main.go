package main

import (
	"net/http"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	v "github.com/go-playground/validator/v10"
)

// Booking contains binded and validated data.
type Account struct {
	Username string `form:"username" binding:"required,max=16,usernamevalid"`
	Password string `form:"password" binding:"required,max=10,passwordvalid"`
}

func main() {
	route := gin.Default()
	_ = New()
	route.POST("/create", GetAccount)
	route.Run(":8085")
}

func GetAccount(c *gin.Context) {
	var b Account
	if err := c.ShouldBindWith(&b, binding.Form); err == nil {
		c.JSON(http.StatusOK, b)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

func FormTagName(fld reflect.StructField) string {
	name := fld.Tag.Get("form")
	if name != "" {
		return name
	}
	name = fld.Tag.Get("json")
	if name != "" {
		return name
	}
	return fld.Name
}

var IsUsernameValid = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
var IsPasswordValid = regexp.MustCompile(`^[a-zA-Z]+[0-9]+[#?!@$%^&*-]+$`).MatchString

const (
	passwordvalid = "passwordvalid"
	usernamevalid = "usernamevalid"
)

var tagFuncMaps = map[string]func(v.FieldLevel) bool{
	passwordvalid: PasswordValidated,
	usernamevalid: UsernameValidated,
}

var tagNameFunc = []func(reflect.StructField) string{
	FormTagName,
}

func New() *v.Validate {
	validate, ok := binding.Validator.Engine().(*v.Validate)
	if !ok {
		return nil
	}

	for key, value := range tagFuncMaps {
		validate.RegisterValidation(key, value)
	}

	for _, value := range tagNameFunc {
		validate.RegisterTagNameFunc(value)
	}

	return validate
}

func PasswordValidated(fl v.FieldLevel) bool {
	field := fl.Field().String()
	if len(field) < 6 && len(field) > 12 {
		return false
	}
	if !IsPasswordValid(field) {
		return false
	}
	return !(field == "")
}

func UsernameValidated(fl v.FieldLevel) bool {
	field := fl.Field().String()
	if !IsUsernameValid(field) {
		return false
	}
	return !(field == "")
}
