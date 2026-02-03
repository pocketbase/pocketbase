package core

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/spf13/cast"
)

func init() {
	Fields[FieldTypeColor] = func() Field {
		return &ColorField{}
	}
}

const FieldTypeColor = "color"

var (
	_ Field        = (*ColorField)(nil)
	_ SetterFinder = (*ColorField)(nil)
)

var colorPatterns = map[string]*regexp.Regexp{
	"hex":   regexp.MustCompile(`^#([A-Fa-f0-9]{3}|[A-Fa-f0-9]{4}|[A-Fa-f0-9]{6}|[A-Fa-f0-9]{8})$`),
	"rgb":   regexp.MustCompile(`^rgba?\s*\(\s*(\d{1,3})\s*[,\s]\s*(\d{1,3})\s*[,\s]\s*(\d{1,3})\s*(?:[,\/]\s*([0-9.]+%?)\s*)?\)$`),
	"hsl":   regexp.MustCompile(`^hsla?\s*\(\s*(\d{1,3}(?:deg|grad|rad|turn)?)\s*[,\s]\s*(\d{1,3})%\s*[,\s]\s*(\d{1,3})%\s*(?:[,\/]\s*([0-9.]+%?)\s*)?\)$`),
	"hwb":   regexp.MustCompile(`^hwb\s*\(\s*(\d{1,3}(?:deg|grad|rad|turn)?)\s+(\d{1,3})%\s+(\d{1,3})%\s*(?:\/\s*([0-9.]+%?)\s*)?\)$`),
	"lab":   regexp.MustCompile(`^lab\s*\(\s*([0-9.]+%?)\s+(-?[0-9.]+)\s+(-?[0-9.]+)\s*(?:\/\s*([0-9.]+%?)\s*)?\)$`),
	"lch":   regexp.MustCompile(`^lch\s*\(\s*([0-9.]+%?)\s+([0-9.]+)\s+(\d{1,3}(?:deg|grad|rad|turn)?)\s*(?:\/\s*([0-9.]+%?)\s*)?\)$`),
	"oklab": regexp.MustCompile(`^oklab\s*\(\s*([0-9.]+%?)\s+(-?[0-9.]+)\s+(-?[0-9.]+)\s*(?:\/\s*([0-9.]+%?)\s*)?\)$`),
	"oklch": regexp.MustCompile(`^oklch\s*\(\s*([0-9.]+%?)\s+([0-9.]+)\s+(\d{1,3}(?:deg|grad|rad|turn)?)\s*(?:\/\s*([0-9.]+%?)\s*)?\)$`),
}

