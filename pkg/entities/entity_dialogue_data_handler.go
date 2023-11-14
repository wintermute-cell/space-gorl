package entities

import (
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/logging"
	"encoding/json"
	"io/ioutil"
)

// Compares the EmotionalProfile to the 'reference' EmotionalProfile. A
// returned value around 0 indicates a good match, while a value < 0 means that
// the EmotionalProfile even surpasses the reference EmotionalProfile. In
// short, the higher the returned value, the worse the fit.
func (ep EmotionalProfile) Difference(reference EmotionalProfile) float32 {
	diff := float32(0)
	diff += reference.Tradition - ep.Tradition
	diff += reference.Security - ep.Security
	diff += reference.Conformity - ep.Conformity
	diff += reference.Achievement - ep.Achievement
	diff += reference.Power - ep.Power
	diff += reference.Hedonism - ep.Hedonism
	diff += reference.Stimulation - ep.Stimulation
	diff += reference.SelfDirection - ep.SelfDirection
	diff += reference.Universalism - ep.Universalism
	diff += reference.Benevolence - ep.Benevolence
	return diff
}

// Adds the EmotionalProfile x to the EmotionalProfile ep and returns the
// result as a new EmotionalProfile.
func (ep EmotionalProfile) Add(x EmotionalProfile) EmotionalProfile {
	new_ep := EmotionalProfile{}
	new_ep.Tradition = ep.Tradition + x.Tradition
	new_ep.Security = ep.Security + x.Security
	new_ep.Conformity = ep.Conformity + x.Conformity
	new_ep.Achievement = ep.Achievement + x.Achievement
	new_ep.Power = ep.Power + x.Power
	new_ep.Hedonism = ep.Hedonism + x.Hedonism
	new_ep.Stimulation = ep.Stimulation + x.Stimulation
	new_ep.SelfDirection = ep.SelfDirection + x.SelfDirection
	new_ep.Universalism = ep.Universalism + x.Universalism
	new_ep.Benevolence = ep.Benevolence + x.Benevolence
	return new_ep
}

type EmotionalProfile struct {
	Tradition     float32 `json:"Tradition,omitempty"`     // Respect for cultural, family, and religious traditions.
	Security      float32 `json:"Security,omitempty"`      // Safety, harmony, and stability of society, of relationships, and of self.
	Conformity    float32 `json:"Conformity,omitempty"`    // Restraint of actions, inclinations, and impulses likely to upset or harm others and violate social expectations or norms.
	Achievement   float32 `json:"Achievement,omitempty"`   // Personal success through demonstrating competence according to social standards.
	Power         float32 `json:"Power,omitempty"`         // Social status and prestige, control or dominance over people and resources.
	Hedonism      float32 `json:"Hedonism,omitempty"`      // Pleasure or sensuous gratification for oneself.
	Stimulation   float32 `json:"Stimulation,omitempty"`   // Excitement, novelty, and challenge in life.
	SelfDirection float32 `json:"SelfDirection,omitempty"` // Independent thought and action; choosing, creating, exploring.
	Universalism  float32 `json:"Universalism,omitempty"`  // Understanding, appreciation, tolerance, and protection for the welfare of all people and for nature.
	Benevolence   float32 `json:"Benevolence,omitempty"`   // Preserving and enhancing the welfare of those with whom one is in frequent personal contact (the 'in-group').
}

type DialogueCondition struct {
	Items []string `json:"Items,omitempty"`
}

type DialogueResponse struct {
	Text      string            `json:"text"`
	Condition DialogueCondition `json:"condition"`
	Influence EmotionalProfile  `json:"influence"`
	Value     float32           `json:"value"`
}

type DialogueNode struct {
	Text       string              `json:"text"`
	Responses  []DialogueResponse  `json:"responses"`
	Conditions []DialogueCondition `json:"conditions"`
	Ep         EmotionalProfile    `json:"Ep"`
	Children   []*DialogueNode      `json:"Children"`
	WasUsed    bool                `json:"WasUsed"`
}

// DialogueDataHandler Entity
type DialogueDataHandlerEntity struct {
	// Required fields
	proto.BaseEntity

	dialogueTreeRoot *DialogueNode

	// Custom Fields
	// Add fields here for any state that the entity should keep track of
	// ...
}

func NewDialogueDataHandlerEntity() *DialogueDataHandlerEntity {
	new_ent := &DialogueDataHandlerEntity{}
	return new_ent
}

func (ent *DialogueDataHandlerEntity) GetRootNode() *DialogueNode {
    return ent.dialogueTreeRoot
}

func (ent *DialogueDataHandlerEntity) Init() {
	// Required initialization
	ent.BaseEntity.Init()

	file, err := ioutil.ReadFile("dialogue/dialogue.json")
	if err != nil {
		logging.Fatal("Error while loading dialogue.json file: %v", err)
	}

    ent.dialogueTreeRoot = &DialogueNode{}
	err = json.Unmarshal(file, ent.dialogueTreeRoot)
	if err != nil {
		logging.Fatal("Error while parsing dialogue.json file: %v", err)
	}

	// Initialization logic for the entity
	// ...
}

func (ent *DialogueDataHandlerEntity) Deinit() {
	// Required de-initialization
	ent.BaseEntity.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *DialogueDataHandlerEntity) Update() {
	// Required update
	ent.BaseEntity.Update()

	// Update logic for the entity
	// ...
}

func (ent *DialogueDataHandlerEntity) Draw() {
	// Draw logic for the entity
	// ...
}
