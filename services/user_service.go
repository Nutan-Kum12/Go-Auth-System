package services

import (
	"Auth/model"
	"Auth/repository"
	"Auth/utils"
	"errors"

	// "fmt"
	// "time"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	Repo repository.UserRepository // service need repository for DB operation
}

func (s UserService) Register(user model.User) error {
	user.Email = utils.NormalizeEmail(user.Email)
	user.Name = utils.NormalizeName(user.Name)
	// Check for Existing User
	_, err := s.Repo.FindUserByEmail(user.Email)
	if err == nil {
		return errors.New(
			"user already exists",
		)
	}
	if err != mongo.ErrNoDocuments {
		return err
	}

	if !utils.ValidateEmail(user.Email) {
		return errors.New("Invalid Email")
	}
	if !utils.ValidatePassword(user.Password) {
		return errors.New("Password must be of length atleast 8 with one uppercase,one lowercase,one digit,one special character")
	}
	if !utils.ValidateName(user.Name) {
		return errors.New("Name must be of length atleast 2 char and atmost 40 char only letters and space")
	}
	// start := time.Now()
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	// fmt.Println(time.Since(start)) //time taken for hashing

	user.Password = hashedPassword
	err = s.Repo.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (s UserService) Login(email string, password string) (string, string, error) {
	email = utils.NormalizeEmail(email)
	//get user from DB
	user, err := s.Repo.FindUserByEmail(email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", errors.New(
				"Invalid credentials",
			)
		}
		return "", "", err
	}
	//compare password
	err = utils.ComparePassword(password, user.Password)
	if err != nil {
		return "", "", errors.New("Invalid credentials")
	}
	//generate access token
	accessToken, err := utils.GenerateAccessToken(
		user.Email,
	)
	if err != nil {
		return "", "", err
	}
	//generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(
		user.Email,
	)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s UserService) GetProfile(email string) (model.User, error) {
	return s.Repo.FindUserByEmail(email)
}
