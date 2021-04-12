package zenity

import (
	"reflect"
	"testing"
)

func TestFileFilters_darwin(t *testing.T) {
	tests := []struct {
		data FileFilters
		want []string
	}{
		{FileFilters{{"", []string{`*.png`}}}, []string{`png`}},
		{FileFilters{{"", []string{`*.pn?`}}}, nil},
		{FileFilters{{"", []string{`*.pn\?`}}}, []string{`pn?`}},
		{FileFilters{{"", []string{`*.[PpNnGg]`}}}, nil},
		{FileFilters{{"", []string{`*.[Pp][Nn][Gg]`}}}, []string{`PNG`}},
		{FileFilters{{"", []string{`*.[Pp][\Nn][G\g]`}}}, []string{`PNG`}},
		{FileFilters{{"", []string{`*.[PNG`}}}, []string{`[PNG`}},
		{FileFilters{{"", []string{`*.]PNG`}}}, []string{`]PNG`}},
		{FileFilters{{"", []string{`*.[[]PNG`}}}, []string{`[PNG`}},
		{FileFilters{{"", []string{`*.[]]PNG`}}}, []string{`]PNG`}},
		{FileFilters{{"", []string{`*.[\[]PNG`}}}, []string{`[PNG`}},
		{FileFilters{{"", []string{`*.[\]]PNG`}}}, []string{`]PNG`}},
	}
	for _, tt := range tests {
		if got := tt.data.darwin(); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("FileFilters.darwin(%+v) = %v, want %v", tt.data, got, tt.want)
		}
	}
}
