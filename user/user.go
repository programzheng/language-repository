package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/programzheng/language-repository/orm"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Account  string `gorm:"size:255; unique" json:"account"`
	Password string `gorm:"size:255" json:"password"`
}

func setup() {
	orm.SetupTableModel(&User{})
}

func createHash(secret string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func NewUser(c *fiber.Ctx) error {
	setup()

	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	//crypt password
	user.Password = createHash(user.Password)

	result := orm.GetDB().Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "建立失敗",
		})
	}

	return c.JSON(user)
}
