package main

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin/v2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type M map[string]interface{}

var ActionIndex = func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("from action index"))
}

var ActionHome = http.HandlerFunc(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("from action home"))
	},
)

var ActionAbout = echo.WrapHandler(
	http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("from action about"))
		},
	),
)

func middlewareOne(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fmt.Println("From middleware One")
		return next(ctx)
	}
}

func middlewareSomething(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("from middleware something")
		next.ServeHTTP(w, r)
	})
}

var (
	argAppName = kingpin.Arg("name", "Application name").Required().String()
	argPort    = kingpin.Arg("port", "Web server port").Default("9000").Int()
)

func main() {
	kingpin.Parse()
	appName := *argAppName
	port := fmt.Sprintf(":%d", *argPort)

	fmt.Printf("Starting %s at %s", appName, port)
	r := echo.New()

	r.Use(middlewareOne)
	r.Use(echo.WrapMiddleware(middlewareSomething))

	r.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	r.GET("/index", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, true)
	})

	r.GET("/index", echo.WrapHandler(http.HandlerFunc(ActionIndex)))
	r.GET("/home", echo.WrapHandler(ActionHome))
	r.GET("/about", ActionAbout)

	r.Logger.Fatal(r.Start(port))
}
