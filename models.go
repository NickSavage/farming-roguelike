package main

import rl "github.com/gen2brain/raylib-go/raylib"

type BoardCoord struct {
	Row    int
	Column int
}

type Game struct {
	Scenes       map[string]*Scene
	Data         map[string]interface{}
	screenWidth  int32
	screenHeight int32
	Run          *Run
	Counter      int32
	Seconds      int32
	ScreenSkip   bool
	WindowOpen   bool
	Technology   map[string]*Technology
}

type BoardSquare struct {
	Tile         Tile
	TileType     string
	Row          int
	Column       int
	Width        int // in tiles
	Height       int // in tiles
	Skip         bool
	Occupied     bool
	MultiSquare  bool
	Technology   *Technology
	IsTechnology bool
	HoverActive  bool
	IsTree       bool
}

type BoardRightClickMenu struct {
	Rectangle   rl.Rectangle
	BoardSquare *BoardSquare
	Items       []BoardMenuItem
}

type BoardMenuItem struct {
	Rectangle       rl.Rectangle
	Text            string
	OnClick         func(*Game)
	CheckIsDisabled func(*Game, *BoardSquare) bool
}

type Technology struct {
	Name              string
	Description       string
	Square            BoardSquare
	Cost              float32
	CanBeBuilt        func(*Game) bool
	OnBuild           func(*Game, *Technology) error
	RoundHandler      []TechnologyRoundHandler
	RoundCounterMax   int
	RoundCounter      int
	RoundHandlerIndex int
	Redraw            bool
	ShowEndRound      bool
}

type TechnologyRoundHandler struct {
	Season        Season
	CostActions   float32
	CostMoney     float32
	OnRoundEnd    func(*Game, *Technology)
	RoundEndText  func(*Game, *Technology) string
	RoundEndValue func(*Game, *Technology) float32
}

type Person struct {
}

type Run struct {
	Technology            []*Technology
	People                []Person
	Products              map[string]*Product
	Money                 float32
	Productivity          float32
	EndRoundMoney         float32
	RoundActions          float32
	RoundActionsRemaining float32
	CurrentRound          int
	CurrentSeason         Season
	Events                []Event
}

type Event struct {
	RoundIndex int
	Name       string
	Effects    []Effect
}
type Effect struct {
	ProductImpacted string
	IsPriceChange   bool
	PriceChange     float32 // percentage
}

type Product struct {
	Name     string
	Quantity float32
	Price    float32
}

// Define the type for the enum
type Season int

// Declare constants using iota
const (
	Spring Season = iota
	Summer
	Autumn
	Winter
)

func (s Season) String() string {
	return [...]string{"Spring", "Summer", "Autumn", "Winter"}[s]
}

func (s *Season) Next() {
	*s = (*s + 1) % 4 // Cycle back to Spring after Winter
}
