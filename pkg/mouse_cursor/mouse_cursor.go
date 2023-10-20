package mousecursor

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type cursor struct {
	sprite rl.Texture2D
}

var c cursor

func Init() {
	c = cursor{
		sprite: rl.LoadTexture("sprites/player/crosshair.png"),
	}
}

func Draw() {
	p := rl.GetMousePosition()
	rl.DrawTexturePro(
		c.sprite,
		rl.NewRectangle(0, 0, float32(c.sprite.Width), float32(c.sprite.Height)),
		rl.NewRectangle(p.X, p.Y, float32(c.sprite.Width)*2, float32(c.sprite.Height)*2),
		rl.NewVector2(float32(c.sprite.Width)*2/2, float32(c.sprite.Height)*2/2),
		0.0,
		rl.White)
	//rl.DrawCircleV(rl.GetMousePosition(), 2, rl.Gray)
}
