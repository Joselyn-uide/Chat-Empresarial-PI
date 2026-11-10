package models

import (
	"errors"
	"strings"
)

type Usuario struct {
	nombreUsuario string
	rol           string
}

// NewUsuario es un constructor para crear un usuario con estado válido
func NewUsuario(nombreUsuario, rol string) (*Usuario, error) {
	nombreUsuario = strings.TrimSpace(nombreUsuario)
	if nombreUsuario == "" {
		return nil, errors.New("El nombre de usuario no puede estar vacío")
	}
	if rol == "" {
		rol = "Empleado" // Rol por defecto
	}
	return &Usuario{
		nombreUsuario: nombreUsuario,
		rol: rol,
	}, nil
}

// GetNombreUsuario es un método asociado a la estructura con receptor para obtener el nombre 
func (u *Usuario) GetNombreUsuario() string {
	return u.nombreUsuario
}

// GetRol permite consultar el rol del usuario 
func (u *Usuario) GetRol() string {
	return u.rol
}
