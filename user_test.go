package main

import (
	"net/http"

	"github.com/gavv/httpexpect/v2"
)

//  referenced doc:https://github.com/gavv/httpexpect/blob/master/_examples/echo_test.go
func UserTest(e *httpexpect.Expect) {
	//create user
	user1 := e.POST("/users").
		WithJSON(map[string]string{
			"username": "sample_user",
		}).Expect().Status(http.StatusCreated).JSON().Object()

	// check
	user1.ValueEqual("username", "test_user")
	user1.Keys().ContainsOnly("ID", "username")
}
