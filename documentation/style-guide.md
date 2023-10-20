<!-- LTeX: language=en-US -->
# Style Guide

## Use method chaining for configuration heavy instantiation
You should attempt to utilize method chaining in situations where a lot of
configuration parameters are required to instantiate some structure.

As an example, this is the wrong way to design an instantiation interface for colliders:
```go
// Lots of arguments, one can't see what each of them means at a glance,
// some might not even be required, and a "default" would be better suited here.
NewConvexColliderAbs(
    poly,
    10, 1,
    category, callbacks,
    true, BodyTypeStatic),
)
```

This is the right way to do it:
```go
// We know what every argument means, the numbers (density and damping)
// disappeared, and instead NewConvexColliderAbs() just sets default values for
// everything, that can then be overridden.
NewConvexColliderAbs(polygon, BodyTypeStatic).
    SetCategory(category).
    SetCallbacks(callbacks).
    SetFixedRotation(true))
```


