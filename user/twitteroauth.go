package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type TwitterOauth10RequestTokenRequest struct {
	OauthCallback   string `json:"oauth_callback"`
	XAuthAccessType string `json:"x_auth_access_type"`
}

type TwitterOauth10RequestTokenResponse struct {
	Status string                 `json:"status"`
	Value  map[string]interface{} `json:"value"`
}

func TwitterOauth10RequestToken(c *fiber.Ctx) error {
	twitterOauthRequestTokenRequest := new(TwitterOauth10RequestTokenRequest)
	if err := c.BodyParser(twitterOauthRequestTokenRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	TwitterOauthRequestTokenResponse, err := sendTwitterOauth10RequestTokenRequest(twitterOauthRequestTokenRequest)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}
	fmt.Printf("%v", TwitterOauthRequestTokenResponse)

	return nil
}

func sendTwitterOauth10RequestTokenRequest(twitterOauthRequestTokenRequest *TwitterOauth10RequestTokenRequest) (map[string]interface{}, error) {
	client := resty.New()

	result := TwitterOauth10RequestTokenResponse{}
	resp, err := client.R().
		SetBody(twitterOauthRequestTokenRequest).
		SetResult(&result).
		Post(os.Getenv("TWITTER_OAUTH_API_REQUEST_TOKEN"))

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

func TwitterOauth20Token(c *fiber.Ctx) {

}
