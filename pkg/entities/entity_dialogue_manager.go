package entities

import (
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/logging"
	"cowboy-gorl/pkg/util"
	"math"
)

// DialogueManager Entity
type DialogueManagerEntity struct {
	// Required fields
	proto.BaseEntity

	// Custom Fields
    gameState *GameStateHandlerEntity
    dataHandler *DialogueDataHandlerEntity
    dialoguePanel *DialoguePanelEntity2D
}

func NewDialogueManagerEntity(gameState *GameStateHandlerEntity) *DialogueManagerEntity {
	new_ent := &DialogueManagerEntity{
        gameState: gameState,
    }
	return new_ent
}

func (ent *DialogueManagerEntity) Init() {
	// Required initialization
	ent.BaseEntity.Init()

    ent.dataHandler = NewDialogueDataHandlerEntity()
    gem.AddEntity(ent, ent.dataHandler)

    ent.dialoguePanel = NewDialoguePanelEntity2D(ent)
    gem.AddEntity(ent, ent.dialoguePanel)

    // TODO: remove in the future
    ent.beginDialogueMode()
}

// selects the next dialogue node that was not already used and fits best with
// the current playerEmotionalProfile.
func selectNextNode(root *DialogueNode, playerEmotionalProfile EmotionalProfile) *DialogueNode {
    // If the root is nil or has no children, return nil
    if root == nil || len(root.Children) == 0 {
        return nil
    }

    // Initialize variables for tracking the best node and the minimum difference
    var bestNode *DialogueNode = root
    minDiff := float32(math.MaxFloat32)

    // Function to traverse and check each node
    var traverse func(*DialogueNode)
    traverse = func(node *DialogueNode) {
        // Check if the node was used and has children
        if node.WasUsed {
            for _, child := range node.Children {
                // Consider only children that haven't been used
                if !child.WasUsed {
                    // Calculate the difference with the player's emotional profile
                    diff := util.Abs(child.Ep.Difference(playerEmotionalProfile))
                    // Update the best node if this is a closer match
                    if diff < minDiff {
                        minDiff = diff
                        bestNode = child
                    }
                }
            }
            // Traverse the children of this node
            for _, child := range node.Children {
                traverse(child)
            }
        }
    }

    // Start the traversal from the root
    traverse(root)

    if bestNode == root {
        logging.Warning(
            "Selected dialogue node is equal to root. " +
            "This could indicate the beginning of the dialogue tree, " +
            "or there is no more unused nodes left. Resetting all " +
            "Nodes back to unused")
        resetDialogueNodes(root)
    }

    // Return the best node found
    return bestNode
}

func resetDialogueNodes(node *DialogueNode) {
    // If the node is nil, return immediately
    if node == nil {
        return
    }

    // Set WasUsed to false for the current node
    node.WasUsed = false

    // Recursively call this function for all children of the current node
    for i := range node.Children {
        resetDialogueNodes(node.Children[i])
    }
}

func (ent *DialogueManagerEntity) beginDialogueMode() {
    // TODO: pick the next node to display
    nextNode := selectNextNode(
        ent.dataHandler.GetRootNode(),
        *ent.gameState.PlayerEmotionalProfile)
    logging.Info("Next selected dialogue Node: %v", nextNode.Text)
    ent.gameState.SetCurrentPlayState(PlayStateDialog)
    ent.dialoguePanel.setCurrentDialogueNode(nextNode)
}

func (ent *DialogueManagerEntity) Deinit() {
	// Required de-initialization
	ent.BaseEntity.Deinit()

	// De-initialization logic for the entity
	// ...
}

func (ent *DialogueManagerEntity) Update() {
	// Required update
	ent.BaseEntity.Update()

	// Update logic for the entity
	// ...
}

func (ent *DialogueManagerEntity) Draw() {
	// Draw logic for the entity
	// ...
}
