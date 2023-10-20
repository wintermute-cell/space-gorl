<!-- LTeX: language=en-US -->
# GUI Style Properties

GUI Style properties are defined as a so called *styledef*. A styledef is a
string like `property:value|other-prop:other-value|...`.

A more concrete example would be `color:255,0,0,255|font:alagard`.

Every GUI widget must have the ability to take in such a styledef during
creation, and store it. The styledef is then read by the drawing backend of the
GUI package, and interpreted as chosen by the drawing backend implementation.

## Using styledefs
To find out more on existing styledefs and their exact specification, see the
[Styledef Specification](/documentation/gui-styledef-specification.md)

## Creating a new styledef property 
The following will show you how to create a new styledef property, such as `color`.

1. (Optional but recommended) First, write down your new property in the
   [Styledef Specification](/documentation/gui-styledef-specification.md)
2. Create a new conversion function for your property name in
   [`/pkg/gui/styledef.go:parseStyleDef()`](/pkg/gui/styledef.go) (or reuse an
   existing one and hook it up to your new property name)
3. Start using the property!

## Motivation
Why have these styledefs? Respecting a styledef is optional for the drawing
backend. Because of this, a prototype GUI might be created using a
simple/generic drawing backend, using properties that the final backend might
not require. Then, after prototyping is done, the drawing backend can be
replaced with the production implementation without having to worry about the
styles. 

Lets look at the following example: For our prototyping gui, we want to be able
to change the font for every individual Label. We can then use the following
line to define a Label:
```go
// specifying a font for the prototyping backend
label := gui.NewLabel("some text", rl.NewVector2(10, 10), "font:alagard")
```

And in our drawing backend, read out the font name, and apply the correct font
before drawing.

Later on, when prototyping is done, we want to use a single font for every
Widget; thus, our drawing backend will not read the `font` property anymore.
But we do not have to change the `font:alagard` styledef we pass into
`gui.NewLabel` but can rather just ignore it.

The same concept also applies backwards. Maybe our current drawing backend
doesn't support some properties that we want to use later, but we can pass them
anyway.

## Why a string based solution?
Specifying styles as a string like `color:255,0,0,255|font:alagard` instead of
a struct for example
```go
type StyleDef struct {
    color raylib.Color
    font raylib.Font
}
```
has certain drawbacks. We lose out on direct type safety, and have to do an
extra parsing step, untangling the string representation back into actual typed
values.

But this approach has one key benefit: *Conciseness*.

The string approach allows us to define a styledef in just a single short line.
And that is of high importance for code readability when one has to define a
huge number of widgets in one place.

Creating an instance of a hypothetical `StyleDef` struct is also possible in
one line like `StyleDef{ <properties go here> }`. But structs do not
differentiate between unset fields, and zero-fields. (There is no discernible
difference between a value set to 0 and one that was never specifically set).
