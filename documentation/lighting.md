<!-- LTeX: language=en-US -->
# Lighting

Lighting is split up in two parts. **Shadow Casting** and **Sprite Lighting**.

## Usage

If a scene requires lighting, it must call `lighting.Enable()` during its
startup and `lighting.Disable()` during its unloading.

### Creating and Deletion of Lights and Occluders
Lights and occluders (shadow casters) must be managed by the creator. To create
a light, invoke:
```go
light := lighting.NewLight2D(
    rl.NewVector2(512, 512), // size
    1024, // resolution
    rl.NewColor(255, 90, 130, 255), // color
    )
```

And remember to later unload the light if you don't need it anymore.
```go
lighting.UnloadLight(light)
```
This way, you can load a scene that creates a light, unload the scene and load
the scene again without there being duplicate lights. If you forget to unload
the light/occluder, it will persist until you unload it.

You create and unload occluders similarly:
```go
occluder = lighting.NewOccluderSprite2D(
    my_sprite,
	my_position,
	my_size,
	my_origin,
	my_rotation,
    )

// This creates a new occluder with zero'd out values and automatic size.
// Useful if the occluder is dynamic and will be updated every frame anyway.
other_occluder = lighting.NewOccluderSprite2DZ(
    my_sprite,
    )
```

And later:
```go
lighting.UnloadOccluder(occluder)
lighting.UnloadOccluder(other_occluder)
```

### Updating Occluders
If an occluder is dynamic, meaning it changes its properties like position,
rotation, etc., then you must supply these new values every time they change,
using `occluder.Update()`.
