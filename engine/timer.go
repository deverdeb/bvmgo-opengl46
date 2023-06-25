package engine

import "time"

type Timer struct {
	// currentTime est le temps courant (en millisecondes)
	currentTime int64
	// previousTime est le temps précédent (en millisecondes)
	previousTime int64
	// elapsedTime est le temps écoulé entre le temps courant et le temps précédent (en millisecondes)
	elapsedTime int64
}

// BuildTimer construit un timer à partir de l'heure courante
func BuildTimer() Timer {
	currentTime := time.Now().UnixMilli()
	return Timer{
		currentTime:  currentTime,
		previousTime: currentTime,
		elapsedTime:  currentTime - currentTime,
	}
}

// NextTimer construit un timer à partir d'un ancien timer
func (timer *Timer) NextTimer() Timer {
	currentTime := time.Now().UnixMilli()
	return Timer{
		currentTime:  currentTime,
		previousTime: timer.currentTime,
		elapsedTime:  currentTime - timer.currentTime,
	}
}

// CurrentTime retourne le temps courant en millisecondes
func (timer *Timer) CurrentTime() int64 {
	return timer.currentTime
}

// PreviousTime retourne le temps précédent en millisecondes
func (timer *Timer) PreviousTime() int64 {
	return timer.previousTime
}

// ElapsedTime retourne le temps écoulé en millisecondes
func (timer *Timer) ElapsedTime() int64 {
	return timer.elapsedTime
}
