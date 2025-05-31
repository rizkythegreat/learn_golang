package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/labstack/echo"
	"github.com/novalagung/gubrak/v2"
)

type M map[string]interface{}

var sc = securecookie.New([]byte("very-secret"), []byte("a-lot-secret-yay"))

func setCookie(c echo.Context, name string, data M) error {
	encoded, err := sc.Encode(name, data)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     name,
		Value:    encoded,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(1 * time.Hour),
	}
	http.SetCookie(c.Response(), cookie)

	return nil
}

func getCookie(c echo.Context, name string) (M, error) {
	cookie, err := c.Request().Cookie(name)
	if err == nil {
		data := M{}
		if err = sc.Decode(name, cookie.Value, &data); err == nil {
			return data, nil
		}
	}
	return nil, err
}

func main() {
	const CookieName = "data"
	e := echo.New()
	e.GET("/gubrak", func(c echo.Context) error {
		data, err := getCookie(c, CookieName)
		if err != nil && err != http.ErrNoCookie && err != securecookie.ErrMacInvalid {
			return err
		}

		if data == nil {
			data = M{"Message": "Hello", "ID": gubrak.RandomString(32)}

			err = setCookie(c, CookieName, data)
			if err != nil {
				return err
			}
		}

		return c.JSON(http.StatusOK, data)
	})
	confAppName := os.Getenv("APP_NAME")
	if confAppName == "" {
		e.Logger.Fatal("APP_NAME config is required")
	}
	confServerPort := os.Getenv("SERVER_PORT")
	if confServerPort == "" {
		e.Logger.Fatal("SERVER_PORT config is required")
	}
	e.GET("/index", func(c echo.Context) (err error) {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"appName": confAppName,
			"server": map[string]interface{}{
				"port": confServerPort,
			},
			"env": map[string]interface{}{
				"SERVER_READ_TIMEOUT_IN_MINUTE":  os.Getenv("SERVER_READ_TIMEOUT_IN_MINUTE"),
				"SERVER_WRITE_TIMEOUT_IN_MINUTE": os.Getenv("SERVER_WRITE_TIMEOUT_IN_MINUTE"),
			},
			"config": map[string]interface{}{
				"appName": os.Getenv("APP_NAME"),
				"server": map[string]interface{}{
					"port": os.Getenv("SERVER_PORT"),
				},
			},
		})
	})

	server := new(http.Server)
	server.Addr = ":" + confServerPort
	if confServerReadTimeout := os.Getenv("SERVER_READ_TIMEOUT_IN_MINUTE"); confServerReadTimeout != "" {
		duration, _ := strconv.Atoi(confServerReadTimeout)
		server.ReadTimeout = time.Duration(duration) * time.Minute
	}

	if confServerWriteTimeout := os.Getenv("SERVER_WRITE_TIMEOUT_IN_MINUTE"); confServerWriteTimeout != "" {
		duration, _ := strconv.Atoi(confServerWriteTimeout)
		server.WriteTimeout = time.Duration(duration) * time.Minute
	}
	e.Logger.Print("Starting", confAppName)
	e.Logger.Fatal(e.StartServer(server))
}
