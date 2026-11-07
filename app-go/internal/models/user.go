package models

import (
	"errors"
)

// User representa la entidad principal del colaborador (POO)
type User struct {
	Username string
	Role     string
}

// NewUser es un constructor para crear un usuario con estado válido
func NewUser(username, role string) (*User, error) {
	if username == "" {
		return nil, errors.New("El nombre de usuario no puede estar vacío")
	}
	if role == "" {
		role = "Empleado" // Rol por defecto
	}
	return &User{
		Username: username,
		Role:     role,
	}, nil
}

// GetUsername es un método asociado a la estructura con receptor para obtener el nombre 
func (u *User) GetUsername() string {
	return u.Username
}

// GetRole permite consultar el rol del usuario 
func (u *User) GetRole() string {
	return u.Role
}
