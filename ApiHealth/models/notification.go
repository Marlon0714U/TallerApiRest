package models

type NotificationMessage struct {
	Tipo         string `json:"tipo"`         // Canal de la notificación
	Mensaje      string `json:"mensaje"`      // Contenido del mensaje
	Destinatario string `json:"destinatario"` // Destinatario del mensaje
}
