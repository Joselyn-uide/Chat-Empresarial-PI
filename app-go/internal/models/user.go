package models

import (
	"errors"
)

type Usuario struct {
	NombreUsuario string
	Rol           string
}

// NewUsuario es un constructor para crear un usuario con estado válido
func NewUsuario(nombreUsuario, rol string) (*Usuario, error) {
	if nombreUsuario == "" {
		return nil, errors.New("El nombre de usuario no puede estar vacío")
	}
	if rol == "" {
		rol = "Empleado" // Rol por defecto
	}
	return &Usuario{
		NombreUsuario: nombreUsuario,
		Rol: rol,
	}, nil
}

// GetNombreUsuario es un método asociado a la estructura con receptor para obtener el nombre 
func (u *Usuario) GetNombreUsuario() string {
	return u.NombreUsuario
}

// GetRol permite consultar el rol del usuario 
func (u *Usuario) GetRol() string {
	return u.Rol
}
