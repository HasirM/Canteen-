package middleware

import (
	// "fmt"
	// "time"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/hasirm/vaddiapp/auth/pkg/config"
	"github.com/hasirm/vaddiapp/auth/pkg/database"
	"github.com/hasirm/vaddiapp/auth/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// var SecretKey = databasetabasedatabase.GetENV("SECRET")
var SecretKey = config.Data().JWTKey

func Validate(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("cid = ?", data["cid"]).First(&user)

	if user.ID == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized!!",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized!!",
		})
	}

	if !user.IsAdmin {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthorized the user doesnt have access to do this operation! Contact admin for more info",
		})
	}
	// createAJWTandRJWT(c)

	c.Status(fiber.StatusAccepted)
	return c.JSON(fiber.Map{
		"message": "Success !!",
	})

}

// func createAJWTandRJWT(c *fiber.Ctx) error{

// 	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
// 		IssuedAt:  time.Now().Unix(),
// 		ExpiresAt: time.Now().Add(time.Second * 10).Unix(), //3 days

// 	})
// 	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
// 		IssuedAt:  time.Now().Unix(),
// 		ExpiresAt: time.Now().Add(time.Minute * 1).Unix(), //3 days

// 	})

// 	accessToken, err := accessClaims.SignedString([]byte(SecretKey))
// 	if err != nil {
// 		c.Status(fiber.StatusInternalServerError)
// 		return c.JSON(fiber.Map{
// 			"message": "Could not Login",
// 		})
// 	}

// 	refreshToken, err := refreshClaims.SignedString([]byte(SecretKey))
// 	if err != nil {
// 		c.Status(fiber.StatusInternalServerError)
// 		return c.JSON(fiber.Map{
// 			"message": "Could not Login",
// 		})
// 	}

// 	accessCookie := fiber.Cookie{
// 		Name:     "ajwt",
// 		Value:    accessToken,
// 		Expires:  time.Now().Add(time.Second * 10),
// 		HTTPOnly: true,
// 	}
// 	refreshCookie := fiber.Cookie{
// 		Name:     "rjwt",
// 		Value:    refreshToken,
// 		Expires:  time.Now().Add(time.Minute * 1),
// 		HTTPOnly: true,
// 	}

// 	c.Cookie(&accessCookie)
// 	c.Cookie(&refreshCookie)

// 	return c.JSON(fiber.Map{
// 		"message": "Success",
// 	})
// }
