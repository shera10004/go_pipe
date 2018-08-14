package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	fmt.Println("echo framework test...")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hello, echo framework^^")
	})

	e.GET("/users/:id", func(c echo.Context) error {
		//Param : /users/value(id)
		id := c.Param("id")
		return c.String(http.StatusOK, id)
	})

	e.GET("/users", func(c echo.Context) error {
		result := fmt.Sprintf("Header:%v ", c.Request().Header.Get("Content-Type"))

		if len(c.ParamNames()) == 0 && len(c.QueryParams()) == 0 {
			return c.String(http.StatusOK, "input your id")
		}

		//QueryParam : /users?name=aValue&age=bValue
		name := c.QueryParam("name")
		age := c.QueryParam("age")

		result = fmt.Sprintf("%s \nname:%v, age:%v", result, name, age)
		return c.String(http.StatusOK, result) //"name:"+name+", age:"+age)
	})

	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
			return
		}
		/*
			if err = c.Validate(u); err != nil {
				return
			}
		*/
		return c.JSON(http.StatusOK, u)
	})

	e.Logger.Fatal(e.Start(":1323"))

}

type (
	User struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	CustomValidator struct {
		//validator *validator.Validate

	}
)

/*
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
//*/