var cssNamedColors = map[string]bool{
	"aliceblue": true, "antiquewhite": true, "aqua": true, "aquamarine": true,
	"azure": true, "beige": true, "bisque": true, "black": true, "blanchedalmond": true,
	"blue": true, "blueviolet": true, "brown": true, "burlywood": true, "cadetblue": true,
	"chartreuse": true, "chocolate": true, "coral": true, "cornflowerblue": true,
	"cornsilk": true, "crimson": true, "cyan": true, "darkblue": true, "darkcyan": true,
	"darkgoldenrod": true, "darkgray": true, "darkgrey": true, "darkgreen": true,
	"darkkhaki": true, "darkmagenta": true, "darkolivegreen": true, "darkorange": true,
	"darkorchid": true, "darkred": true, "darksalmon": true, "darkseagreen": true,
	"darkslateblue": true, "darkslategray": true, "darkslategrey": true, "darkturquoise": true,
	"darkviolet": true, "deeppink": true, "deepskyblue": true, "dimgray": true, "dimgrey": true,
	"dodgerblue": true, "firebrick": true, "floralwhite": true, "forestgreen": true,
	"fuchsia": true, "gainsboro": true, "ghostwhite": true, "gold": true, "goldenrod": true,
	"gray": true, "grey": true, "green": true, "greenyellow": true, "honeydew": true,
	"hotpink": true, "indianred": true, "indigo": true, "ivory": true, "khaki": true,
	"lavender": true, "lavenderblush": true, "lawngreen": true, "lemonchiffon": true,
	"lightblue": true, "lightcoral": true, "lightcyan": true, "lightgoldenrodyellow": true,
	"lightgray": true, "lightgrey": true, "lightgreen": true, "lightpink": true,
	"lightsalmon": true, "lightseagreen": true, "lightskyblue": true, "lightslategray": true,
	"lightslategrey": true, "lightsteelblue": true, "lightyellow": true, "lime": true,
	"limegreen": true, "linen": true, "magenta": true, "maroon": true, "mediumaquamarine": true,
	"mediumblue": true, "mediumorchid": true, "mediumpurple": true, "mediumseagreen": true,
	"mediumslateblue": true, "mediumspringgreen": true, "mediumturquoise": true,
	"mediumvioletred": true, "midnightblue": true, "mintcream": true, "mistyrose": true,
	"moccasin": true, "navajowhite": true, "navy": true, "oldlace": true, "olive": true,
	"olivedrab": true, "orange": true, "orangered": true, "orchid": true, "palegoldenrod": true,
	"palegreen": true, "paleturquoise": true, "palevioletred": true, "papayawhip": true,
	"peachpuff": true, "peru": true, "pink": true, "plum": true, "powderblue": true,
	"purple": true, "rebeccapurple": true, "red": true, "rosybrown": true, "royalblue": true,
	"saddlebrown": true, "salmon": true, "sandybrown": true, "seagreen": true, "seashell": true,
	"sienna": true, "silver": true, "skyblue": true, "slateblue": true, "slategray": true,
	"slategrey": true, "snow": true, "springgreen": true, "steelblue": true, "tan": true,
	"teal": true, "thistle": true, "tomato": true, "turquoise": true, "violet": true,
	"wheat": true, "white": true, "whitesmoke": true, "yellow": true, "yellowgreen": true,
	"transparent": true,
}

// ColorField defines "color" type field for storing color values.
// Supports:
// - Hex: #RGB, #RRGGBB, #RRGGBBAA
// - RGB/RGBA: rgb(255 0 0), rgba(255 0 0 / 0.5)
// - HSL/HSLA: hsl(120deg 100% 50%), hsla(120 100% 50% / 0.5)
// - HWB: hwb(120deg 30% 50%)
// - LAB: lab(50% 40 30)
// - LCH: lch(50% 40 120deg)
// - OKLAB: oklab(0.5 0.1 0.1)
// - OKLCH: oklch(0.5 0.1 120deg)
// - CSS named colors: red, blue, rebeccapurple, etc.
//
// The respective zero record field value is empty string.
type ColorField struct {
	Name string `form:"name" json:"name"`
	Id string `form:"id" json:"id"`
	System bool `form:"system" json:"system"`
	Hidden bool `form:"hidden" json:"hidden"`
	Presentable bool `form:"presentable" json:"presentable"`
	Required bool `form:"required" json:"required"`
}

// Type implements [Field.Type] interface method.
func (f *ColorField) Type() string {
	return FieldTypeColor
}

// GetId implements [Field.GetId] interface method.
func (f *ColorField) GetId() string {
	return f.Id
}

// SetId implements [Field.SetId] interface method.
func (f *ColorField) SetId(id string) {
	f.Id = id
}

// GetName implements [Field.GetName] interface method.
func (f *ColorField) GetName() string {
	return f.Name
}

// SetName implements [Field.SetName] interface method.
func (f *ColorField) SetName(name string) {
	f.Name = name
}

// GetSystem implements [Field.GetSystem] interface method.
func (f *ColorField) GetSystem() bool {
	return f.System
}

// SetSystem implements [Field.SetSystem] interface method.
func (f *ColorField) SetSystem(system bool) {
	f.System = system
}

