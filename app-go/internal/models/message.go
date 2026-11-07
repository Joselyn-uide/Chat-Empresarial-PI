package models

import (
	"errors"
	"strings"
	"time"
)

// Mensaje representa la estructura de un mensaje enviado en el chat
type Mensaje struct {
	Remitente string
	Contenido string
	FechaHora time.Time
}

// NewMensaje es la función constructora para crear y validar un mensaje de texto
func NewMensaje(remitente, contenido string) (*Mensaje, error) {
	remitente = strings.TrimSpace(remitente)
	contenido = strings.TrimSpace(contenido)

	if remitente == "" {
		return nil, errors.New("El remitente del mensaje no puede estar vacío")
	}
	if contenido == "" {
		return nil, errors.New("El contenido del mensaje no puede estar vacío")
	}

	return &Mensaje{
		Remitente: remitente,
		Contenido: contenido,
		FechaHora: time.Now(),
	}, nil
}

// GetRemitente retorna el remitente del mensaje
func (m *Mensaje) GetRemitente() string {
	return m.Remitente
}

// GetContenido retorna el contenido del mensaje
func (m *Mensaje) GetContenido() string {
	return m.Contenido
}

// GetFechaHora retorna la fecha y hora en que se envió el mensaje
func (m *Mensaje) GetFechaHora() time.Time {
	return m.FechaHora
}
