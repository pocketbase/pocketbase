package inflector_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/tools/inflector"
)

func TestSingularize(t *testing.T) {
	scenarios := []struct {
		word     string
		expected string
	}{
		{"abcnese", "abcnese"},
		{"deer", "deer"},
		{"sheep", "sheep"},
		{"measles", "measles"},
		{"pox", "pox"},
		{"media", "media"},
		{"bliss", "bliss"},
		{"sea-bass", "sea-bass"},
		{"Statuses", "Status"},
		{"Feet", "Foot"},
		{"Teeth", "Tooth"},
		{"abcmenus", "abcmenu"},
		{"Quizzes", "Quiz"},
		{"Matrices", "Matrix"},
		{"Vertices", "Vertex"},
		{"Indices", "Index"},
		{"Aliases", "Alias"},
		{"Alumni", "Alumnus"},
		{"Bacilli", "Bacillus"},
		{"Cacti", "Cactus"},
		{"Fungi", "Fungus"},
		{"Nuclei", "Nucleus"},
		{"Radii", "Radius"},
		{"Stimuli", "Stimulus"},
		{"Syllabi", "Syllabus"},
		{"Termini", "Terminus"},
		{"Viri", "Virus"},
		{"Faxes", "Fax"},
		{"Crises", "Crisis"},
		{"Axes", "Axis"},
		{"Shoes", "Shoe"},
		{"abcoes", "abco"},
		{"Houses", "House"},
		{"Mice", "Mouse"},
		{"abcxes", "abcx"},
		{"Movies", "Movie"},
		{"Series", "Series"},
		{"abcquies", "abcquy"},
		{"Relatives", "Relative"},
		{"Drives", "Drive"},
		{"aardwolves", "aardwolf"},
		{"Analyses", "Analysis"},
		{"Diagnoses", "Diagnosis"},
		{"People", "Person"},
		{"Men", "Man"},
		{"Children", "Child"},
		{"News", "News"},
		{"Netherlands", "Netherlands"},
		{"Tableaus", "Tableau"},
		{"Currencies", "Currency"},
		{"abcs", "abc"},
		{"abc", "abc"},
	}

	for _, s := range scenarios {
		t.Run(s.word, func(t *testing.T) {
			result := inflector.Singularize(s.word)
			if result != s.expected {
				t.Fatalf("Expected %q, got %q", s.expected, result)
			}
		})
	}
}
