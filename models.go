package main

import (
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
	UnlockSave        []UnlockSave
	Unlocks           map[string]*Unlock
	ProductStats      map[ProductType]*ProductStat
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

type RunSaveFile struct {
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
	Seeds          []TechnologySave
}

type TechnologyType = string

const (
	PlantSpace     TechnologyType = "PlantSpace"
	Seed           TechnologyType = "Seed"
	Field          TechnologyType = "Field"
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
	Type          ProductType
	Quantity      float32
	Price         float32
	Yield         float32
	TotalProduced float32
	TotalEarned   float32
}

type ProductType string

const (
	Chicken ProductType = "Chicken"
	Wheat   ProductType = "Wheat"
	Flour   ProductType = "Flour"
	Potato  ProductType = "Potato"
	Carrot  ProductType = "Carrot"
	Solar   ProductType = "Solar"
	Cow     ProductType = "Cow"
	Beef    ProductType = "Beef"
)

type ProductStat struct {
	ProductType          ProductType
	MaxProduction        float32
	TotalProduction      float32
	CurrentRunProduction float32
	MaxEarned            float32
	TotalEarned          float32
	CurrentRunEarned     float32
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
	CostActions    int    `json:"costActions"`
	OtherCost      bool   `json:"otherCost"`
	Dependency     string `json:"dependency"`
	TechnologyName string
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
	Technology                   *Technology
	Unlocked                     bool
	CostActions                  int
	OtherCost                    bool
	OtherCostFunction            func(*Game) bool
	OtherCostDescriptionFunction func(*Game) string
	DependencyName               string
	DependencyMet                bool
	RunSpentActions              int
}
