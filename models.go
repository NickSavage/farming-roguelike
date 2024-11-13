package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"nsavage/farming-roguelike/engine"
)

type Game struct {
	Scenes            map[string]*engine.Scene
	Data              map[string]interface{}
	screenWidth       int32
	screenHeight      int32
	SidebarWidth      int32
	Run               *Run
	Counter           int32
	Seconds           int32
	ScreenSkip        bool
	ButtonSkip        int32
	WindowOpen        bool
	Technology        map[string]*Technology
	InitialData       map[string]InitialData
	GameOver          bool
	GameOverTriggered bool
	KeyBindingJSONs   []KeyBindingJSON
	ExistingSave      bool
	ActiveRun         bool
	UnlockBaseData    []UnlockJSON
	Unlocks           map[string]*Unlock
}

type Run struct {
	Game                      *Game
	Technology                []*Technology
	Products                  map[ProductType]*Product
	Money                     float32
	Yield                     float32
	Productivity              float32
	EndRoundMoney             float32
	MoneyRequirement          float32
	MoneyRequirementStart     float32
	MoneyRequirementRate      float32
	CurrentRound              int
	CurrentYear               int
	CurrentSeason             Season
	CurrentRoundShopPlants    []*Technology
	CurrentRoundShopBuildings []*Technology
	NextSeason                Season
	EventChoices              []Event
	Events                    []Event
	PossibleEvents            []Event
	triggerFunctions          map[string]func(*Game)
	EventTracker              map[string]bool // track if its been called or not
	TechnologySpaces          []*TechnologySpace
	ActionsRemaining          int
	ActionsMaximum            int
	AutoSellRoundEnd          bool //whether the player wants to autosell
}

type SaveFile struct {
	Money                 float32                  `json:"money"`
	MoneyRequirementStart float32                  `json:"money_requirement_start"`
	MoneyRequirementRate  float32                  `json:"money_requirement_rate"`
	Yield                 float32                  `json:"yield"`
	Productivity          float32                  `json:"productivity"`
	CurrentRound          int                      `json:"current_round"`
	CurrentYear           int                      `json:"current_year"`
	CurrentSeason         Season                   `json:"current_season"`
	ActionsRemaining      int                      `json:"actions_remaining"`
	ActionsMaximum        int                      `json:"actions_maximum"`
	EventTracker          map[string]bool          `json:"event_tracker"`
	Technology            []TechnologySave         `json:"technology_save"`
	Products              map[ProductType]*Product `json:"products"`
	Events                []EventSave              `json:"event_save"`
	Unlocks               []UnlockSave             `json:"unlock_save"`
}

type BoardSquare struct {
	Tile              Tile
	TileType          string
	Row               int
	Column            int
	Width             int // in tiles
	Height            int // in tiles
	Skip              bool
	Occupied          bool
	MultiSquare       bool
	IsTechnologySpace bool
	TechnologySpace   *TechnologySpace
	HoverActive       bool
	IsTree            bool
}

type TechnologySpace struct {
	Game             *Game
	ID               int
	Technology       *Technology
	TechnologyType   TechnologyType
	Row              int
	Column           int
	Width            int // in tiles
	Height           int // in tiles
	IsFilled         bool
	Active           bool // whether the game displays or not
	SelectDirections engine.SelectDirections
	Selected         bool
}

func (space *TechnologySpace) Render() {

	g := space.Game
	scene := space.Game.Scenes["Board"]
	if !space.Active {
		return
	}
	boxColor := rl.Blue
	if space.Selected {
		boxColor = rl.Green
	}

	vec := g.GetVecFromCoords(engine.BoardCoord{Row: space.Row, Column: space.Column})
	x := vec.X
	y := vec.Y
	width := float32(space.Width * TILE_WIDTH)
	height := float32(space.Height * TILE_HEIGHT)
	rect := rl.NewRectangle(x, y, width, height)
	rl.DrawRectangleRec(rect, boxColor)
	if !space.IsFilled {
		return
	}
	mousePosition := rl.GetMousePosition()

	if !g.WindowOpen && rl.CheckCollisionPointRec(mousePosition, rect) {

		space.Technology.Tile.Color = rl.Green
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			result := g.HandleClickTech(space.Technology)
			message := engine.Message{
				Text:  result,
				Vec:   rl.Vector2{X: x, Y: y},
				Timer: 30,
			}
			scene.Messages = append(scene.Messages, message)
		}

	} else {
		if space.Technology.ReadyToHarvest {
			space.Technology.Tile.Color = rl.Blue
		} else if space.Technology.ReadyToTouch {
			space.Technology.Tile.Color = rl.Red
		} else {
			space.Technology.Tile.Color = rl.White
		}
	}

	if space.Technology.TileFillSpace {
		for i := range space.Width {
			for j := range space.Height {
				DrawTile(
					space.Technology.Tile,
					float32(float32(x)+float32(i*TILE_WIDTH)),
					float32(float32(y)+float32(j*TILE_HEIGHT)),
				)
			}
		}

	} else {
		DrawTile(space.Technology.Tile, float32(x), float32(y))

	}

}

