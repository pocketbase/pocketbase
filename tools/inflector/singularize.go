package inflector

import (
	"log"
	"regexp"

	"github.com/pocketbase/pocketbase/tools/store"
)

var compiledPatterns = store.New[string, *regexp.Regexp](nil)

// note: the patterns are extracted from popular Ruby/PHP/Node.js inflector packages
var singularRules = []struct {
	pattern     string // lazily compiled
	replacement string
}{
	{"(?i)([nrlm]ese|deer|fish|sheep|measles|ois|pox|media|ss)$", "${1}"},
	{"(?i)^(sea[- ]bass)$", "${1}"},
	{"(?i)(s)tatuses$", "${1}tatus"},
	{"(?i)(f)eet$", "${1}oot"},
	{"(?i)(t)eeth$", "${1}ooth"},
	{"(?i)^(.*)(menu)s$", "${1}${2}"},
	{"(?i)(quiz)zes$", "${1}"},
	{"(?i)(matr)ices$", "${1}ix"},
	{"(?i)(vert|ind)ices$", "${1}ex"},
	{"(?i)^(ox)en", "${1}"},
	{"(?i)(alias)es$", "${1}"},
	{"(?i)(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$", "${1}us"},
	{"(?i)([ftw]ax)es", "${1}"},
	{"(?i)(cris|ax|test)es$", "${1}is"},
	{"(?i)(shoe)s$", "${1}"},
	{"(?i)(o)es$", "${1}"},
	{"(?i)ouses$", "ouse"},
	{"(?i)([^a])uses$", "${1}us"},
	{"(?i)([m|l])ice$", "${1}ouse"},
	{"(?i)(x|ch|ss|sh)es$", "${1}"},
	{"(?i)(m)ovies$", "${1}ovie"},
	{"(?i)(s)eries$", "${1}eries"},
	{"(?i)([^aeiouy]|qu)ies$", "${1}y"},
	{"(?i)([lr])ves$", "${1}f"},
	{"(?i)(tive)s$", "${1}"},
	{"(?i)(hive)s$", "${1}"},
	{"(?i)(drive)s$", "${1}"},
	{"(?i)([^fo])ves$", "${1}fe"},
	{"(?i)(^analy)ses$", "${1}sis"},
	{"(?i)(analy|diagno|^ba|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$", "${1}${2}sis"},
	{"(?i)([ti])a$", "${1}um"},
	{"(?i)(p)eople$", "${1}erson"},
	{"(?i)(m)en$", "${1}an"},
	{"(?i)(c)hildren$", "${1}hild"},
	{"(?i)(n)ews$", "${1}ews"},
	{"(?i)(n)etherlands$", "${1}etherlands"},
	{"(?i)eaus$", "eau"},
	{"(?i)(currenc)ies$", "${1}y"},
	{"(?i)^(.*us)$", "${1}"},
	{"(?i)s$", ""},
}

// Singularize converts the specified word into its singular version.
//
// For example:
//
//	inflector.Singularize("people") // "person"
func Singularize(word string) string {
	if word == "" {
		return ""
	}

	for _, rule := range singularRules {
		re := compiledPatterns.GetOrSet(rule.pattern, func() *regexp.Regexp {
			re, err := regexp.Compile(rule.pattern)
			if err != nil {
				return nil
			}
			return re
		})
		if re == nil {
			// log only for debug purposes
			log.Println("[Singularize] failed to retrieve/compile rule pattern " + rule.pattern)
			continue
		}

		if re.MatchString(word) {
			return re.ReplaceAllString(word, rule.replacement)
		}
	}

	return word
}
