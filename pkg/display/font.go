package display

import "github.com/realency/arke/pkg/bits"

// Font represents a mapping from a rune to a graphical representation in the form of binary pixels.
type Font func(r rune) *bits.Matrix
