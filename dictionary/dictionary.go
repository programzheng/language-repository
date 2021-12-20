package dictionary

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/programzheng/language-repository/orm"
	"github.com/programzheng/language-repository/user"
	"gorm.io/gorm"
)

type Dictionary struct {
	gorm.Model
	UserID             uint               `json:"user_id"`
	DictionaryLanguage DictionaryLanguage `json:"language"`
}

type DictionaryLanguage struct {
	gorm.Model
	DictionaryID uint   `json:"dictionary_id"`
	English      string `gorm:"size:255" json:"english" validate:"required"`
	Chinese      string `gorm:"size:255" json:"chinese" validate:"required"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func Setup() {
	orm.SetupTableModel(&Dictionary{})
	orm.SetupTableModel(&DictionaryLanguage{})
}

func defaultWhere(c *fiber.Ctx, tx *gorm.DB) *gorm.DB {
	tx = tx.Where("user_id IN ?", []uint{0})
	user := user.GetUserByJwtToken(c)
	//where current user's
	if user != nil {
		tx = tx.Or("user_id IN ?", []uint{user.ID})
	}
	return tx
}

func ValidateStruct(input interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(input)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func GetDictionaries(c *fiber.Ctx) error {
	//search text
	search := c.Query("search")

	tx := orm.GetDB().Joins("DictionaryLanguage")

	//external search
	if search != "" {
		tx = tx.Where("DictionaryLanguage.english LIKE ?", search+"%").Or("DictionaryLanguage.chinese LIKE ?", search+"%")
	}

	tx = defaultWhere(c, tx)

	dictionaries := new([]Dictionary)
	result := tx.Find(dictionaries)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(dictionaries)
}

func GetDictionary(c *fiber.Ctx) error {
	id := c.Params("id")

	dictionary := new(Dictionary)
	result := orm.GetDB().Preload("DictionaryLanguage").Find(dictionary, id)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(dictionary)
}

func NewDictionary(c *fiber.Ctx) error {
	dictionaryLanguage := new(DictionaryLanguage)
	if err := c.BodyParser(dictionaryLanguage); err != nil {
		return c.Status(503).SendString(err.Error())
	}
	dictionary := new(Dictionary)
	dictionary.DictionaryLanguage = *dictionaryLanguage

	//validate the dictionary
	errors := ValidateStruct(*dictionary)
	if errors != nil {
		return c.JSON(errors)
	}

	dictionary.UserID = 0
	user := user.GetUserByJwtToken(c)

	if user != nil {
		dictionary.UserID = user.ID
	}

	result := orm.GetDB().Create(&dictionary)
	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"results": dictionary,
	})
}

func UpdateDictionary(c *fiber.Ctx) error {
	id := c.Params("id")

	dictionary := new(Dictionary)
	tx := defaultWhere(c, orm.GetDB())
	tx.Find(&dictionary, id)
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
	tx := defaultWhere(c, orm.GetDB())
	tx.Find(&dictionary, id)
	if dictionary.ID == 0 {
		return c.Status(500).SendString("No Found with ID")
	}
	orm.GetDB().Delete(&dictionary)

	return c.SendString("Successfully deleted")
}
