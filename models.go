package main

// import rl "github.com/gen2brain/raylib-go/raylib"

type BoardCoord struct {
	Row    int
	Column int
}

type Game struct {
	Scenes            map[string]*Scene
	Data              map[string]interface{}
	screenWidth       int32
	screenHeight      int32
	SidebarWidth      int32
	Run               *Run
	Counter           int32
	Seconds           int32
	ScreenSkip        bool
	WindowOpen        bool
	Technology        map[string]*Technology
	InitialData       map[string]InitialData
	GameOver          bool
	GameOverTriggered bool
}

type Run struct {
	Technology             []*Technology
	People                 []Person
	Products               map[ProductType]*Product
	Money                  float32
	Yield                  float32
	Productivity           float32
	EndRoundMoney          float32
	MoneyRequirement       float32
	MoneyRequirementStart  float32
	MoneyRequirementRate   float32
	CurrentRound           int
	CurrentYear            int
	CurrentSeason          Season
	CurrentRoundShopPlants []*Technology
	NextSeason             Season
	EventChoices           []Event
	Events                 []Event
	PossibleEvents         []Event
	EventTracker           EventTracker
	TechnologySpaces       []*TechnologySpace
	ActionsRemaining       int
	ActionsMaximum         int
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
	Technology     *Technology
	TechnologyType TechnologyType
	Row            int
	Column         int
	Width          int // in tiles
	Height         int // in tiles
	IsFilled       bool
	Active         bool // whether the game displays or not
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
	OnRoundEnd      func(*Game, *Technology)
	RoundEndProduce func(*Game, *Technology) float32
	ShopButton      func(*Game) *ShopButton
	ToBeDeleted     bool
	Space           *TechnologySpace
	ReadyToHarvest  bool
	ReadyToTouch    bool
	TempYield       float32
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
	LandClearageTriggered bool
	LandClearageFinished  bool
}

type EventJSON struct {
	Name        string
	Description string
}

type Event struct {
	RoundIndex  int
	Name        string
	Description string
	Effects     []Effect
	BlankEvent  bool
	OnTrigger   func(*Game)
}
type Effect struct {
	ProductImpacted ProductType
	IsPriceChange   bool
	PriceChange     float32 // percentage
	EventTrigger    func(*Game)
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
	Name           string         `json:"name"`
	Price          float32        `json:"price"`
	ProductType    string         `json:"productType"`
	TechnologyType string         `json:"technologyType"`
	CostMoney      float32        `json:"costMoney"`
	Production     float32        `json:"production"`
	Rarity         string         `json:"rarity"`
	Description    string         `json:"description"`
	TileConfig     TechTileConfig `json:"tile"`
	ShopIcon       string         `json:"shopIcon"`
}

type TechTileConfig struct {
	ID        string `json:"id"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	FillSpace bool   `json:"fillSpace"`
}
