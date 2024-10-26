package main

// import rl "github.com/gen2brain/raylib-go/raylib"

type BoardCoord struct {
	Row    int
	Column int
}

type Game struct {
	Scenes       map[string]*Scene
	Data         map[string]interface{}
	screenWidth  int32
	screenHeight int32
	SidebarWidth int32
	Run          *Run
	Counter      int32
	Seconds      int32
	ScreenSkip   bool
	WindowOpen   bool
	Technology   map[string]*Technology
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
}

type Technology struct {
	Name              string
	ProductType       ProductType
	TechnologyType    TechnologyType
	Tile              Tile
	TileWidth         int
	TileHeight        int
	TileFillSpace     bool
	Description       string
	Square            BoardSquare
	CostMoney         float32
	OnBuild           func(*Game, *Technology) error
	RoundHandler      []TechnologyRoundHandler
	RoundCounterMax   int
	RoundCounter      int
	RoundHandlerIndex int
	Redraw            bool
	ShowEndRound      bool
	ToBeDeleted       bool
	Space             *TechnologySpace
}

type TechnologyType int

const (
	PlantSpace TechnologyType = iota
	BuildingSpace
)

type TechnologyRoundHandler struct {
	Season          Season
	CostMoney       float32
	OnRoundEnd      func(*Game, *Technology)
	RoundEndProduce func(*Game, *Technology) float32
}

type Person struct {
}

type Run struct {
	Technology       []*Technology
	People           []Person
	Products         map[ProductType]*Product
	Money            float32
	Yield            float32
	Productivity     float32
	EndRoundMoney    float32
	CurrentRound     int
	CurrentSeason    Season
	Events           []Event
	TechnologySpaces []*TechnologySpace
}

type Event struct {
	RoundIndex int
	Name       string
	Effects    []Effect
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
