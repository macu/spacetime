package recaptcha

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"treetime/pkg/env"
	"treetime/pkg/utils/net"
	"treetime/pkg/utils/secrets"
)

func VerifyRecaptcha(r *http.Request) (bool, error) {

	var token = r.FormValue("g-recaptcha-response")

	if token == "" {
		if env.IsLocal() {
			// Allow bypass on local
			return true, nil
		}
		return false, fmt.Errorf("reCAPTCHA token undefined")
	}

	var err error

	var recaptchaSecretKey = env.GetRecaptchaSiteKey()
	if recaptchaSecretKey == "" {
		if env.IsAppEngine() {
			recaptchaSecretKey, err = secrets.LoadSecret(os.Getenv("RECAPTCHA_SECRET"))
			if err != nil {
				return false, fmt.Errorf("loading reCAPTCHA secret key: %w", err)
			}
		} else {
			return false, fmt.Errorf("reCAPTCHA secret key undefined")
		}
	}

	ip := net.GetUserIP(r)

	form := url.Values{}
	form.Set("secret", recaptchaSecretKey)
	form.Set("response", token)
	form.Set("remoteip", ip)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.PostForm("https://www.google.com/recaptcha/api/siteverify", form)
	if err != nil {
		return false, fmt.Errorf("error fetching URL: %w", err)
	}

	defer res.Body.Close()

	var response = struct {
		Success bool `json:"success"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return false, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Success, nil
}
