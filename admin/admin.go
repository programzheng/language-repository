package admin

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/programzheng/language-repository/orm"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Account  string `gorm:"size:255; unique" json:"account"`
	Password string `gorm:"size:255" json:"password"`
}

func setup() {
	orm.SetupTableModel(&Admin{})
}

func createHash(secret string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func NewAdmin(c *fiber.Ctx) error {
	setup()

	admin := new(Admin)
	if err := c.BodyParser(admin); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	//crypt password
	admin.Password = createHash(admin.Password)

	result := orm.GetDB().Create(&admin)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "建立失敗",
		})
	}

	return c.JSON(admin)
}
