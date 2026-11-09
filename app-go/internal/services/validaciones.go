// Package services
// En este paquete reunimos las validaciones de entrada del cliente.
// Revisamos los datos que escribe el usuario (IP, puerto, nombre y mensajes)
// antes de que salgan por el socket, y si algo está mal devolvemos un error
// claro en lugar de dejar que el programa falle más adelante.
package services

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Definimos los límites en constantes para poder cambiarlos fácilmente.
const (
	MaxLongitudMensaje = 500
	MinLongitudUsuario = 3
	MaxLongitudUsuario = 20
)

// Errores propios de validación, para poder compararlos con errors.Is.
var (
	ErrMensajeVacio        = errors.New("el mensaje no puede estar vacío ni tener solo espacios")
	ErrMensajeMuyLargo     = fmt.Errorf("el mensaje no puede superar los %d caracteres", MaxLongitudMensaje)
	ErrCaracteresInvalidos = errors.New("el texto contiene caracteres no permitidos")
	ErrUsuarioInvalido     = fmt.Errorf("el usuario debe tener entre %d y %d caracteres y solo letras, números, '_', '-' o '.'", MinLongitudUsuario, MaxLongitudUsuario)
	ErrIPInvalida          = errors.New("la dirección IP no tiene un formato válido (ejemplo: 127.0.0.1)")
	ErrPuertoInvalido      = errors.New("el puerto debe ser un número entre 1 y 65535")
)

// ValidarMensaje limpia los espacios del inicio y del final y revisa que el
// mensaje no esté vacío, no sea demasiado largo y no tenga caracteres de control.
// Si todo está bien, devuelve el mensaje ya limpio.
func ValidarMensaje(texto string) (string, error) {
	texto = strings.TrimSpace(texto)
	if texto == "" {
		return "", ErrMensajeVacio
	}
	// Contamos runas y no bytes, para que las tildes y la ñ cuenten como un solo carácter
	if utf8.RuneCountInString(texto) > MaxLongitudMensaje {
		return "", ErrMensajeMuyLargo
	}
	if !utf8.ValidString(texto) {
		return "", ErrCaracteresInvalidos
	}
	for _, r := range texto {
		// Los caracteres de control (como \n o \x00) podrían romper el
		// protocolo, porque el servidor separa los mensajes por líneas
		if unicode.IsControl(r) {
			return "", ErrCaracteresInvalidos
		}
	}
	return texto, nil
}

// ValidarUsuario revisa que el nombre de usuario tenga una longitud adecuada
// y solo use letras, números, guion bajo, guion o punto.
func ValidarUsuario(nombre string) (string, error) {
	nombre = strings.TrimSpace(nombre)
	largo := utf8.RuneCountInString(nombre)
	if largo < MinLongitudUsuario || largo > MaxLongitudUsuario {
		return "", ErrUsuarioInvalido
	}
	for _, r := range nombre {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' && r != '.' {
			return "", ErrUsuarioInvalido
		}
	}
	return nombre, nil
}

// ValidarIP acepta "localhost" o una dirección IP con formato correcto.
func ValidarIP(ip string) (string, error) {
	ip = strings.TrimSpace(ip)
	if strings.EqualFold(ip, "localhost") {
		return "localhost", nil
	}
	if net.ParseIP(ip) == nil {
		return "", ErrIPInvalida
	}
	return ip, nil
}

// ValidarPuerto revisa que el puerto sea un número dentro del rango permitido.
func ValidarPuerto(puerto string) (string, error) {
	puerto = strings.TrimSpace(puerto)
	numero, err := strconv.Atoi(puerto)
	if err != nil || numero < 1 || numero > 65535 {
		return "", ErrPuertoInvalido
	}
	return strconv.Itoa(numero), nil
}