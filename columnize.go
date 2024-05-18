package columnize

import (
	"strings"
	"unsafe"

	runewidth "github.com/mattn/go-runewidth"
)

type config struct {
	prefix    string
	separator string
	rows      [][]string
}

type Option func(*config)

// WithPrefix sets the prefix before every line.
func WithPrefix(prefix string) Option {
	return func(c *config) {
		c.prefix = prefix
	}
}

// WithSeparator sets the separator between all columns.
func WithSeparator(separator string) Option {
	return func(c *config) {
		c.separator = separator
	}
}

// WithHeaders sets the first row as headers.
func WithHeaders(headers ...string) Option {
	return func(c *config) {
		c.rows = [][]string{headers}
	}
}

var _spaces = strings.Repeat(" ", 1024)

// Columnize is the public-facing interface that takes a list of strings and
// returns nicely aligned column-formatted text.
func Columnize(rows [][]string, opts ...Option) string {
	if len(rows) == 0 {
		return ""
	}

	config := config{
		prefix:    "",
		separator: "  ",
		rows:      nil,
	}
	for _, opt := range opts {
		opt(&config)
	}

	if config.rows != nil {
		rows = append(config.rows, rows...)
	}

	columns := 0
	for _, row := range rows {
		columns = max(columns, len(row))
	}

	// examine list of strings and determine how wide each column should be
	widths := make([]int, columns)
	for _, row := range rows {
		for i, elem := range row {
			widths[i] = max(widths[i], runewidth.StringWidth(elem))
		}
	}

	widthMax := 0
	widthSum := 0
	for _, w := range widths {
		widthMax = max(widthMax, w)
		widthSum += w
	}

	if widthMax > len(_spaces) {
		// NOTE: very rare case hopefully
		_spaces = strings.Repeat(" ", widthMax)
	}

	// estimate buffer size
	size := (len(config.separator)*(len(widths)-1) + widthSum) * len(rows)
	b := make([]byte, 0, size)
	for _, row := range rows {
		columns := len(row)
		b = append(b, config.prefix...)
		for i, elem := range row {
			if i == columns-1 {
				b = append(b, elem...)
				b = append(b, '\n')
			} else {
				cnt := widths[i] - runewidth.StringWidth(elem)
				b = append(b, elem...)
				b = append(b, _spaces[:cnt]...)
				b = append(b, config.separator...)
			}
		}
	}
	return unsafe.String(unsafe.SliceData(b), len(b)-1)
}
