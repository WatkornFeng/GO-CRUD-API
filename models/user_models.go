package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}

// If your model’s ID is uint (from gorm.Model),
// but you pass a uint64 without conversion,
// Go’s strict typing may cause compile errors or subtle runtime issues.
func GetUser(db *gorm.DB, id uint) (*User, error) {
	var user User
	result := db.First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
func GetAllUsers(db *gorm.DB) ([]User, error) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func UpdateUserName(db *gorm.DB, id uint, name string) (*User, error) {
	result := db.Model(&User{}).Where("id = ?", id).Update("name", name)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	// fetch the updated user to return
	var user User
	if err := db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// for testing result.Error
// type FakeUser struct {
// 	ID uint
// }
//result := db.Delete(&FakeUser{}, id)

func DeleteUser(db *gorm.DB, id uint) error {
	result := db.Delete(&User{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
