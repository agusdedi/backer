package user

import (
	"backer/config"
	"strings"
)

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	ImageURL   string `json:"image_url"`
	Token      string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		ImageURL:   buildImageURL(user.AvatarFileName),
		Token:      token,
	}

	return formatter
}

func buildImageURL(fileName string) string {
	if fileName == "" {
		return ""
	}

	if strings.HasPrefix(fileName, "http") {
		return fileName
	}

	return config.AppConfig.ImageBaseURL + "/" + fileName
}
