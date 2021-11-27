package dictionary

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/programzheng/language-repository/orm"
	"gorm.io/gorm"
)

type Dictionary struct {
	gorm.Model
	MemberID uint
	English  string `json:"english"`
	Chinese  string `json:"chinese"`
}

func setup() {
	orm.SetupTableModel(&Dictionary{})
}

func GetDictionaries(c *fiber.Ctx) error {
	dictionary := new([]Dictionary)
	result := orm.GetDB().Find(dictionary)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(dictionary)
}

func GetDictionary(c *fiber.Ctx) error {
	id := c.Params("id")

	dictionary := new(Dictionary)
	result := orm.GetDB().Find(dictionary, id)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(dictionary)
}

func NewDictionary(c *fiber.Ctx) error {
	setup()

	dictionary := new(Dictionary)
	if err := c.BodyParser(dictionary); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	dictionary.MemberID = 0

	result := orm.GetDB().Create(&dictionary)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(dictionary)
}

func UpdateDictionary(c *fiber.Ctx) error {
	id := c.Params("id")

	dictionary := new(Dictionary)
	orm.GetDB().Find(&dictionary, id)
	if dictionary.ID == 0 {
		return c.Status(500).SendString("No Found with ID")
	}

	if err := c.BodyParser(dictionary); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	result := orm.GetDB().Save(&dictionary)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(dictionary)
}

func DeleteDictionary(c *fiber.Ctx) error {
	id := c.Params("id")

	dictionary := new(Dictionary)
	orm.GetDB().Find(&dictionary, id)
	if dictionary.ID == 0 {
		return c.Status(500).SendString("No Found with ID")
	}
	orm.GetDB().Delete(&dictionary)

	return c.SendString("Successfully deleted")
}
