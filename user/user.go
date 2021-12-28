package user

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/programzheng/language-repository/orm"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Account     string      `gorm:"size:255; unique"`
	Password    string      `gorm:"size:255" json:"password"`
	UserProfile UserProfile `json:"profile"`
}

type UserProfile struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	Email  string `gorm:"size:255" json:"email" validate:"email"`
}

func Setup() {
	orm.SetupTableModel(&User{})
	orm.SetupTableModel(&UserProfile{})
}

func createHash(secret string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func setAccount(user *User) *User {
	accountType := os.Getenv("USER_ACCOUNT_TYPE")
	switch accountType {
	case "email":
		user.Account = user.UserProfile.Email
	}

	return user
}

func NewUser(c *fiber.Ctx) error {
	userProfile := new(UserProfile)
	if err := c.BodyParser(userProfile); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	user.UserProfile = *userProfile

	user = setAccount(user)
	if user.Account == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "帳號格式有誤",
		})
	}

	//check account
	orm.GetDB().Where("account = ?", user.Account).Find(user)
	if user.ID > 0 {
		return c.JSON(&fiber.Map{
			"status":  "fail",
			"message": "帳號重複，請再確認",
		})
	}

	//crypt password
	user.Password = createHash(user.Password)

	result := orm.GetDB().Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": "建立失敗",
		})
	}

	return c.JSON(user)
}
