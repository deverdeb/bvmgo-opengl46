package main

/*
Adapted from this tutorial: http://www.learnopengl.com/#!Getting-started/Camera

Shows how to create a basic controllable FPS camera. This has been refactored into
classes to allow better reuse going forward.
*/

import (
	"log"
	"ogl46/engine"
)

func main() {
	app := engine.Application{}
	app.SetVSync(true)
	app.SetStage(ExampleStage)
	err := app.Execute()
	if err != nil {
		log.Printf("error found:\n - %v", err)
	}
}
