package env

// this package holds environment variables

var isAppEngine bool
var cacheControlVersionStamp string
var recaptchaSiteKey string
var recaptchaSecret string
var mailjetApiKey string
var mailjetSecret string

func IsAppEngine() bool {
	return isAppEngine
}

func IsLocal() bool {
	return !isAppEngine
}

func SetIsAppEngine(val bool) {
	isAppEngine = val
}

func GetCacheControlVersionStamp() string {
	if IsLocal() {
		return ""
	}
	return cacheControlVersionStamp
}

func SetCacheControlVersionStamp(val string) {
	cacheControlVersionStamp = val
}

func GetRecaptchaSiteKey() string {
	return recaptchaSiteKey
}

func SetRecaptchaSiteKey(val string) {
	recaptchaSiteKey = val
}

func GetRecaptchaSecret() string {
	return recaptchaSecret
}

func SetRecaptchaSecret(val string) {
	recaptchaSecret = val
}

func GetMailjetApiKey() string {
	return mailjetApiKey
}

func SetMailjetApiKey(val string) {
	mailjetApiKey = val
}

func GetMailjetSecret() string {
	return mailjetSecret
}

func SetMailjetSecret(val string) {
	mailjetSecret = val
}
