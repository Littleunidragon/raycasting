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
	radius       = 3
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
	dir float64
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
	if ebiten.IsKeyPressed(ebiten.KeyA) {
			p.dir -= 0.05
		if p.dir< 0 {
			p.dir += 2 * math.Pi
		}	
		p.vel.x = math.Cos(p.dir) * 5
		p.vel.y = math.Sin(p.dir) * 5
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.dir += 0.05
	if p.dir > 2 {
		p.dir -= 2 * math.Pi
	}	
	p.vel.x = math.Cos(p.dir) * 5
	p.vel.y = math.Sin(p.dir) * 5
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

// func (p *player)raycast(screen *ebiten.Image) {
// 	mx, my := p.pos.x, p.pos.y
// 	mvx, mvy := p.vel.x, p.vel.y
// 	for r := 0 ; r < 1; r++ {
// 		ebitenutil.DrawLine(screen, p.pos.x, p.pos.y, mx+ mvx * 200,my + mvy * 200, color.RGBA{100,100,100,255})
// 	}
// }

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
	ebitenutil.DrawLine(screen,g.p.pos.x, g.p.pos.y, g.p.pos.x + g.p.vel.x *5, g.p.pos.y + g.p.vel.y *5, color.RGBA{255,0,0,255})
	//g.p.raycast(screen)	
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
