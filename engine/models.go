package engine

import (
	"github.com/gen2brain/raylib-go/raylib"
	//	"log"
)

type GameInterface interface {
	GetRun() interface{}
	GetScenes() map[string]*Scene
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
	Name                   string
	Active                 bool
	AutoDisable            bool
	DrawScene              func(GameInterface)
	UpdateScene            func(GameInterface)
	Buttons                []Button
	skip                   bool
	Data                   map[string]interface{}
	Camera                 rl.Camera2D
	Windows                map[string]*Window
	RenderMenu             bool
	Messages               []Message
	KeyBindingFunctions    map[string]func(GameInterface)
	KeyBindings            map[string]*KeyBinding
	Components             []UIComponent
	SelectedComponentIndex int
}

type Window struct {
	Name                   string
	DrawWindow             func(GameInterface, *Window)
	Display                bool
	Buttons                []Button
	Components             []UIComponent
	SelectedComponentIndex int
	Rectangle              rl.Rectangle
	X                      int32
	Y                      int32
	Width                  int32
	Height                 int32
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
	OnPress      func(GameInterface)
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
