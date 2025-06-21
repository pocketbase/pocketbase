package security

import (
	cryptoRand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp/syntax"
	"strings"
)

const defaultMaxRepeat = 6

var anyCharNotNLPairs = []rune{'A', 'Z', 'a', 'z', '0', '9'}

// RandomStringByRegex generates a random string matching the regex pattern.
// If optFlags is not set, fallbacks to [syntax.Perl].
//
// NB! While the source of the randomness comes from [crypto/rand] this method
// is not recommended to be used on its own in critical secure contexts because
// the generated length could vary too much on the used pattern and may not be
// as secure as simply calling [security.RandomString].
// If you still insist on using it for such purposes, consider at least
// a large enough minimum length for the generated string, e.g. `[a-z0-9]{30}`.
//
// This function is inspired by github.com/pipe01/revregexp, github.com/lucasjones/reggen and other similar packages.
func RandomStringByRegex(pattern string, optFlags ...syntax.Flags) (string, error) {
	var flags syntax.Flags
	if len(optFlags) == 0 {
		flags = syntax.Perl
	} else {
		for _, f := range optFlags {
			flags |= f
		}
	}

	r, err := syntax.Parse(pattern, flags)
	if err != nil {
		return "", err
	}

	var sb = new(strings.Builder)

	err = writeRandomStringByRegex(r, sb)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

func writeRandomStringByRegex(r *syntax.Regexp, sb *strings.Builder) error {
	// https://pkg.go.dev/regexp/syntax#Op
	switch r.Op {
	case syntax.OpCharClass:
		c, err := randomRuneFromPairs(r.Rune)
		if err != nil {
			return err
		}
		_, err = sb.WriteRune(c)
		return err
	case syntax.OpAnyChar, syntax.OpAnyCharNotNL:
		c, err := randomRuneFromPairs(anyCharNotNLPairs)
		if err != nil {
			return err
		}
		_, err = sb.WriteRune(c)
		return err
	case syntax.OpAlternate:
		idx, err := randomNumber(len(r.Sub))
		if err != nil {
			return err
		}
		return writeRandomStringByRegex(r.Sub[idx], sb)
	case syntax.OpConcat:
		var err error
		for _, sub := range r.Sub {
			err = writeRandomStringByRegex(sub, sb)
			if err != nil {
				break
			}
		}
		return err
	case syntax.OpRepeat:
		return repeatRandomStringByRegex(r.Sub[0], sb, r.Min, r.Max)
	case syntax.OpQuest:
		return repeatRandomStringByRegex(r.Sub[0], sb, 0, 1)
	case syntax.OpPlus:
		return repeatRandomStringByRegex(r.Sub[0], sb, 1, -1)
	case syntax.OpStar:
		return repeatRandomStringByRegex(r.Sub[0], sb, 0, -1)
	case syntax.OpCapture:
		return writeRandomStringByRegex(r.Sub[0], sb)
	case syntax.OpLiteral:
		_, err := sb.WriteString(string(r.Rune))
		return err
	default:
		return fmt.Errorf("unsupported pattern operator %d", r.Op)
	}
}

func repeatRandomStringByRegex(r *syntax.Regexp, sb *strings.Builder, min int, max int) error {
	if max < 0 {
		max = defaultMaxRepeat
	}

	if max < min {
		max = min
	}

	n := min
	if max != min {
		randRange, err := randomNumber(max - min)
		if err != nil {
			return err
		}
		n += randRange
	}

	var err error
	for i := 0; i < n; i++ {
		err = writeRandomStringByRegex(r, sb)
		if err != nil {
			return err
		}
	}

	return nil
}

func randomRuneFromPairs(pairs []rune) (rune, error) {
	if len(pairs)%2 != 0 {
		return 0, fmt.Errorf("invalid pairs slice: odd number of elements")
	}

	// Pre-calculate the cumulative size of all ranges to make the selection process more efficient.
	cumulativeSizes := make([]int, len(pairs)/2)
	totalRunes := 0
	for i := 0; i < len(pairs); i += 2 {
		start, end := pairs[i], pairs[i+1]
		if start > end {
			return 0, fmt.Errorf("invalid range: start '%c' > end '%c'", start, end)
		}
		totalRunes += int(end - start + 1)
		cumulativeSizes[i/2] = totalRunes
	}

	if totalRunes == 0 {
		return 0, errors.New("no runes to choose from")
	}

	// Select a random number in the range of total runes.
	runeNumber, err := randomNumber(totalRunes)
	if err != nil {
		return 0, fmt.Errorf("failed to generate random number: %w", err)
	}

	// Find which range the selected number falls into using the pre-calculated cumulative sizes.
	for i, size := range cumulativeSizes {
		if runeNumber < size {
			startRune := pairs[i*2]
			previousSize := 0
			if i > 0 {
				previousSize = cumulativeSizes[i-1]
			}
			return startRune + rune(runeNumber-previousSize), nil
		}
	}

	// This part should be unreachable if the logic is correct.
	// It indicates a bug in this function or in randomNumber.
	panic("unreachable: failed to find a rune")
}

func randomNumber(maxSoft int) (int, error) {
	randRange, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(int64(maxSoft)))

	return int(randRange.Int64()), err
}
