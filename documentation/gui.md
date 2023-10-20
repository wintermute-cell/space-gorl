<!-- LTeX: language=en-US -->
# GUI

## Abstract Structure
The `gui` package is split into two logical parts: The logic processing
(registering button clicks, moving the scroll panel, etc.) in `pkg/gui/gui.go`
and the rendering (drawing colored buttons, rendering text in different fonts,
providing styling resources, etc.) in `pkg/gui/backend.go`.

This means that `gui` code flows like this:
```
Define gui Widget in scene (pass position, size, etc. and styling info)
->
Gui widgets functionality is processed (reacting to user input) and the state
of the Widget is updated.
->
The Widget, along with its current state is passed on to the backend, where it
is drawn. The style info provided in the first step will also be unpacked here,
and can be interpreted by the backend.
```

Let's look at the `Button` Widget for example. This is an excerpt of the
`doRecursiveDraw()` function:
```go
case *Button:
    w.update_button() // <- logical update
    backend_button(*w) // <- backend render step
    // a component with children might draw its children here,
    // to be able to draw below and above its children.
    backend_button_finalize(*w) // <- another backend render step
```
(Here, `w` is a pointer to a `Button` widget.)

> [!NOTE]
> All styling resources, such as fonts, icons, etc. are to be provided by
> `/pkg/gui/backend.go`

## Usage

### Available Widgets

#### BaseWidget
This is not a real Widget. It provides core properties, and is inherited by
every other Widget listed below.

**Properties**:
- `position`: A Vector2 storing the position of the Widget.
- `size`: A Vector2 storing the size of the Widget.

**Methods**:
- `SetPosition(p rl.Vector2)`
- `GetPosition() rl.Vector2`
- `SetSize(s rl.Vector2)`
- `GetSize() rl.Vector2`
- `Bounds() rl.Rectangle`

#### Label

**Definition**:
A widget used to display a piece of text.

**Usage**:
```go
label := NewLabel("Sample Text", rl.Vector2{X: 10, Y: 20}, "color:255,255,255")
```

#### Button

**Definition**:
A widget that represents an interactive button with the capability to detect hover and click states.

**Usage**:
```go
btn := NewButton("Click Me", rl.Vector2{X: 30, Y: 40}, rl.Vector2{X: 100, Y: 30}, myCallbackFunction, "color:255,0,0|bgColor:0,0,255")
```

#### Scroll Panel

**Definition**:
A scrollable container that holds and arranges other widgets inside it.

**Methods**:
- `AddChild(child Widget)`: Adds a child widget to the ScrollPanel.
- `RemoveChild(target Widget)`: Removes a specified child widget from the ScrollPanel.

**Usage**:
```go
scrollPanel := NewScrollPanel(rl.NewRectangle(10, 10, 200, 200), rl.NewRectangle(0, 0, 400, 400), "color:150,150,150")
scrollPanel.AddChild(myLabelWidget)
scrollPanel.RemoveChild(myLabelWidget)
```

## Internal Implementation
The following will describe how the `gui` package is structured and implemented.
Afterward a guide on how to implement new Widgets will follow.

The `gui` package has one main object, the `Gui`. This `Gui` object is a
container for other objects, and has a `Draw()` function. This draw function
should be called every frame while the GUI is active.

Beyond that, every widget is either atomic (such as a `Button`) or another
container (such as a `ScrollPanel`). Using container Widgets (and their `AddChild` functions), a kind of tree structure is created.

Example:
```
Gui
|- Button
|
|- Button
|
|- ScrollPanel
|   |- Label
|   |- Button
|
|- ScrollPanel
|   |- Label
|   |- Button
|
|- Label
|
```

### Walking the tree
To ensure every Widget is drawn in the correct order, the Gui-Objects `Draw()`
function actually starts a recursive drawing process, beginning at the `Gui`
node.

Here is an excerpt from that recursive function:
```go
func doRecursiveDraw(container Container) {
	for _, widget := range container.Children {
		switch w := any(widget).(type) {
		case *Label:
			w.update_label()
			backend_label(*w)
			backend_label_finalize(*w)
		case *ScrollPanel:
			w.update_scroll_panel()
			backend_scroll_panel(*w)
			doRecursiveDraw(*w.container) // draw the panels children
			backend_scroll_panel_finalize(*w)
    }
}
```

As one may see, every container Widget type (one that has its own children),
must call `doRecursiveDraw()` on its own children.

### How to create a new Widget
1. Define the widget in `pkg/gui/gui.go`:
```go
// ----------------
//    MY WIDGET   |
// ----------------

// widget definiton
type MyWidget struct {
	BaseWidget // Every widget must inherit BaseWidget
	style_info string // Every widget must have a style_info field

    // here you can add more fields you Widget might need,
    // such as a callback function for a button, or general state information
    // for example.
}

// update function
func (my_widget *MyWidget) update_my_widget() {
    // here you can perform general Widget functionality, and modify the state
    // of the Widget accordingly. For example checking for a mouse click in
    // a button, and if one is registered, calling the callback and setting the
    // buttons state to pressed, so it can be rendered in a "pressed" color
    // later.
}

// constructor
func NewMyWidget(/* taking in all required information for MyWidget */) *MyWidget {
	return &MyWidget{
        // Fill out the required fields
    }
}

// Now may come any additional functions the Widget might require, such as an
// AddChild() function for a container widget for example.
```

2. After that, implement a rendering function in `pkg/gui/backend.go`:
```go 
// not taking in a pointer, this step will only draw, not affect state!
func backend_my_widget(my_widget MyWidget) {
    style := parseStyleDef(label.style_info) // first, unpack the styling information if you need it

    // an example for using the provided styling info
    color := rl.Black // we choose a default fallback value
    if c, ok := style["color"]; ok && c != nil { // we check if the field is present in the style map
        color = c.(rl.Color) // and we overwrite the fallback with the given style
    }

    // Now we draw the Widget according to our needs (optionally respecting
    // some or all style properties)
    rl.DrawSomething(...)
}

// we may also have another backend function, in case we need to draw to steps.
func backend_my_widget_finalize(my_widget MyWidget) {
	// implement more drawing logic
}
```

3. The last step is to add the new Widget to the `doRecursiveDraw()` function
in `pkg/gui/gui.go`.
```go
func doRecursiveDraw(container Container) {
	for _, widget := range container.Children {
		switch w := any(widget).(type) {
		case *Label:
			w.update_label()
			backend_label(*w)
			backend_label_finalize(*w)
        /* ... other widgets ... */
		case *ScrollPanel:
			w.update_scroll_panel()
			backend_scroll_panel(*w)
			doRecursiveDraw(*w.container)
			backend_scroll_panel_finalize(*w)
        /* ... other widgets ... */

        case *MyWidget:                 // <==== our new widget!
            w.update_my_widget() // first, running the widget functionality
			backend_my_widget(*w) // then, draw the widget
			doRecursiveDraw(*w.container) // optionally, draw children
			backend_my_widget_finalize(*w) // optionally, do another draw step

        /* ... other widgets ... */
```



### Swapping out the backend
The entire implementation of `pkg/gui/backend.go` can be swapped out, following
this contract:

1. The backend must provide all required resources, such as fonts, icons, etc.
2. The backend may or may not respect styling information passed by the widgets.
(read the [Styledef Docs](/documentation/gui-styledef.md).)
3. The backend must contain a `backend_<widget-name>(w Widget)` function, and a 
`backend_<widget-name>_finalize(w Widget)` function for **every** Widget!
4. That's it.
