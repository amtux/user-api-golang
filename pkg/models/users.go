package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	tableName      struct{}  `pg:"users"`
	Id             int       `pg:"id,pk,unique,notnull"`
	Email          string    `json:"email,omitempty" pg:"email,unique,notnull"`
	FirstName      string    `json:"firstName,omitempty" pg:"first_name"`
	LastName       string    `json:"lastName,omitempty" pg:"last_name"`
	HashedPassword []byte    `pg:"enc_password,notnull"`
	CreatedAt      time.Time `pg:"created_at,notnull"`
	Active         bool      `pg:"active,notnull"`
	Password       string    `pg:"-"`
}

type UserDB struct {
	DB *pg.DB
}

func (udb *UserDB) Create(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return err
	}
	user.HashedPassword = hashedPassword
	user.CreatedAt = time.Now()
	user.Active = true

	_, err = udb.DB.Model(user).Returning("id").Insert()
	fmt.Println(err)
	return err
}

func (udb *UserDB) GetBasedOnId(id int) (User, error) {
	query := fmt.Sprintf("id = %d", id)
	var u User
	err := udb.DB.Model(&u).Where(query).Select()

	return u, err
}
func (udb *UserDB) GetBasedOnEmail(email string) (User, error) {
	var u User
	err := udb.DB.Model(&u).Where("email = ?", email).Select()
	fmt.Printf("in error %v\n", err)
	return u, err
}

func (udb *UserDB) Auth(user *User) (int, error) {
	u, err := udb.GetBasedOnEmail(user.Email)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(user.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, errors.New("Invalid Credentials")
		} else {
			return 0, err
		}
	}
	return u.Id, nil
}

func (udb *UserDB) Update(user User) error {
	_, err := udb.DB.Model(&user).WherePK().Update()
	return err
}
