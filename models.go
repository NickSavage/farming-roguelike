package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Game struct {
	Scenes       map[string]*Scene
	Data         map[string]interface{}
	screenWidth  int32
	screenHeight int32
	Run          *Run
	Counter      int32
	Seconds      int32
}

type BoardSquare struct {
	Tile        Tile
	TileType    string
	Row         int
	Column      int
	Width       int // in tiles
	Height      int // in tiles
	Skip        bool
	Occupied    bool
	MultiSquare bool
	Technology  *Technology
}

type BoardRightClickMenu struct {
	Rectangle   rl.Rectangle
	BoardSquare *BoardSquare
}

type Technology struct {
	Name            string
	Description     string
	Tile            BoardSquare
	Cost            float32
	OnRoundEnd      func(*Game, *Technology)
	OnBuild         func(*Game, *Technology)
	RoundEndText    func(*Game, *Technology) string
	RoundEndValue   func(*Game, *Technology) float32
	RoundCounterMax int
	RoundCounter    int
}

type Person struct {
}

type Run struct {
	Technology            []*Technology
	People                []Person
	Money                 float32
	Productivity          float32
	EndRoundMoney         float32
	RoundActions          int
	RoundActionsRemaining int
	CurrentRound          int
	Events                []Event
}

type Event struct {
	RoundIndex int
	Name       string
	Effects    []Effect
}

type Effect struct{}
