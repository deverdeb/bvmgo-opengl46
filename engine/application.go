package engine

import (
	"fmt"
	"log"
	"ogl46/engine/assetsmngr"
	"ogl46/engine/graphic"
	"ogl46/engine/input"
)

type Application struct {
	// Etat de l'application
	isLaunched bool
	isStopped  bool
	isPause    bool

	// Scène en cours
	stage Stage

	// Window
	window input.Window
	// Renderer 2D
	renderer2d graphic.Renderer2d

	// Managers pour les ressources
	musicManager   *assetsmngr.MusicManager
	soundManager   *assetsmngr.SoundManager
	textureManager *assetsmngr.TextureManager
	fontManager    *assetsmngr.FontManager
}

func (app *Application) Window() *input.Window {
	return &app.window
}

func (app *Application) Renderer2d() *graphic.Renderer2d {
	return &app.renderer2d
}

func (app *Application) Execute() error {
	if app.isLaunched {
		return nil
	}
	app.isLaunched = true
	app.isStopped = false

	// Initialiser la partie vidéo
	err := app.window.Launch()
	if err != nil {
		return fmt.Errorf("échec du lancement de la fenêtre.\n - %w", err)
	}
	defer app.window.Close()

	if app.stage == nil {
		app.stage = &EmptyStage{}
	}

	app.initialize()

	// Boucle applicative
	timer := BuildTimer()
	for !app.isStopped {
		timer = timer.NextTimer()
		app.processStage(&timer)
	}

	app.release()

	app.isLaunched = false
	return nil
}

func (app *Application) initialize() {

	// Initialiser le renderer 2D
	if err := app.renderer2d.Initialize(); err != nil {
		log.Fatalf("échec de l'initialisation du renderer 2D\n - %v", err)
	}

	app.musicManager = assetsmngr.NewMusicManager()
	app.soundManager = assetsmngr.NewSoundManager()
	app.textureManager = assetsmngr.NewTextureManager()
	app.fontManager = assetsmngr.NewFontManager()
	app.stage.Initialize(app)
}

func (app *Application) release() {
	app.stage.Release(app)
	app.fontManager.ReleaseAll()
	app.textureManager.ReleaseAll()
	app.soundManager.ReleaseAll()
	app.musicManager.ReleaseAll()
	app.renderer2d.Release()
	app.window.Close()
}

func (app *Application) processStage(timer *Timer) {
	// Exécution
	if !app.isPause {
		app.stage.Execute(app, timer)
	}

	// Affichage
	app.window.Clear()
	app.stage.Display(app, timer)
	app.window.Swap()

	// Traitement des évènements
	app.window.ProcessEvents(func(event input.Event) {
		app.stage.ProcessEvent(app, event, timer)
	})
}

func (app *Application) IsLaunched() bool {
	return app.isLaunched
}
func (app *Application) IsStopped() bool {
	return app.isStopped
}

func (app *Application) Stop() {
	app.isStopped = true
}

func (app *Application) Stage() Stage {
	return app.stage
}

func (app *Application) SetStage(stage Stage) {
	if stage == app.stage {
		return
	}
	if app.isLaunched {
		app.stage.Release(app)
	}
	app.stage = stage
	if app.stage == nil {
		app.stage = &EmptyStage{}
	}
	if app.isLaunched {
		app.stage.Initialize(app)
	}
}

func (app *Application) MusicManager() *assetsmngr.MusicManager {
	return app.musicManager
}

func (app *Application) SoundManager() *assetsmngr.SoundManager {
	return app.soundManager
}

func (app *Application) TextureManager() *assetsmngr.TextureManager {
	return app.textureManager
}

func (app *Application) FontManager() *assetsmngr.FontManager {
	return app.fontManager
}

func (app *Application) VSync() bool {
	return app.window.VSync()
}

func (app *Application) SetVSync(vsync bool) {
	app.window.SetVSync(vsync)
}

func (app *Application) FullScreen() bool {
	return app.window.FullScreen()
}

func (app *Application) SetFullScreen(fulscreen bool) {
	app.window.SetFullScreen(fulscreen)
}

func (app *Application) Size() (width int, height int) {
	return app.window.Size()
}

func (app *Application) Keyboard() *input.Keyboard {
	return app.window.Keyboard()
}
func (app *Application) Mouse() *input.Mouse {
	return app.window.Mouse()
}
