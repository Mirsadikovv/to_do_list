package check

import (
	"errors"
	"regexp"
	"time"
)

func ValidatePhone(phone string) error {
	phoneRegex := regexp.MustCompile(`^\+998[0-9]{9}$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("phone is not valid")
	}
	return nil
}

func ValidateMail(mail string) error {
	mailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@(gmail\.com|mail\.ru)$`)
	if mailRegex.MatchString(mail) {
		return errors.New("mail is not valid")
	}
	return nil
}

func CheckDeadline(timestamp string) (float64, error) {
	layout := time.RFC3339

	date, err := time.Parse(layout, timestamp)
	if err != nil {
		return -1, errors.New("wrong timestamp format")
	}

	now := time.Now()
	timeUntilDeadline := date.Sub(now)

	if timeUntilDeadline < 0 {
		if timeUntilDeadline >= -72*time.Hour {
			return 1, nil
		}
		return 0, nil
	}

	return timeUntilDeadline.Hours(), nil
}
