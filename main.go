package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 630
	screenHeight = screenWidth
	wallX        = screenHeight / 21
	wallY        = wallX
	radius       = 5
)

var (
	worldmap = [21][21]int{
		{1, 1, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 1, 1, 1, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 0, 1, 0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0},
		{0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 1},
		{0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 2, 0, 1, 1, 0, 0, 0, 1, 1, 0, 1},
		{0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1},
		{1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1, 1, 0},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 1, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 1},
		{1, 0, 1, 1, 1, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 1},
		{1, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0},
	}
)

type Coord struct {
	x, y float64
}

type player struct {
	pos Coord
	vel Coord
}

type game struct {
	p    player
	last time.Time
}

// draw map
func mapReader(screen *ebiten.Image) {
	for i := range worldmap {
		for j := range worldmap[i] {
			if worldmap[i][j] == 1 {
				ebitenutil.DrawRect(screen, float64(wallX*j), float64(wallY*i), wallX, wallY, color.RGBA{0, 0, 0, 255})
			}
		}
	}
}

// player starting position
func (p *player) newPlayer() {
	for i := range worldmap {
		for j := range worldmap[i] {
			if worldmap[i][j] == 2 {
				p.pos.x = float64(wallX * j)
				p.pos.y = float64(wallY * i)
			}
		}
	}
}

// draw player
func (p *player) drawPlayer(screen *ebiten.Image) {
	ebitenutil.DrawCircle(screen, p.pos.x, p.pos.y, 5, color.RGBA{255, 0, 0, 255})
}

// move player
func (p *player) movePlayer(dtms float64) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		p.pos.y -= p.vel.y * dtms
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		p.pos.y += p.vel.y * dtms
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.pos.x -= p.vel.x * dtms
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.pos.x += p.vel.x * dtms
	}
}

// outer collision
func (p *player) oCollision() {
	switch {
	case p.pos.x >= float64(screenWidth)-radius:
		p.pos.x = float64(screenWidth) - radius
		p.vel.x *= -1
	case p.pos.x <= radius:
		p.pos.x = radius
		p.vel.x *= -1
	case p.pos.y >= float64(screenHeight)-radius:
		p.pos.y = float64(screenHeight) - radius
		p.vel.y *= -1
	case p.pos.y <= radius:
		p.pos.y = radius
		p.vel.y *= -1
	}
}

func DrawLineDDA(screen *ebiten.Image, p0x, p0y, p1x, p1y float64, color color.Color) {
	if math.Abs(p1x-p0x) >= math.Abs(p1y-p0y) {
		if p0x > p1x {
			p0x, p1x = p1x, p0x
			p0y, p1y = p1y, p0y
		}
		y := p0y
		for x := p0x; x <= p1x; x++ {
			screen.Set(int(x), int(y), color)
			y += (p1y - p0y) / (p1x - p0x)
		}
	} else {
		if p0y > p1y {
			p0x, p1x = p1x, p0x
			p0y, p1y = p1y, p0y
		}
		x := p0x
		for y := p0y; y <= p1y; y++ {
			screen.Set(int(x), int(y), color)
			x += (p1x - p0x) / (p1y - p0y)
		}
	}
}

// rotate point
func rotatePoint(p *Coord, angle float64) *Coord {
	s := math.Sin(angle)
	c := math.Cos(angle)
	xnew := p.x*c - p.y*s
	ynew := p.x*s + p.y*c
	return &Coord{xnew, ynew}
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }
func (g *game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	g.p.vel.x = 0.1
	g.p.vel.y = 0.1
	g.p.movePlayer(dt)
	g.p.oCollision()
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{255, 255, 255, 255})
	mapReader(screen)
	g.p.drawPlayer(screen)
	DrawLineDDA(screen, g.p.pos.x, g.p.pos.y, g.p.pos.x+g.p.vel.x+100.1, g.p.pos.y+g.p.vel.y+100.1, color.RGBA{255, 0, 0, 255})
}
func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Into the encoded world")
	var p player
	p.newPlayer()
	if err := ebiten.RunGame(&game{p, time.Now()}); err != nil {
		log.Fatal(err)
	}
}
