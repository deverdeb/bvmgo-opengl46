package input

// EventSource représente la source d'un évènement du moteur.
type EventSource uint8

const (
	// WINDOW indique un évènement provenant de la fenêtre
	WINDOW EventSource = iota
	// KEYBOARD indique un évènement provenant du clavier
	KEYBOARD
	// MOUSE_BUTTON indique un évènement provenant d'un bouton de la souris
	MOUSE_BUTTON
	// MOUSE_MOVE indique un évènement provenant d'un déplacement de la souris
	MOUSE_MOVE
	// MOUSE_SCROLL indique un évènement provenant de la molette de la souris
	MOUSE_SCROLL
)

// Event représente un évènement remonté par le moteur.
type Event interface {
	// Source retourne la source de l'évènement
	Source() EventSource
}
