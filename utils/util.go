package utils

import "unicode"

func ValidatePassword(input *string) (validPassword string, errorList []string) {
	password := ""

	if input != nil {
		password = *input
	}

	if len(password) < 6 || len(password) > 64 {
		errorList = append(errorList, "password should be 6 to 64 characters")
	}

	containsCapital, containsNumber, containsSpecialAlphaNumeric := false, false, false
	for _, c := range password {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			containsSpecialAlphaNumeric = true
		}
		if unicode.IsUpper(c) {
			containsCapital = true
		}
		if unicode.IsNumber(c) {
			containsNumber = true
		}

		if containsSpecialAlphaNumeric && containsCapital && containsNumber {
			break
		}
	}

	if !containsCapital {
		errorList = append(errorList, "password should contain a capital letter")
	}
	if !containsNumber {
		errorList = append(errorList, "password should contain a number")
	}
	if !containsSpecialAlphaNumeric {
		errorList = append(errorList, "password should contain a special alphanumeric character")

	}

	if len(errorList) == 0 {
		validPassword = password
	}

	return validPassword, errorList
}
