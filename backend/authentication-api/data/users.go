package data

import (
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-hclog"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// User define la estructura de un usuario para el API
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

// UserSignin defines user when is on Signin phase
type UserSignin struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
	Rol      int    `json:"rol"`
}

// UserCreate defines data user structure when realices a signup
type UserCreate struct {
	ID        int    `json:"id"`
	Firstname string `json:"Firstname" validate:"required"`
	Lastname  string `json:"Lastname" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Rol       int    `json:"rol"`
}

//Users is una colección de User
type Users []*User

//UserService representa una implementación de mysql
type UserService struct {
	DB *sql.DB
	l  hclog.Logger
}

// New creates a new user service
func New(d *sql.DB, l hclog.Logger) *UserService {
	return &UserService{d, l}
}

// GetUsers retorna una lista de usuarios
func (s *UserService) GetUsers() (Users, error) {
	users := Users{}
	rows, err := s.DB.Query("SELECT id, name FROM users")
	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := &User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

//GetUserByID returns an user given an id
func (s *UserService) GetUserByID(id int) (User, error) {
	user := User{}
	rows, err := s.DB.Query("SELECT id, name FROM users WHERE id = (?)", id)
	if err != nil {
		return user, ErrProductNotFound
	}

	for rows.Next() {
		user = User{}
		err = rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return user, err
		}

		return user, err
	}

	return user, ErrProductNotFound
}

//GetUserByPhone returns an user given an phone number
func (s *UserService) GetUserByPhone(phone string) (UserSignin, error) {
	s.l.Info("Getting user from database with", "phone", phone)

	user := UserSignin{}
	rows, err := s.DB.Query("SELECT phone, password, rol FROM users WHERE phone = (?)", phone)
	if err != nil {
		return user, ErrProductNotFound
	}

	for rows.Next() {
		user = UserSignin{}
		err = rows.Scan(&user.Phone, &user.Password, &user.Rol)
		if err != nil {
			return user, err
		}

		return user, err
	}

	return user, ErrProductNotFound
}

//CreateUser crea un usuario
func (s *UserService) CreateUser(pUser *UserCreate) error {
	s.l.Info("Creating", "user", pUser)
	saltedPassword, err := s.hashAndSalt([]byte(pUser.Password))
	if err != nil {
		return err
	}
	_, err = s.DB.Exec("INSERT INTO users (firstname, lastname, password, phone, rol) VALUES (?, ?, ?, ?, ?)",
		pUser.Firstname,
		pUser.Lastname,
		saltedPassword,
		pUser.Phone,
		pUser.Rol)

	return err
}

//DeleteUser elimina un usuario dado un id
func (s *UserService) DeleteUser(id int) error {
	_, err := s.DB.Exec("DELETE FROM users WHERE id = (?)", id)
	if err != nil {
		return err
	}
	return nil
}
