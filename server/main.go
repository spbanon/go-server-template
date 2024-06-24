package main

import (
	"net/http"
	"time"

	"jwt_fiber_template/dbhandler"
	"jwt_fiber_template/rabbitmq"
	_ "jwt_fiber_template/server/docs"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v4/pgxpool"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @host 127.0.0.1:1337
// @BasePath /

//	@tag.name			login
//	@tag.description	Request ID

//	@tag.name			robots.txt
//	@tag.description	Return ok

var tasksSender *rabbitmq.FileSender
var dbHandler *pgxpool.Pool

func getRobotsTxt(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}

func login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")

	if user != "john" || pass != "doe" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := jwt.MapClaims{
		"name":  "John Doe",
		"admin": true,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func accessible(c *fiber.Ctx) error {
	return c.SendString("Accessible")
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func main() {

	var err error
	if !fiber.IsChild() {
		for {
			tasksSender, err = rabbitmq.NewFileSender("127.0.0.1", "2222", "user", "password", "test")
			if err == nil {
				break
			}
		}
		defer tasksSender.Close()
		for {
			dbHandler, err = dbhandler.NewDBConnection("127.0.0.1", "1111", "user", "password", "db", true)
			if err == nil {
				break
			}
		}
		defer dbHandler.Close()
	}

	app := fiber.New()
	app.Post("/login", login)
	app.Get("/", accessible)
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	app.Get("/restricted", restricted)
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/robots.txt", getRobotsTxt)
	app.Listen("127.0.0.1:1337")

}
