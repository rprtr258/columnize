package columnize

import (
	crand "crypto/rand"
	"fmt"
	"runtime"
	"testing"
)

func TestColumnize(t *testing.T) {
	for name, tc := range map[string]struct {
		input    [][]string
		opts     []Option
		expected string
	}{
		"ListOfStringsInput": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "y", "z"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"x         y         z",
		},
		"EmptyLinesOutput": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{""},
				{"x", "y", "z"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"\n" +
				"x         y         z",
		},
		"LeadingSpacePreserved": {
			input: [][]string{
				{"", "Column B", "Column C"},
				{"x", "y", "z"},
			},
			opts: nil,
			expected: "   Column B  Column C\n" +
				"x  y         z",
		},
		"ColumnWidthCalculator": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"Longer than A", "Longer than B", "Longer than C"},
				{"short", "short", "short"},
			},
			opts: nil,
			expected: "Column A       Column B       Column C\n" +
				"Longer than A  Longer than B  Longer than C\n" +
				"short          short          short",
		},
		"ColumnWidthCalculatorNonASCII": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"⌘⌘⌘⌘⌘⌘⌘⌘", "Longer than B", "Longer than C"},
				{"short", "short", "short"},
			},
			opts: nil,
			expected: "Column A  Column B       Column C\n" +
				"⌘⌘⌘⌘⌘⌘⌘⌘  Longer than B  Longer than C\n" +
				"short     short          short",
		},
		"VariedInputSpacing": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "y", "z"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"x         y         z",
		},
		"VariedInputSpacing_NoTrim": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "y", "  z"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"x         y           z",
		},
		"UnmatchedColumnCounts": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"Value A", "Value B"},
				{"Value A", "Value B", "Value C", "Value D"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"Value A   Value B\n" +
				"Value A   Value B   Value C   Value D",
		},
		"AlternateDelimiter": {
			input: [][]string{
				{"Column | A", "Column | B", "Column | C"},
				{"Value A", "Value B", "Value C"},
			},
			opts: nil,
			expected: "Column | A  Column | B  Column | C\n" +
				"Value A     Value B     Value C",
		},
		"AlternateSpacingString": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "y", "z"},
			},
			opts: []Option{WithSeparator("    ")},
			expected: "Column A    Column B    Column C\n" +
				"x           y           z",
		},
		"SimpleFormat": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "y", "z"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"x         y         z",
		},
		"AlternatePrefixString": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "y", "z"},
			},
			opts: []Option{WithPrefix("  ")},
			expected: "  Column A  Column B  Column C\n" +
				"  x         y         z",
		},
		"EmptyFieldReplacement": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "<none>", "z"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"x         <none>    z",
		},
		"EmptyConfigValues": {
			input: [][]string{
				{"Column A", "Column B", "Column C"},
				{"x", "y", "z"},
			},
			opts: nil,
			expected: "Column A  Column B  Column C\n" +
				"x         y         z",
		},
		"No Input": {
			input:    [][]string{},
			opts:     nil,
			expected: "",
		},
		"WithHeaders": {
			input: [][]string{
				{"x", "y", "z"},
			},
			opts: []Option{WithHeaders("Column A", "Column B", "Column C")},
			expected: "Column A  Column B  Column C\n" +
				"x         y         z",
		},
	} {
		t.Run(name, func(t *testing.T) {
			actual := Columnize(tc.input, tc.opts...)
			if actual != tc.expected {
				t.Fatalf("Expected:\n%s\n\nGot:\n%s", tc.expected, actual)
			}
		})
	}
}

func TestStringWidth(t *testing.T) {
	for s, expected := range map[string]int{
		"":  0,
		"x": 1,
	} {
		t.Run(s, func(t *testing.T) {
			if actual := stringWidth(s); actual != expected {
				t.Errorf("Expected %d, got %d", expected, actual)
			}
		})
	}
}

func BenchmarkColumnWidthCalculator(b *testing.B) {
	input := [][]string{
		{"UUID A", "UUID B", "UUID C", "Column D", "Column E"},
	}
	for i := 0; i < 1000; i++ {
		buf := make([]byte, 16)
		if _, err := crand.Read(buf); err != nil {
			b.Fatal(err.Error())
		}

		uuid := fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
			buf[0:4],
			buf[4:6],
			buf[6:8],
			buf[8:10],
			buf[10:16])

		input = append(input, []string{uuid[:8], uuid[:12], uuid, "short", "short"})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runtime.KeepAlive(Columnize(input))
	}
}
