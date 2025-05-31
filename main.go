package main

import (
	"fmt"
	"net/http"

	"os"

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

var (
	app         = kingpin.New("App", "Simple app")
	flagAppName = app.Flag("name", "Application name").Required().String()
	flagPort    = app.Flag("port", "Web server port").Short('p').Default("9000").Int()
)

var (
	commandAdd             = app.Command("add", "add new user")
	commandAddFlagOverride = commandAdd.Flag("override", "override existing user").Short('o').Bool()
	commandAddArgUser      = commandAdd.Arg("user", "username").Required().String()
)

var (
	commandUpdate           = app.Command("update", "update user")
	commandUpdateArgOldUser = commandUpdate.Arg("old", "old username").Required().String()
	commandUpdateArgNewUser = commandUpdate.Arg("new", "new username").Required().String()
)

var (
	commandDelete          = app.Command("delete", "delete user")
	commandDeleteFlagForce = commandDelete.Flag("force", "force deletion").Short('f').Bool()
	commandDeleteArgUser   = commandDelete.Arg("user", "username").Required().String()
)

func main() {
	commandAdd.Action(func(ctx *kingpin.ParseContext) error {
		// more code here ...
		user := *commandAddArgUser
		override := *commandAddFlagOverride
		fmt.Printf("adding user %s, override %t \n", user, override)

		return nil
	})

	commandUpdate.Action(func(ctx *kingpin.ParseContext) error {
		// more code here ...
		oldUser := *commandUpdateArgOldUser
		newUser := *commandUpdateArgNewUser
		fmt.Printf("updating user from %s %s \n", oldUser, newUser)

		return nil
	})

	commandDelete.Action(func(ctx *kingpin.ParseContext) error {
		// more code here ...
		user := *commandDeleteArgUser
		force := *commandDeleteFlagForce
		fmt.Printf("deleting user %s, force %t \n", user, force)

		return nil
	})
	kingpin.MustParse(app.Parse(os.Args[1:]))
	appName := *flagAppName
	port := fmt.Sprintf(":%d", *flagPort)

	fmt.Printf("Starting %s at %s", appName, port)
	r := echo.New()

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
