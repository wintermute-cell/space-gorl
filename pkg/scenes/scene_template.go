package scenes

import (
	"cowboy-gorl/pkg/entities/proto"
)

// This checks at compile time if the interface is implemented
var _ Scene = (*TemplateScene)(nil)

// Template Scene
type TemplateScene struct {
	// Required fields
	entity_manager *proto.EntityManager

	// Custom Fields
	// Add fields here for any state that the scene should keep track of
	// ...
}

func (scn *TemplateScene) Init() {
	// Required initialization
	scn.entity_manager = proto.NewEntityManager()

	// Initialization logic for the scene
	// ...
}

func (scn *TemplateScene) Deinit() {
	// De-initialization logic for the scene
}

func (scn *TemplateScene) DrawGUI() {
	// Draw the GUI for the scene
}

func (scn *TemplateScene) Draw() {
	// Draw the scene
}
