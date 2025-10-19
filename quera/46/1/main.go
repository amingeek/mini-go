package main

import (
	"fmt"
	"net/http"
	"strconv"
	"unicode"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type user struct {
	name string
	age  int
	job  string
}

var users []user

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func search(name string) bool {
	for _, u := range users {
		if u.name == name {
			return true
		}
	}
	return false
}

func addUser(fName, lName, job string, age int) bool {
	name := fName + "_" + lName
	if search(name) {
		return false
	}
	users = append(users, user{name: name, age: age, job: job})
	return true
}

func register(c *gin.Context) {
	fName := c.PostForm("firstname")
	lName := c.PostForm("lastname")
	job := c.DefaultPostForm("job", "Unknown")
	ageStr := c.DefaultPostForm("age", "18")

	if fName == "" {
		c.JSON(http.StatusBadRequest, Result{http.StatusBadRequest, "firstname is required"})
		return
	}
	if lName == "" {
		c.JSON(http.StatusBadRequest, Result{http.StatusBadRequest, "lastname is required"})
		return
	}
	if !isInt(ageStr) {
		c.JSON(http.StatusBadRequest, Result{http.StatusBadRequest, "age should be integer"})
		return
	}

	age, _ := strconv.Atoi(ageStr)

	if addUser(fName, lName, job, age) {
		c.JSON(http.StatusOK, Result{http.StatusOK, fmt.Sprintf("%s %s registered successfully", fName, lName)})
	} else {
		c.JSON(http.StatusConflict, Result{http.StatusConflict, fmt.Sprintf("%s %s registered before", fName, lName)})
	}
}

func hello(c *gin.Context) {
	fName := c.PostForm("firstname")
	lName := c.PostForm("lastname")
	if addUser(fName, lName, "Unknown", 18) {
		c.JSON(http.StatusOK, Result{http.StatusOK, fmt.Sprintf("Hello %s %s; Job: %job%; Age: %age")})
	}
}