// GetHidden implements [Field.GetHidden] interface method.
func (f *ColorField) GetHidden() bool {
	return f.Hidden
}

// SetHidden implements [Field.SetHidden] interface method.
func (f *ColorField) SetHidden(hidden bool) {
	f.Hidden = hidden
}

// ColumnType implements [Field.ColumnType] interface method.
func (f *ColorField) ColumnType(app App) string {
	return "TEXT DEFAULT '' NOT NULL"
}

// PrepareValue implements [Field.PrepareValue] interface method.
func (f *ColorField) PrepareValue(record *Record, raw any) (any, error) {
	value := cast.ToString(raw)
	value = strings.TrimSpace(value)
	return value, nil
}

// ValidateValue implements [Field.ValidateValue] interface method.
func (f *ColorField) ValidateValue(ctx context.Context, app App, record *Record) error {
	value, ok := record.GetRaw(f.Name).(string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	return f.ValidatePlainValue(value)
}

// ValidatePlainValue validates color format.
func (f *ColorField) ValidatePlainValue(value string) error {
	if f.Required {
		if err := validation.Required.Validate(value); err != nil {
			return err
		}
	}

	if value == "" {
		return nil
	}

	value = strings.TrimSpace(strings.ToLower(value))

	// Check hex format
	if strings.HasPrefix(value, "#") {
		if colorPatterns["hex"].MatchString(value) {
			return nil
		}
		return validation.NewError("validation_invalid_hex_color", "Invalid hex color format.")
	}

	// Check modern CSS color formats
	colorTypes := []struct {
		prefix  string
		pattern string
		name    string
	}{
		{"oklch(", "oklch", "OKLCH"},
		{"oklab(", "oklab", "OKLAB"},
		{"lch(", "lch", "LCH"},
		{"lab(", "lab", "LAB"},
		{"hwb(", "hwb", "HWB"},
		{"hsla(", "hsl", "HSL"},
		{"hsl(", "hsl", "HSL"},
		{"rgba(", "rgb", "RGB"},
		{"rgb(", "rgb", "RGB"},
	}

	for _, ct := range colorTypes {
		if strings.HasPrefix(value, ct.prefix) {
			if colorPatterns[ct.pattern].MatchString(value) {
				// Additional validation for RGB values (0-255 range)
				if ct.pattern == "rgb" {
					matches := colorPatterns["rgb"].FindStringSubmatch(value)
					if len(matches) >= 4 {
						for i := 1; i <= 3; i++ {
							val, _ := strconv.Atoi(matches[i])
							if val > 255 {
								return validation.NewError("validation_invalid_rgb_range", "RGB values must be between 0 and 255.")
							}
						}
					}
				}
				return nil
			}
			return validation.NewError(
				"validation_invalid_"+ct.pattern+"_color",
				"Invalid "+ct.name+" color format.",
			)
		}
	}

	// Check named color
	if cssNamedColors[value] {
		return nil
	}

	return validation.NewError(
		"validation_invalid_color",
		"Invalid color format. Supported formats: hex (#RGB, #RRGGBB), rgb(a), hsl(a), hwb, lab, lch, oklab, oklch, or CSS color names.",
	)
}

// ValidateSettings implements [Field.ValidateSettings] interface method.
func (f *ColorField) ValidateSettings(ctx context.Context, app App, collection *Collection) error {
	return validation.ValidateStruct(f,
		validation.Field(&f.Id, validation.By(DefaultFieldIdValidationRule)),
		validation.Field(&f.Name, validation.By(DefaultFieldNameValidationRule)),
	)
}

// FindSetter implements the [SetterFinder] interface.
func (f *ColorField) FindSetter(key string) SetterFunc {
	if key == f.Name {
		return func(record *Record, raw any) {
			value := cast.ToString(raw)
			record.SetRaw(f.Name, strings.TrimSpace(value))
		}
	}
	return nil
}