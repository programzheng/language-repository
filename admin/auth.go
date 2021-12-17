package admin

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
	admin := new(Admin)
	if err := c.BodyParser(admin); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	inputPassword := admin.Password

	//get admin
	result := orm.GetDB().Where("account = ?", admin.Account).First(&admin)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "找不到使用者",
		})
	}
	hashPassword := admin.Password

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
		"account": admin.Account,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_ADMIN_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func GetAdminByJwtToken(c *fiber.Ctx) *Admin {
	localsAdmin := c.Locals("admin")
	if localsAdmin == nil {
		return nil
	}
	jwtAdmin := localsAdmin.(*jwt.Token)
	claims := jwtAdmin.Claims.(jwt.MapClaims)
	account := claims["account"].(string)

	admin := new(Admin)
	//get admin
	result := orm.GetDB().Where("account = ?", account).First(&admin)
	if result.Error != nil {
		return nil
	}

	return admin
}
