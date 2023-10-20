<!-- LTeX: language=en-US -->
# GUI Style Properties Specification

This document provides a specification for various style properties, their expected formats, and their associated data types.

### Font Color

**Property Name:** `color`

**Format:** `R,G,B,A`

- **R, G, B, A:** Integer values from 0 to 255, indicating the red, green, blue and alpha components respectively.

**Example:** `color:255,0,0,255`

**Associated Data Type:** `raylib.Color`

---

### Font

**Property Name:** `font`

**Format:** A string representing a font name.

**Example:** `font:alagard`

**Associated Data Type:** `string`

---

### Font Scale

**Property Name:** `font-scale`

**Format:** A floating point value indicating the font scaling factor.

**Example:** `font-scale:2.0`

**Associated Data Type:** `float32`

---

### Background Color

**Property Name:** `background`

**Format:** `R,G,B,A`

- **R, G, B, A:** Integer values from 0 to 255, indicating the red, green, blue and alpha components respectively.

**Example:** `color:255,0,0,255`

**Associated Data Type:** `raylib.Color`

---

### Background Color Hover

**Property Name:** `background-hover`

**Format:** `R,G,B,A`

- **R, G, B, A:** Integer values from 0 to 255, indicating the red, green, blue and alpha components respectively.

**Example:** `color:255,0,0,255`

**Associated Data Type:** `raylib.Color`

---

### Background Color Pressed

**Property Name:** `background-pressed`

**Format:** `R,G,B,A`

- **R, G, B, A:** Integer values from 0 to 255, indicating the red, green, blue and alpha components respectively.

**Example:** `color:255,0,0,255`

**Associated Data Type:** `raylib.Color`

---

### Debugging Flag

**Property Name:** `debug`

**Format:** A boolean value, indicating whether the widget should draw debugging information.

**Example:** `debug:true`

**Associated Data Type:** `bool`


> [!NOTE]
> Developers are advised to ensure values adhere to the specified format for
> consistent rendering. Any deviation might result in unexpected behaviors or
> errors. Additional properties can be added in the future based on project
> requirements.
