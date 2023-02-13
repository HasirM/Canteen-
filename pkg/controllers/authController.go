package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hasirm/vaddiapp/auth/pkg/database"
	"github.com/hasirm/vaddiapp/auth/pkg/models"
	"golang.org/x/crypto/bcrypt"
	// jwtware "github.com/gofiber/jwt/v3"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		panic(err)
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	// admin, err := strconv.ParseBool(data["admin"])
	// if err != nil {
	// 	fmt.Println("Something wrong with parsing string to bool", err)
	// }

	user := models.User{
		Name:     data["name"],
		CID:    data["cid"],
		Password: password,
		Phone:    data["phone"],
	}

	// if result.Error != nil {
	// 	return c.JSON(fiber.Map{
	// 		"message": "an error occured please try again",
	// 	})
	// }

	if err := database.DB.Where("name = ?", data["name"]).First(&user).Error; err == nil {
		return c.JSON(fiber.Map{
			"message": "username already exists",
		})
	}

	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err == nil {
		return c.JSON(fiber.Map{
			"message": "College ID already exists",
		})
	}

	if err := database.DB.Where("phone = ?", data["phone"]).First(&user).Error; err == nil {
		return c.JSON(fiber.Map{
			"message": "phone number already exists",
		})
	}

	database.DB.Create(&user)

	return c.JSON(fiber.Map{
		"":        "Registered Successfully",
		"message": user,
	})
}

func createAJWTandRJWT(c *fiber.Ctx) error {

	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Second * 10).Unix(), //3 days

	})
	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(), //3 days

	})

	accessToken, err := accessClaims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not Login",
		})
	}

	refreshToken, err := refreshClaims.SignedString([]byte(SecretKey))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "Could not Login",
		})
	}

	accessCookie := fiber.Cookie{
		Name:     "ajwt",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Second * 10),
		HTTPOnly: true,
	}
	refreshCookie := fiber.Cookie{
		Name:     "rjwt",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Minute * 1),
		HTTPOnly: true,
	}

	c.Cookie(&accessCookie)
	c.Cookie(&refreshCookie)

	return nil

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("cid = ?", data["cid"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Incorrect Password",
		})
	}

	createAJWTandRJWT(c)

	return c.JSON(fiber.Map{
		"message": "Successfully Logged In",
	})

}

func Logout(c *fiber.Ctx) error {
	acookie := fiber.Cookie{
		Name:     "ajwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), //expired one hour ago
		HTTPOnly: true,
	}

	rcookie := fiber.Cookie{
		Name:     "rjwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), //expired one hour ago
		HTTPOnly: true,
	}

	c.Cookie(&acookie)
	c.Cookie(&rcookie)

	return c.JSON(fiber.Map{
		"message": "Successfully Logged Out",
	})
}

func Refresh(c *fiber.Ctx) error {
	rcookie := c.Cookies("rjwt")

	_, rErr := jwt.ParseWithClaims(rcookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if rErr != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated, Please re-login",
		})
	} else {
		createAJWTandRJWT(c)
	}

	return c.JSON(fiber.Map{
		"message": "JWT Tokens Updated",
	})
}


func Product(c *fiber.Ctx) error {
	cookie := c.Cookies("ajwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated, User does not have access to this site",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("ID = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}


func Order(c *fiber.Ctx) error {
	cookie := c.Cookies("ajwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated, User does not have access to this site",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("ID = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Wallet(c *fiber.Ctx) error {
	cookie := c.Cookies("ajwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated, User does not have access to this site",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	database.DB.Where("ID = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

// func Accessible(c *fiber.Ctx) error {
// 	return c.SendString("Unauthenticated Page")
// }

// func Restricted(c *fiber.Ctx) error {
// 	return c.SendString("Welcome User")
// }

// func Home(c *fiber.Ctx) error {
// 	return c.SendString("Unauthenticated Page")
// }

// func Home1(c *fiber.Ctx) error {
// 	return c.SendString(" Welcome User")
// }

// func Home2(c *fiber.Ctx) error {
// 	return c.SendString(" Welcome User")
// }

// func Home3(c *fiber.Ctx) error {
// 	return c.SendString(" Welcome User")
// }
