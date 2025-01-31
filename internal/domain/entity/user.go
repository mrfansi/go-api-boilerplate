package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Name      string    `json:"name" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;default:'user'"`
	Active    bool      `json:"active" gorm:"not null;default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUser(email, password, name string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New(),
		Email:     email,
		Password:  string(hashedPassword),
		Name:      name,
		Role:      "user",
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Update(name string) {
	u.Name = name
	u.UpdatedAt = time.Now()
}

func (u *User) SetRole(role string) {
	u.Role = role
	u.UpdatedAt = time.Now()
}

func (u *User) SetActive(active bool) {
	u.Active = active
	u.UpdatedAt = time.Now()
}