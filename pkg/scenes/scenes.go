package scenes

// Scene is an interface that every scene in the game should implement
type Scene interface {
	Init()
	Deinit()
	DrawGUI()
	Draw()
}
