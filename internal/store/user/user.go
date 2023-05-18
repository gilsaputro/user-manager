package user

import (
	"errors"
	"math"

	"github.com/jinzhu/gorm"

	"gilsaputro/user-manager/pkg/postgres"
)

// UserStoreMethod is set of methods for interacting with a user storage system
type UserStoreMethod interface {
	CreateUser(userinfo UserStoreInfo) error
	UpdateUser(userinfo UserStoreInfo) error
	DeleteUser(userid int) error
	GetUserInfoByUsername(username string) (UserStoreInfo, error)
	GetUserInfoByID(userid int) (UserStoreInfo, error)
	GetAllUserInfoWithPagging(size, cursor int) ([]UserStoreInfo, int, error)
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
		Email:    userinfo.Email,
	}

	return db.Create(user).Error
}

// UpdateUser is func to edit / update user info into database
func (u *UserStore) UpdateUser(userinfo UserStoreInfo) error {
	db, err := u.getDB()
	if err != nil {
		return err
	}

	var user postgres.User

	err = db.Where("username = ? AND id = ?", userinfo.Username, userinfo.UserId).First(&user).Error
	if err != nil {
		return err
	}

	user.Password = userinfo.Password
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

	return db.Unscoped().Delete(&user).Error
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
		Username:    user.Username,
		UserId:      int(user.ID),
		Fullname:    user.Fullname,
		Email:       user.Email,
		CreatedDate: user.CreatedAt.String(),
	}, nil
}

// GetAllUserInfoWithPagging is func to get all data user info in database
func (u *UserStore) GetAllUserInfoWithPagging(size, cursor int) ([]UserStoreInfo, int, error) {
	db, err := u.getDB()
	var totalCount int
	if err != nil {
		return []UserStoreInfo{}, 0, err
	}

	var users []postgres.User
	query := db.Limit(size).Offset((cursor - 1) * size).Order("id asc")

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Model(&postgres.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	totalCursor := int(math.Ceil(float64(totalCount) / float64(size)))
	var nextCursor int
	if cursor < totalCursor {
		nextCursor = cursor + 1
	} else {
		nextCursor = 0 // 0 menunjukkan tidak ada halaman berikutnya
	}

	var listUser []UserStoreInfo
	for _, user := range users {
		listUser = append(listUser, UserStoreInfo{
			UserId:      int(user.ID),
			Username:    user.Username,
			Fullname:    user.Fullname,
			Email:       user.Email,
			CreatedDate: user.CreatedAt.String(),
		})
	}

	return listUser, nextCursor, err
}
