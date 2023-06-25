package engine

const measuresNumber = 10

// Mesure vide
var emptyMeasure = measure{
	elapsedTime: 0.,
	frames:      0,
}

type measure struct {
	// elapsedTime est le temps écoulé pour cette mesure en millisecondes
	elapsedTime int64
	// frames est le nombre d'images compté durant le temps écoulé
	frames int64
}

type FpsCounter struct {
	measures       [measuresNumber]measure
	currentMeasure uint8
}

func BuildFpsCounter() FpsCounter {
	counter := FpsCounter{
		currentMeasure: 0,
	}
	for idx := 0; idx < measuresNumber; idx++ {
		counter.measures[idx] = emptyMeasure
	}
	return counter
}

// Increment ajoute une frame et son temps de rendu (en millisecondes) au compteur de FPS
func (counter *FpsCounter) Increment(elapsedTime int64) {
	counter.measures[counter.currentMeasure].frames += 1
	counter.measures[counter.currentMeasure].elapsedTime += elapsedTime
	if counter.measures[counter.currentMeasure].elapsedTime >= 1000 {
		counter.currentMeasure = (counter.currentMeasure + 1) % measuresNumber
		counter.measures[counter.currentMeasure] = emptyMeasure
	}
}

// ComputeFps retourne le nombre de FTP calculé depuis les dernières mesures
func (counter *FpsCounter) ComputeFps() float32 {
	elapsedTime := int64(0)
	frames := int64(0)
	for idx := 0; idx < measuresNumber; idx++ {
		elapsedTime += counter.measures[idx].elapsedTime
		frames += counter.measures[idx].frames
	}
	if elapsedTime == 0 {
		return 0
	} else {
		return float32(frames*1000) / float32(elapsedTime)
	}
}
