package entities

import (
	"cowboy-gorl/pkg/entities/gem"
	"cowboy-gorl/pkg/entities/proto"
	"cowboy-gorl/pkg/physics"
	"cowboy-gorl/pkg/util"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EnemyType int32
const (
    EnemyTypeBasic EnemyType = iota
    // ... other types
)

type EnemyData struct {
    Type          EnemyType
    Sprite          rl.Texture2D
    Size rl.Vector2
    Hitpoints   float32
    ScoreWorth int32
    Damage        int32
    MovementPattern func(collider *physics.Collider) // Function to determine movement
    LinearDamping float32
    AngularDamping float32
    // ... other properties
}

// WaveManager Entity
type WaveManagerEntity struct {
	// Required fields
	proto.BaseEntity

	// Custom Fields
    attackTarget rl.Vector2

    enemyPool    map[EnemyType][]*EnemyEntity2D
    activeEnemies []*EnemyEntity2D
    enemyData    map[EnemyType]EnemyData

    spawnTimer *util.Timer
}

func (ent *WaveManagerEntity) SpawnEnemy(enemyType EnemyType) {
    // Reuse or create a new enemy based on type
    var enemy *EnemyEntity2D
    spawnPosition := rl.NewVector2(float32(rand.Int31n(320)), float32(rand.Int31n(20)))
    if len(ent.enemyPool[enemyType]) > 0 {
        enemy = ent.enemyPool[enemyType][len(ent.enemyPool[enemyType])-1]
        enemy.SetPosition(spawnPosition)
        ent.enemyPool[enemyType] = ent.enemyPool[enemyType][:len(ent.enemyPool[enemyType])-1]
    } else {
        // TODO: what values to initialize with?
        enemy = NewEnemyEntity2D(spawnPosition, 0, rl.Vector2One()) // Function to create a new enemy
    }

    // Initialize enemy based on EnemyData
    data := ent.enemyData[enemyType]
    enemy.SetData(data)
    gem.AddEntity(ent, enemy)

    ent.activeEnemies = append(ent.activeEnemies, enemy)
}


func NewWaveManagerEntity() *WaveManagerEntity {
	new_ent := &WaveManagerEntity{
        attackTarget: rl.NewVector2(160, 160),
        enemyPool: make(map[EnemyType][]*EnemyEntity2D),
        enemyData: make(map[EnemyType]EnemyData),
        spawnTimer: util.NewTimer(0.01),
    }

    new_ent.enemyData[EnemyTypeBasic] = EnemyData{
        Type: EnemyTypeBasic,
        Sprite: rl.LoadTexture("sprites/ships/basic.png"),
        Size: rl.NewVector2(4, 8),
        Hitpoints: 3,
        ScoreWorth: 5,
        Damage: 10,
        MovementPattern: func(collider *physics.Collider) {
            // Calculate separation, alignment, and cohesion forces
            separationForce := CalculateSeparationForce(collider, new_ent.activeEnemies, 10)
            alignmentForce := CalculateAlignmentForce(new_ent.activeEnemies, 10)
            cohesionForce := CalculateCohesionForce(collider, new_ent.activeEnemies, 10)

            ApplyTorqueToRotate(collider, new_ent.attackTarget, 0.01, alignmentForce, cohesionForce)
            ApplyForceForAcceleration(collider, 10, separationForce)
        },
        LinearDamping: 5,
        AngularDamping: 50,
    }

	return new_ent
}

func (ent *WaveManagerEntity) Init() {
	// Required initialization
	ent.BaseEntity.Init()

	// Initialization logic for the entity
	// ...
}

func (ent *WaveManagerEntity) Deinit() {
	// Required de-initialization
	ent.BaseEntity.Deinit()

	// De-initialization logic for the entity
	// ...
}
func (ent *WaveManagerEntity) Update() {
    // ... Update logic, including spawning enemies

    if ent.spawnTimer.Check() {
        ent.SpawnEnemy(EnemyTypeBasic)
    }

    toRemove := []*EnemyEntity2D{}
    for _, e := range ent.activeEnemies {
        if !e.isAlive {
            toRemove = append(toRemove, e)
            ent.RecycleEnemy(e)
        }
    }

    for _, e := range toRemove {
        idx := util.SliceIndex(ent.activeEnemies, e)
        ent.activeEnemies = util.SliceDelete(ent.activeEnemies, idx, idx + 1)
    }
}

func (ent *WaveManagerEntity) Draw() {
}

func (ent *WaveManagerEntity) RecycleEnemy(enemy *EnemyEntity2D) {
    // Return the enemy to the pool
    gem.RemoveEntity(enemy)
    ent.enemyPool[enemy.GetData().Type] = append(ent.enemyPool[enemy.GetData().Type], enemy)
}

