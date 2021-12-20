package user

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/programzheng/language-repository/orm"
	"golang.org/x/crypto/bcrypt"
)

func checkHash(hash string, secret string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err
}

func Login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	inputPassword := user.Password

	//get user
	result := orm.GetDB().Where("account = ?", user.Account).First(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "找不到使用者",
		})
	}
	hashPassword := user.Password

	//check password
	err := checkHash(hashPassword, inputPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "密碼錯誤",
		})
	}

	//jwt
	// Create the Claims
	claims := jwt.MapClaims{
		"account": user.Account,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_USER_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func Auth(c *fiber.Ctx) error {
	user := GetUserByJwtToken(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "fail",
		})
	}
	return c.JSON(fiber.Map{
		"status": "success",
	})
}

func GetUserByJwtToken(c *fiber.Ctx) *User {
	localsUser := c.Locals("user")
	if localsUser == nil {
		return nil
	}
	jwtUser := localsUser.(*jwt.Token)
	claims := jwtUser.Claims.(jwt.MapClaims)
	account := claims["account"].(string)

	user := new(User)
	//get user
	result := orm.GetDB().Where("account = ?", account).First(&user)
	if result.Error != nil {
		return nil
	}

	return user
}
