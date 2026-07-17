package auth

import "time"

const sessionTokenCookieName = "session_token"

const sessionTokenCookieExpiry = time.Hour * 24 * 30
const sessionTokenCookieRenewIfExpiresIn = time.Hour * 24 * 29

const PasswordMinLength = 8
const passwordBcryptCost = 15

const signupRequestTokenLength = 15
const signupTokenExpiry = time.Hour

const userHandleMaxLength = 25
const userDisplayNameMaxLength = 50

// require alphabetic first character, then alphanumeric or underscore
const UserHandlePattern = `^[a-zA-Z][a-zA-Z0-9_]*$`
