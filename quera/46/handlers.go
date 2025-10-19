package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Message string `json:"message"`
}

type user struct {
	name string
	age  int
	job  string
}

var users []user

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func makeName(fName string, lName string) string {
	return fName + "_" + lName
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
	name := makeName(fName, lName)
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
		c.JSON(http.StatusBadRequest, Result{"firtname is required"}) // توجه کنید: firtname NOT firstname
		return
	}
	if lName == "" {
		c.JSON(http.StatusBadRequest, Result{"lastname is required"})
		return
	}
	if !isInt(ageStr) {
		c.JSON(http.StatusBadRequest, Result{"age should be integer"})
		return
	}

	age, _ := strconv.Atoi(ageStr)

	if addUser(fName, lName, job, age) {
		c.JSON(http.StatusOK, Result{fmt.Sprintf("%s %s registered successfully", fName, lName)})
	} else {
		c.JSON(http.StatusConflict, Result{fmt.Sprintf("%s %s registered before", fName, lName)})
	}
}

func hello(c *gin.Context) {
	fName := c.Param("firstname")
	lName := c.Param("lastname")

	if fName == "" {
		c.JSON(http.StatusBadRequest, Result{"firstname is required"})
		return
	}
	if lName == "" {
		c.JSON(http.StatusBadRequest, Result{"lastname is required"})
		return
	}

	name := makeName(fName, lName)

	for _, u := range users {
		if u.name == name {
			c.String(http.StatusOK, fmt.Sprintf("Hello %s %s; Job: %s; Age: %d", fName, lName, u.job, u.age))
			return
		}
	}

	c.String(http.StatusNotFound, fmt.Sprintf("%s %s is not registered", fName, lName))
}
