package users

import "regexp"

const (
	maxNameLength = 20
	minNameLength = 3

	minPassLength        = 6
	maxPassLength        = 120
	invalidPassLengthMsg = "length should be between 6 and 120"

	inUseMsg             = "is already in use"
	invalidNameLengthMsg = "length should be between 3 and 20"

	minEmailLength        = 5
	maxEmailLength        = 254
	invalidEmailLengthMsg = "length should be between 5 and 254"
	invalidEmailFormatMsg = "format is invalid"
)

var emailRegex = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}" +
		"[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

// Validate is a method that validates a user's fields
func (u *User) Validate() map[string][]string {
	validationErrors := map[string][]string{}

	validateEmailFormat(validationErrors, u.Email)
	validateEmailLength(validationErrors, u.Email)
	validateEmailUniqueness(validationErrors, u.Email)

	validateNameLength(validationErrors, u.FirstName, "first_name")
	validateNameLength(validationErrors, u.LastName, "last_name")

	validatePasswordLength(validationErrors, u.Password)

	return validationErrors
}

func validateEmailUniqueness(errorMsgs map[string][]string, e string) {
	exists := EmailExists(e)

	if exists {
		errorMsgs["email"] = append(errorMsgs["email"], inUseMsg)
	}
}

func validateEmailLength(errorMsgs map[string][]string, email string) {
	length := len(email)

	if length < minEmailLength || length > maxEmailLength {
		errorMsgs["email"] = append(errorMsgs["email"], invalidEmailLengthMsg)
	}
}

func validateEmailFormat(errorMsgs map[string][]string, e string) {
	if !emailRegex.MatchString(e) {
		errorMsgs["email"] = append(errorMsgs["email"], invalidEmailFormatMsg)
	}
}

func validateNameLength(errorMsgs map[string][]string, name, key string) {
	if validateLength(name, minNameLength, maxNameLength) {
		errorMsgs[key] = append(errorMsgs[key], invalidNameLengthMsg)
	}
}

func validatePasswordLength(errorMsgs map[string][]string, pass string) {
	if validateLength(pass, minPassLength, maxPassLength) {
		errorMsgs["password"] = append(errorMsgs["password"], invalidPassLengthMsg)
	}
}

func validateLength(name string, min, max int) bool {
	length := len(name)

	if length < min || length > max {
		return true
	}

	return false
}
