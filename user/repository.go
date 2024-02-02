package user

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindById(id int) (User, error)
	Update(user User) (User, error)
	FindAllUsers() ([]User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		fmt.Println("error")
		log.Fatal(err.Error())
		return user, err
	}
	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {

	var user User

	err := r.db.Where("Email = ? ", email).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindById(id int) (User, error) {
	var user User

	err := r.db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) Update(user User) (User, error) {

	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil

}

func (r *repository) FindAllUsers() ([]User, error) {
	var users []User
	err := r.db.Find(&users).Error
	if err != nil {
		return users, err
	}
	return users, nil
}
