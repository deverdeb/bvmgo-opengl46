package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl64"
)

// EventProcessor permet de transmettre les évènements pour traitement.
type EventProcessor func(event Event)

// EventManager a pour but de gérer les évènements entrant (clavier, souris, action sur la fenêtre, ...).
type EventManager struct {
	// window représente la fenêtre à suivre
	window *Window

	// processor est la fonction vers laquelle rediriger les évènements.
	processor EventProcessor
}

// Initialize positionne les différentes fonctions de traitement des évènements.
func (mng *EventManager) Initialize(window *Window) {
	if window == nil {
		return
	}
	mng.window = window
	glfwWindow := mng.window.glfwWindow

	// Window Callback
	glfwWindow.SetFramebufferSizeCallback(mng.resizeCallback)
	glfwWindow.SetMaximizeCallback(mng.maximizeCallback)
	glfwWindow.SetIconifyCallback(mng.iconifyCallback)
	glfwWindow.SetCloseCallback(mng.closeCallback)
	glfwWindow.SetFocusCallback(mng.focusCallback)

	// Keyboard Callback
	glfwWindow.SetKeyCallback(mng.keyCallBack)

	// Mouse Callback
	glfwWindow.SetCursorPosCallback(mng.mouseCursorCallback)
	glfwWindow.SetMouseButtonCallback(mng.mouseButtonCallback)
	glfwWindow.SetScrollCallback(mng.mouseScrollCallback)
}

// Release enlève les différentes fonctions de traitement des évènements.
func (mng *EventManager) Release() {
	if mng.window == nil || mng.window.glfwWindow == nil {
		return
	}
	glfwWindow := mng.window.glfwWindow

	// Window Callback
	glfwWindow.SetFramebufferSizeCallback(nil)
	glfwWindow.SetMaximizeCallback(nil)
	glfwWindow.SetIconifyCallback(nil)
	glfwWindow.SetCloseCallback(nil)
	glfwWindow.SetFocusCallback(nil)

	// Keyboard Callback
	glfwWindow.SetKeyCallback(nil)

	// Mouse Callback
	glfwWindow.SetCursorPosCallback(nil)
	glfwWindow.SetMouseButtonCallback(nil)
	glfwWindow.SetScrollCallback(nil)
}

func (mng *EventManager) ProcessEvents(processor EventProcessor) {
	if mng.processor == nil {
		mng.processor = processor
		defer func() {
			mng.processor = nil
		}()
		glfw.PollEvents()
	}
}

//////////////////////////////// KEYBOARD

func (mng *EventManager) keyCallBack(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, _ glfw.ModifierKey) {
	event := &KeyboardEvent{
		keyboard: &mng.window.keyboard,
		key:      key,
		scanCode: scancode,
		action:   convertGlfwActionToKeyAction(action),
	}
	mng.window.keyboard.updateKeyboardWithKeyEvent(event)
	mng.processor(event)
}

//////////////////////////////// MOUSE

func (mng *EventManager) mouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	event := &MouseButtonEvent{
		mouseEvent: mouseEvent{mouse: &mng.window.mouse},
		button:     button,
		action:     convertGlfwActionToMouseButtonAction(action),
	}
	mng.window.mouse.updateWithMouseButtonEvent(event)
	mng.processor(event)
}

func (mng *EventManager) mouseCursorCallback(_ *glfw.Window, xpos float64, ypos float64) {
	event := &MouseMoveEvent{
		mouseEvent: mouseEvent{mouse: &mng.window.mouse},
		pos:        mgl64.Vec2{xpos, ypos},
		shift:      mgl64.Vec2{xpos - mng.window.mouse.pos.X(), ypos - mng.window.mouse.pos.Y()},
	}
	mng.window.mouse.updateWithMouseMoveEvent(event)
	mng.processor(event)
}

func (mng *EventManager) mouseScrollCallback(_ *glfw.Window, horizontalShift float64, verticalShift float64) {
	event := &MouseScrollEvent{
		mouseEvent: mouseEvent{mouse: &mng.window.mouse},
		offset:     mgl64.Vec2{mng.window.mouse.scrollWheels.X() + horizontalShift, mng.window.mouse.scrollWheels.Y() + verticalShift},
		shift:      mgl64.Vec2{horizontalShift, verticalShift},
	}
	mng.window.mouse.updateWithMouseScrollEvent(event)
	mng.processor(event)
}

func (mng *EventManager) resizeCallback(_ *glfw.Window, width int, height int) {
	//gl.Viewport(0, 0, int32(width), int32(height))
	event := WindowEvent{
		window: mng.window,
		width:  width,
		heigth: height,
		action: WINDOW_RESIZED,
	}
	//gl.Viewport(0, 0, int32(width), int32(height))
	mng.processor(&event)
}

//////////////////////////////// WINDOW

func (mng *EventManager) closeCallback(glfwWindow *glfw.Window) {
	width, height := glfwWindow.GetFramebufferSize()
	event := WindowEvent{
		window: mng.window,
		width:  width,
		heigth: height,
		action: WINDOW_CLOSED,
	}
	mng.processor(&event)
}

func (mng *EventManager) maximizeCallback(glfwWindow *glfw.Window, maximized bool) {
	var action WindowAction
	width, height := glfwWindow.GetFramebufferSize()
	if maximized {
		action = WINDOW_MAXIMIZED
	} else {
		action = WINDOW_RESTORED
	}
	//gl.Viewport(0, 0, int32(width), int32(height))
	event := WindowEvent{
		window: mng.window,
		width:  width,
		heigth: height,
		action: action,
	}
	mng.processor(&event)
}

func (mng *EventManager) iconifyCallback(glfwWindow *glfw.Window, iconified bool) {
	var action WindowAction
	width, height := glfwWindow.GetFramebufferSize()
	if iconified {
		action = WINDOW_ICONIFIED
	} else {
		action = WINDOW_RESTORED
		//gl.Viewport(0, 0, int32(width), int32(height))
	}
	event := WindowEvent{
		window: mng.window,
		width:  width,
		heigth: height,
		action: action,
	}
	mng.processor(&event)
}

func (mng *EventManager) focusCallback(glfwWindow *glfw.Window, focused bool) {
	width, height := glfwWindow.GetFramebufferSize()
	var action WindowAction
	if focused {
		action = WINDOW_FOCUS
	} else {
		action = WINDOW_UNFOCUS
	}
	event := WindowEvent{
		window: mng.window,
		width:  width,
		heigth: height,
		action: action,
	}
	mng.processor(&event)
}