func (space *TechnologySpace) OnClick() {
	if space.IsFilled {
		space.Technology.OnClick(space.Game, space.Technology)
	}

}

func (space *TechnologySpace) Rect() rl.Rectangle {

	vec := space.Game.GetVecFromCoords(engine.BoardCoord{Row: space.Row, Column: space.Column})
	return rl.Rectangle{
		X:      vec.X,
		Y:      vec.Y,
		Width:  float32(space.Width * TILE_WIDTH),
		Height: float32(space.Height) * TILE_HEIGHT,
	}

}
func (space *TechnologySpace) Select() {
	log.Printf("selected")
	space.Selected = true

}

func (space *TechnologySpace) Unselect() {
	space.Selected = false
}

func (space *TechnologySpace) IsSelected() bool {
	return space.Selected
}

func (space *TechnologySpace) Directions() *engine.SelectDirections {
	return &space.SelectDirections
}

type Technology struct {
	Name            string
	ProductType     ProductType
	TechnologyType  TechnologyType
	Rarity          string
	Tile            Tile
	TileWidth       int
	TileHeight      int
	TileFillSpace   bool
	ShopIcon        string
	Description     string
	Square          BoardSquare
	CostMoney       float32
	CostActions     int
	InitialPrice    float32
	BaseProduction  float32
	CanBuild        func(*Game, *Technology) bool
	OnBuild         func(*Game, *Technology) error
	OnClick         func(*Game, *Technology) string
	ShopOnClick     func(*Game)
	OnRoundEnd      func(*Game, *Technology)
	RoundEndProduce func(*Game, *Technology) float32
	ShopButton      func(*Game) *ShopButton
	ToBeDeleted     bool
	Space           *TechnologySpace
	ReadyToHarvest  bool
	ReadyToTouch    bool
	TempYield       float32
	Input           Input
	Unlocked        bool
}

type Input struct {
	ProductType    ProductType
	MaximumInput   float32
	OutputPerInput float32
}

type TechnologySave struct {
	Name           string
	ReadyToHarvest bool
	ReadyToTouch   bool
	TempYield      float32
	SpaceID        int
}

type TechnologyType = string

const (
	PlantSpace     TechnologyType = "PlantSpace"
	BuildingSpace  TechnologyType = "BuildingSpace"
	CellTowerSpace TechnologyType = "CellTowerSpace"
)

type Person struct {
}

type EventTracker struct {
	LandClearage bool
	HireHelp     bool
	CellTower    bool
}

type EventJSON struct {
	Name        string
	Description string
	Repeatable  bool
	CostMoney   float32
	Severity    float32
}

type Event struct {
	RoundIndex  int
	Name        string
	Description string
	Effects     []Effect
	BlankEvent  bool
	OnTrigger   func(*Game)
	Repeatable  bool
	CostMoney   float32
	Severity    float32
}

type EventSave struct {
	RoundIndex  int
	Name        string
	Description string
	Effects     []Effect
	BlankEvent  bool
	Repeatable  bool
}

type Effect struct {
	ProductImpacted ProductType
	IsPriceChange   bool
	PriceChange     float32 // percentage
}

type Product struct {
	Type        ProductType
	Quantity    float32
	Price       float32
	Yield       float32
	TotalEarned float32
}

type ProductType string

const (
	Chicken ProductType = "Chicken"
	Wheat   ProductType = "Wheat"
	Potato  ProductType = "Potato"
	Carrot  ProductType = "Carrot"
	Solar   ProductType = "Solar"
	Cow     ProductType = "Cow"
	Beef    ProductType = "Beef"
)

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

type InitialData struct {
	Name            string         `json:"name"`
	Price           float32        `json:"price"`
	ProductType     string         `json:"productType"`
	TechnologyType  string         `json:"technologyType"`
	CostMoney       float32        `json:"costMoney"`
	CostActions     int            `json:"costActions"`
	Production      float32        `json:"production"`
	Rarity          string         `json:"rarity"`
	Description     string         `json:"description"`
	TileConfig      TechTileConfig `json:"tile"`
	ShopIcon        string         `json:"shopIcon"`
	Input           Input          `json:"input"`
	Unlock          *UnlockJSON    `json:"unlock"`
	DefaultUnlocked bool
}

type UnlockJSON struct {
	CostActions      int    `json:"costActions"`
	OtherCost        bool   `json:"otherCost"`
	OtherDescription string `json:"otherDescription"`
	Dependency       string `json:"dependency"`
	TechnologyName   string
}

type UnlockSave struct {
	TechnologyName string `json:"technologyName"`
	Unlocked       bool   `json:"Unlocked"`
}

type TechTileConfig struct {
	ID        string `json:"id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	FillSpace bool   `json:"fillSpace"`
}

type Unlock struct {
	Technology        *Technology
	Unlocked          bool
	CostActions       int
	OtherCost         bool
	OtherDescription  string
	OtherCostFunction func(*Game) bool
	DependencyName    string
	DependencyMet     bool
}
