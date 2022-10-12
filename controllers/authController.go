package controllers

import (
	// "fmt"

	"login/database"
	"login/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	user_password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	var exists bool = false

	if err := database.DB.Raw(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)",
		data["email"]).
		Scan(&exists).Error; err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	if exists {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	user := models.User{
		Name:  data["name"],
		Email: data["email"],
	}

	database.DB.Create(&user)

	password := models.Password{
		UserId:   user.Id,
		Password: user_password,
	}

	database.DB.Preload("passwords").Create(&password)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	var pwd models.Password

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		//c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	database.DB.Select("password").Where("user_id = ?", user.Id).First(&pwd)

	if err := bcrypt.CompareHashAndPassword(pwd.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), //1 day
	})

	token, err := claims.SignedString([]byte(SecretKey))

	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func Delete(c *fiber.Ctx) error {

	id := c.Params("id")

	var user models.User

	var pwd models.Password

	database.DB.Select("id").Where("id = ?", id).First(&user)
	database.DB.Delete(&user)

	database.DB.Select("user_id").Where("user_id = ?", id).First(&pwd)
	database.DB.Delete(&pwd)

	return c.JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}

func User(c *fiber.Ctx) error {

	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Please Login to continue",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	var pwd models.Password

	database.DB.Select("id", "name", "email").Where("id = ?", claims.Issuer).First(&user)

	database.DB.Select("password").Where("id = ?", claims.Issuer).First(&pwd)

	return c.JSON(fiber.Map{
		"id":    claims.Issuer,
		"name":  user.Name,
		"email": user.Email,
	})
}

func Update(c *fiber.Ctx) error {

	id := c.Params("id")

	var user models.User

	//var pwd models.Password

	var data map[string]string

	//database.DB.First(&data, id)

	if err := c.BodyParser(&user); err != nil {
		return err
	}
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	database.DB.Model(&user).Debug().Select("*").Where("id = ?", id).Updates(&user)

	user_password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	database.DB.Debug().Select("password").Where("user_id = ?", id).Updates(models.Password{Password: user_password})

	return c.JSON(fiber.Map{
		"id":       id,
		"name":     data["name"],
		"email":    data["email"],
		"password": data["password"],
	})

}

func GetUser(c *fiber.Ctx) error {

	id := c.Params("id")

	var user models.User

	database.DB.Select("id", "name", "email").First(&user, id)

	return c.JSON(user)
}
