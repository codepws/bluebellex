package login

import (
	"bluebell_backend/models/login"
)

func Login(signInRequest *login.SignInRequest) (user *login.User, err error) {

	user = new(login.User)
	err = nil
	//
	user.UserName = signInRequest.UserName
	user.UserID = uint64(111111)
	user.VIP = 55
	user.Gender = 1
	user.Email = "123@a.com"
	user.Addresses = nil

	return
}

func Register(signUpRequest *login.SignUpRequest) (user *login.User, err error) {

	user = new(login.User)
	err = nil
	//
	user.UserName = signUpRequest.UserName
	user.UserID = uint64(111111)
	user.VIP = 55
	user.Gender = signUpRequest.Gender
	user.Email = signUpRequest.Email
	user.Addresses = signUpRequest.Addresses

	return
}
