package zenity

import (
	"reflect"
	"testing"
)

func TestFileFilters_name(t *testing.T) {
	tests := []struct {
		data FileFilters
		want string
	}{
		{FileFilters{{"", []string{`*.png`}, true}}, "*.png"},
		{FileFilters{{"", []string{`*.png`, `*.jpg`}, true}}, "*.png *.jpg"},
		{FileFilters{{"Image files", []string{`*.png`, `*.jpg`}, true}}, "Image files"},
	}
	for i, tt := range tests {
		tt.data.name()
		if got := tt.data[0].Name; got != tt.want {
			t.Errorf("FileFilters.name[%d] = %q; want %q", i, got, tt.want)
		}
	}
}

func TestFileFilters_simplify(t *testing.T) {
	tests := []struct {
		data []string
		want []string
	}{
		{[]string{``}, nil},
		{[]string{`*.\?`}, nil},
		{[]string{`*.png`}, []string{"*.png"}},
		{[]string{`*.pn?`}, []string{"*.pn?"}},
		{[]string{`*.pn;`}, []string{"*.pn?"}},
		{[]string{`*.[PpNnGg]`}, []string{"*.?"}},
		{[]string{`*.[Pp][Nn][Gg]`}, []string{"*.PNG"}},
		{[]string{`*.[Pp][\Nn][G\g]`}, []string{"*.PNG"}},
		{[]string{`*.[PNG`}, []string{"*.[PNG"}},
		{[]string{`*.]PNG`}, []string{"*.]PNG"}},
		{[]string{`*.[[]PNG`}, []string{"*.[PNG"}},
		{[]string{`*.[]]PNG`}, []string{"*.]PNG"}},
		{[]string{`*.[\[]PNG`}, []string{"*.[PNG"}},
		{[]string{`*.[\]]PNG`}, []string{"*.]PNG"}},
		{[]string{`public.png`}, []string{"public.png"}},
	}
	for i, tt := range tests {
		filters := FileFilters{FileFilter{Patterns: tt.data}}
		filters.simplify()
		if got := filters[0].Patterns; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("FileFilters.simplify[%d] = %q; want %q", i, got, tt.want)
		}
	}
}

func TestFileFilters_casefold(t *testing.T) {
	tests := []struct {
		data []string
		want []string
	}{
		{[]string{`*.png`}, []string{`*.[pP][nN][gG]`}},
		{[]string{`*.pn?`}, []string{`*.[pP][nN]?`}},
		{[]string{`*.pn;`}, []string{`*.[pP][nN];`}},
		{[]string{`*.pn\?`}, []string{`*.[pP][nN]\?`}},
		{[]string{`*.[PpNnGg]`}, []string{`*.[PppPNnnNGggG]`}},
		{[]string{`*.[Pp][Nn][Gg]`}, []string{`*.[PppP][NnnN][GggG]`}},
		{[]string{`*.[Pp][\Nn][G\g]`}, []string{`*.[PppP][\NnnN][Gg\gG]`}},
		{[]string{`*.[PNG`}, []string{`*.[PpNnGg`}},
		{[]string{`*.]PNG`}, []string{`*.][Pp][Nn][Gg]`}},
		{[]string{`*.[[]PNG`}, []string{`*.[[][Pp][Nn][Gg]`}},
		{[]string{`*.[]]PNG`}, []string{`*.[]][Pp][Nn][Gg]`}},
		{[]string{`*.[\[]PNG`}, []string{`*.[\[][Pp][Nn][Gg]`}},
		{[]string{`*.[\]]PNG`}, []string{`*.[\]][Pp][Nn][Gg]`}},
	}
	for i, tt := range tests {
		filters := FileFilters{FileFilter{Patterns: tt.data}}
		filters[0].CaseFold = true
		filters.casefold()
		if got := filters[0].Patterns; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("FileFilters.casefold[%d] = %q; want %q", i, got, tt.want)
		}
	}
}

func TestFileFilters_types(t *testing.T) {
	tests := []struct {
		data []string
		want []string
	}{
		{[]string{``}, nil},
		{[]string{`*.png`}, []string{"", "png"}},
		{[]string{`*.pn?`}, nil},
		{[]string{`*.pn;`}, []string{"", "pn;"}},
		{[]string{`*.pn\?`}, []string{"", "pn?"}},
		{[]string{`*.[PpNnGg]`}, nil},
		{[]string{`*.[Pp][Nn][Gg]`}, []string{"", "PNG"}},
		{[]string{`*.[Pp][\Nn][G\g]`}, []string{"", "PNG"}},
		{[]string{`*.[PNG`}, []string{"", "[PNG"}},
		{[]string{`*.]PNG`}, []string{"", "]PNG"}},
		{[]string{`*.[[]PNG`}, []string{"", "[PNG"}},
		{[]string{`*.[]]PNG`}, []string{"", "]PNG"}},
		{[]string{`*.[\[]PNG`}, []string{"", "[PNG"}},
		{[]string{`*.[\]]PNG`}, []string{"", "]PNG"}},
		{[]string{`public.png`}, []string{"", "public.png"}},
		{[]string{`-public-.png`}, []string{"", "png"}},
	}
	for i, tt := range tests {
		filters := FileFilters{FileFilter{Patterns: tt.data}}
		if got := filters.types(); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("FileFilters.types[%d] = %v; want %v", i, got, tt.want)
		}
	}
}
