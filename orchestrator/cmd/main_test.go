package main

import (
	"fmt"
	. "github.com/Eun/go-hit"
	"net/http"
	"testing"
	"time"
)

func TestExpressions(t *testing.T) {
	var id int64
	Test(
		t,
		Description("Registrate"),
		Post("http://localhost:8080/register"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{
			"login":    "test1",
			"password": "test",
		}),
		Expect().Status().Equal(http.StatusOK),
	)

	var token string
	Test(t,
		Description("Login"),
		Post("http://localhost:8080/login"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{
			"login":    "test1",
			"password": "test",
		}),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".token").In(&token))

	Test(t,
		Description("Post expression"),
		Post("http://localhost:8080/expression"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Authorization").Add(token),
		Send().Body().JSON(map[string]interface{}{
			"expression": "2 + 2 * 2",
			"settings": map[string]interface{}{
				"plus":  1,
				"minus": 1,
				"mult":  1,
				"div":   1,
			},
		}),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".id").In(&id))

	time.Sleep(5 * time.Second)
	var result string
	Test(t,
		Description("Post expression"),
		Get("http://localhost:8080/expression?id="+fmt.Sprintf("%d", id)),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("Authorization").Add(token),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".result").In(&result))
	if result != "6.00" {
		t.Fatal("Incorrect")
	}
}
