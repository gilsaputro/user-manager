package user

import (
	"errors"

	"github.com/jinzhu/gorm"

	"gilsaputro/user-manager/pkg/postgres"
)

// UserStoreMethod is set of methods for interacting with a user storage system
type UserStoreMethod interface {
	CreateUser(userinfo UserStoreInfo) error
	UpdateUser(userid int, userinfo UserStoreInfo) error
	DeleteUser(userid int) error
	GetUserInfoByUsername(username string) (UserStoreInfo, error)
	GetUserInfoByID(userid int) (UserStoreInfo, error)
	GetAllUserInfoWithPagging(userid, size, cursor int) ([]UserStoreInfo, error)
}

// UserStore is list dependencies user store
type UserStore struct {
	pg postgres.PostgresMethod
}

// NewUserStore is func to generate UserStoreMethod interface
func NewUserStore(pg postgres.PostgresMethod) UserStoreMethod {
	return &UserStore{
		pg: pg,
	}
}

func (u *UserStore) getDB() (*gorm.DB, error) {
	db := u.pg.GetDB()
	if db == nil {
		return nil, errors.New("Database Client is not init")
	}

	return db, nil
}

// CreateUser is func to store / create user info into database
func (u *UserStore) CreateUser(userinfo UserStoreInfo) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}
	user := &postgres.User{
		Username: userinfo.Username,
		Password: userinfo.Password,
		Fullname: userinfo.Fullname,
		Email:    userinfo.Fullname,
	}

	return db.Create(user).Error
}

// UpdateUser is func to edit / update user info into database
func (u *UserStore) UpdateUser(userid int, userinfo UserStoreInfo) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	var user postgres.User

	err = db.Where("username = ? AND id = ?", userinfo.Username, userid).First(&user).Error
	if err != nil {
		return err
	}

	user.Fullname = userinfo.Fullname
	user.Email = userinfo.Email

	return db.Save(&user).Error
}

// GetUserID is func to get user id by username and password
func (u *UserStore) GetUserInfoByUsername(username string) (UserStoreInfo, error) {
	var user postgres.User
	db, err := u.getDB()
	if err != nil {
		return UserStoreInfo{}, err
	}

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return UserStoreInfo{}, err
	}

	return UserStoreInfo{
		UserId:   int(user.ID),
		Username: user.Username,
		Password: user.Password,
	}, nil
}

// DeleteUser is func to delete user info on database
func (u *UserStore) DeleteUser(userid int) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	user := postgres.User{
		Model: gorm.Model{
			ID: uint(userid),
		},
	}

	return db.Delete(&user).Error
}

// GetUserByID is func to get user info by id on database
func (u *UserStore) GetUserInfoByID(userid int) (UserStoreInfo, error) {
	var user postgres.User
	db, err := u.getDB()
	if err != nil {
		return UserStoreInfo{}, err
	}

	if err := db.First(&user, userid).Error; err != nil {
		return UserStoreInfo{}, err
	}

	return UserStoreInfo{
		Username: user.Username,
		Password: user.Password,
	}, nil
}

// GetAllUserInfoWithPagging is func to get all data user info in database
func (u *UserStore) GetAllUserInfoWithPagging(userid, size, cursor int) ([]UserStoreInfo, error) {
	db, err := u.getDB()
	if err != nil {
		return []UserStoreInfo{}, err
	}

	var users []postgres.User
	err = db.Find(&users).Limit(size).Offset((cursor - 1) * size).Error
	if err != nil {
		return nil, err
	}

	var listUser = make([]UserStoreInfo, len(users))
	for _, user := range users {
		listUser = append(listUser, UserStoreInfo{
			UserId:   int(user.ID),
			Username: user.Username,
		})
	}

	return listUser, err
}
