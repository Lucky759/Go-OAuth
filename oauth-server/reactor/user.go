package reactor

import (
	"oauth-server/structs"
)

// UserInfo UserInfo
func UserInfo(id uint) (*structs.User, error) {
	user := structs.User{}
	dbr := mainDB.Where("id=?", id).First(&user)
	return &user, dbr.Error
}

// UserFindWithUsername UserFindWithUsername
func UserFindWithUsername(username string) (*structs.User, error) {
	user := structs.User{}
	dbr := mainDB.Where("username=?", username).First(&user)
	if dbr.Error == nil {
		return &user, nil
	}
	return nil, dbr.Error
}

// UserFindWithMobile UserFindWithMobile
func UserFindWithMobile(mobile string) (*structs.User, error) {
	user := structs.User{}
	dbr := mainDB.Where("mobile=?", mobile).First(&user)
	if dbr.Error == nil {
		return &user, nil
	}
	return nil, dbr.Error
}

// UserFindWithEmail UserFindWithEmail
func UserFindWithEmail(email string) (*structs.User, error) {
	user := structs.User{}
	dbr := mainDB.Where("email=?", email).First(&user)
	if dbr.Error == nil {
		return &user, nil
	}
	return nil, dbr.Error
}

// UserCreateOption UserCreateOption
type UserCreateOption struct {
	Username, Mobile, Email, Password, Avatar, Intro, GiteaSSOToken, GiteaUsername *string
}

// UserUpdateOption UserUpdateOption
type UserUpdateOption struct {
	Username, Mobile, Email, Password, Avatar, Intro, GiteaSSOToken, Experience *string
	StatusID, ParttimeStatusID                                                  *uint
}