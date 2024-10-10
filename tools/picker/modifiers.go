package picker

import (
	"fmt"

	"github.com/pocketbase/pocketbase/tools/tokenizer"
)

var Modifiers = map[string]ModifierFactoryFunc{}

type ModifierFactoryFunc func(args ...string) (Modifier, error)

type Modifier interface {
	// Modify executes the modifier and returns a new modified value.
	Modify(value any) (any, error)
}

func initModifer(rawModifier string) (Modifier, error) {
	t := tokenizer.NewFromString(rawModifier)
	t.Separators('(', ')', ',', ' ')
	t.IgnoreParenthesis(true)

	parts, err := t.ScanAll()
	if err != nil {
		return nil, err
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid or empty modifier expression %q", rawModifier)
	}

	name := parts[0]
	args := parts[1:]

	factory, ok := Modifiers[name]
	if !ok {
		return nil, fmt.Errorf("missing or invalid modifier %q", name)
	}

	return factory(args...)
}
