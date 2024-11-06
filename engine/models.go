package engine

import (
	"github.com/gen2brain/raylib-go/raylib"
	//	"log"
)

type GameInterface interface {
	GetRun() interface{}
	// Scenes() map[string]*Scene
}

type Game struct {
	GameScenes map[string]*Scene
}

func (g Game) GetRun() interface{} {
	return "not implemented"
}

func (g *Game) Scenes() map[string]*Scene {
	return g.GameScenes
}

type UIComponent interface {
	Render()
	OnClick()
	Rect() rl.Rectangle
}

type Scene struct {
	Name                string
	Active              bool
	AutoDisable         bool
	DrawScene           func(*Game)
	UpdateScene         func(*Game)
	Buttons             []Button
	skip                bool
	Data                map[string]interface{}
	Camera              rl.Camera2D
	Windows             map[string]*Window
	RenderMenu          bool
	Messages            []Message
	KeyBindingFunctions map[string]func(*Game)
	KeyBindings         map[string]*KeyBinding
	Components          []UIComponent
}

type Window struct {
	Name       string
	DrawWindow func(GameInterface, *Window)
	Display    bool
	Buttons    []Button
	Components []UIComponent
}
type Message struct {
	Text  string
	Vec   rl.Vector2
	Timer int
}
type KeyBinding struct {
	Current      int32
	Default      int32
	Name         string
	FunctionName string
	Scene        string
	Configurable bool
	OnPress      func(*Game)
}

type KeyBindingJSON struct {
	Default      int32  `json:"default"`
	Name         string `json:"name"`
	FunctionName string `json:"functionName"`
	Scene        string `json:"scene"`
	Configurable bool   `json:"configurable"`
}

type BoardCoord struct {
	Row    int
	Column int
}
