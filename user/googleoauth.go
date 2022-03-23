package user

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/programzheng/language-repository/orm"
)

const provider string = "google_oauth"

type GoogleOauthRequest struct {
	IDToken string `json:"id_token"`
}

type GoogleGetGoogleUserInfoByIDTokenResponse struct {
	Status string                 `json:"status"`
	Value  map[string]interface{} `json:"value"`
}

func getGoogleUserInfoByIDToken(IDToken string) (map[string]interface{}, error) {
	client := resty.New()

	result := GoogleGetGoogleUserInfoByIDTokenResponse{}
	resp, err := client.R().
		SetBody(map[string]interface{}{"id_token": IDToken}).
		SetResult(&result).
		Post(os.Getenv("GOOGLE_OAUTH_API_GET_UNIQUE_ID_URL"))

	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != fiber.StatusOK {
		t := make(map[string]interface{})
		t["error"] = "unknow error"
		err := json.Unmarshal([]byte(resp.Body()), &t)
		if err != nil {
			log.Printf("%v", err)
		}
		return nil, errors.New(t["error"].(string))
	}

	return result.Value, nil
}

func GoogleOauth(c *fiber.Ctx) error {
	googleOauthRequest := new(GoogleOauthRequest)
	if err := c.BodyParser(googleOauthRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	userInfo, err := getGoogleUserInfoByIDToken(googleOauthRequest.IDToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	uniqueID := userInfo["unique_id"].(string)

	u, err := GetUserByProvider(provider, uniqueID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if u.ID == 0 {
		u, err = NewOauthUser(provider, uniqueID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"status":  "fail",
				"message": err.Error(),
			})
		}
	}
	go UpdateGoogleUserInfo(u, userInfo["claims"].(map[string]interface{}))

	t, err := generateJwt(u)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{"token": t})
}

func UpdateGoogleUserInfo(u *User, claims map[string]interface{}) error {
	up := u.UserProfile
	update := make(map[string]interface{})
	if email, ok := claims["email"].(string); ok {
		update["email"] = email
	}
	result := orm.GetDB().Model(up).Where("user_id = ?", u.ID).Updates(update)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
