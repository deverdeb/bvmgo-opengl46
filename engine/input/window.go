package input

import (
	"fmt"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log/slog"
	"runtime"
)

type Window struct {
	// glfwWindow est la fenêtre GLFW.
	glfwWindow *glfw.Window

	// Gestionnaire d'évènements
	eventManager EventManager

	// mouse représente l'état de la souris
	mouse Mouse

	// keyboard représente l'état du clavier
	keyboard Keyboard

	// vsync indique si la synchronisation verticale est activée
	vsync bool

	//
	fullscreen    bool
	initialX      int
	initialY      int
	initialWidth  int
	initialHeight int
}

// Launch initialise le moteur d'affichage et affiche la fenêtre
func (window *Window) Launch() error {
	slog.Debug("launch window")
	if window.glfwWindow != nil {
		return fmt.Errorf("window is already launched")
	}

	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize glfw\n - %w", err)
	}

	// Create glfw window
	err = window.createGlfwWindow()
	if err != nil {
		window.Close()
		return fmt.Errorf("failed to create window\n - %w", err)
	}

	// Mettre en place le gestionnaire d'évènements
	window.mouse.glfwWindow = window.glfwWindow
	window.eventManager.Initialize(window)

	// Initialize Glow
	window.glfwWindow.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		window.Close()
		return fmt.Errorf("failed to initialize OpenGL (Glow)\n - %w", err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	vendor := gl.GoStr(gl.GetString(gl.VENDOR))
	slog.Info("OpenGL", "version", version, "vendor", vendor)

	// Positionner les propriétés VSync et Fullscreen
	window.SetVSync(window.VSync())

	slog.Debug("window is launched")
	return nil
}

// Close demande à la fenêtre de se fermer.
func (window *Window) Close() {
	if window.glfwWindow != nil {
		slog.Debug("release window")
		// Libérer le gestionnaire d'évènements
		window.eventManager.Release()
		// Fermer la fenêtre glfw
		window.releaseGlfwWindow()
	}
}

func (window *Window) Clear() {
	// Z-Buffer
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	// Vider les buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Disable(gl.BLEND)

}

func (window *Window) Swap() {
	if window.glfwWindow != nil {
		window.glfwWindow.SwapBuffers()
	}
}

// IsLaunched indique si la fenêtre est lancée
func (window *Window) IsLaunched() bool {
	return window.glfwWindow != nil
}

// createGlfwWindow met en place la fenêtre GLFW
func (window *Window) createGlfwWindow() error {
	windowWidth := 800
	windowHeight := 600

	// Configure OpengGL for GLFW
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	var mode = glfw.GetPrimaryMonitor().GetVideoMode()

	glfw.WindowHint(glfw.RedBits, mode.RedBits)
	glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
	glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
	glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)

	// Création de la fenêtre GLFW
	var err error
	window.glfwWindow, err = glfw.CreateWindow(windowWidth, windowHeight, "<my title>", nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create glfw window\n - %w", err)
	}

	return nil
}

// releaseGlfwWindow libère la fenêtre GLFW
func (window *Window) releaseGlfwWindow() {
	// Libérer la fenêtre glfw
	if window.glfwWindow != nil {
		glfwWindow := window.glfwWindow
		window.glfwWindow = nil
		glfwWindow.Destroy()
		// Quitter glfw
		glfw.Terminate()
	}
}

// ProcessEvents demande le traitement des évènements reçus par la fenêtre
func (window *Window) ProcessEvents(eventCallback EventProcessor) {
	window.eventManager.ProcessEvents(eventCallback)
}

// Size retourne la dimension de la fenêtre
func (window *Window) Size() (width int, height int) {
	return window.glfwWindow.GetSize()
}

func (window *Window) VSync() bool {
	return window.vsync
}

func (window *Window) SetVSync(vsync bool) {
	window.vsync = vsync
	if window.IsLaunched() {
		if vsync {
			glfw.SwapInterval(1)
		} else {
			glfw.SwapInterval(0)
		}
	}
}

func (window *Window) FullScreen() bool {
	return window.fullscreen
}

func (window *Window) SetFullScreen(fullscreen bool) {
	if fullscreen == window.fullscreen {
		window.fullscreen = fullscreen
		if window.IsLaunched() {
			window.applyFullScreen()
		}
	}
}

func (window *Window) applyFullScreen() {
	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()
	if window.fullscreen {
		// Conserver position et taille précédente
		window.initialX, window.initialY = window.glfwWindow.GetPos()
		window.initialWidth, window.initialHeight = window.glfwWindow.GetSize()
		// Switch en plein écran
		window.glfwWindow.SetMonitor(monitor, 0, 0, mode.Width, mode.Height, mode.RefreshRate)
	} else {
		// Restaurer position et taille précédente
		window.glfwWindow.SetMonitor(nil, window.initialX, window.initialY, window.initialWidth, window.initialHeight, mode.RefreshRate)
	}
	// Reforcer la synchronisation vertical
	window.SetVSync(window.VSync())
}

func (window *Window) Mouse() *Mouse {
	return &window.mouse
}

func (window *Window) Keyboard() *Keyboard {
	return &window.keyboard
}
