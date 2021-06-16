package zenity

import (
	"reflect"
	"testing"
)

func TestFileFilters_simplify(t *testing.T) {
	tests := []struct {
		data FileFilters
		want []string
	}{
		{FileFilters{{"", []string{``}}}, nil},
		{FileFilters{{"", []string{`*.png`}}}, []string{"*.png"}},
		{FileFilters{{"", []string{`*.pn?`}}}, []string{"*.pn?"}},
		{FileFilters{{"", []string{`*.pn;`}}}, []string{"*.pn?"}},
		{FileFilters{{"", []string{`*.pn\?`}}}, nil},
		{FileFilters{{"", []string{`*.[PpNnGg]`}}}, []string{"*.?"}},
		{FileFilters{{"", []string{`*.[Pp][Nn][Gg]`}}}, []string{"*.PNG"}},
		{FileFilters{{"", []string{`*.[Pp][\Nn][G\g]`}}}, []string{"*.PNG"}},
		{FileFilters{{"", []string{`*.[PNG`}}}, []string{"*.[PNG"}},
		{FileFilters{{"", []string{`*.]PNG`}}}, []string{"*.]PNG"}},
		{FileFilters{{"", []string{`*.[[]PNG`}}}, []string{"*.[PNG"}},
		{FileFilters{{"", []string{`*.[]]PNG`}}}, []string{"*.]PNG"}},
		{FileFilters{{"", []string{`*.[\[]PNG`}}}, []string{"*.[PNG"}},
		{FileFilters{{"", []string{`*.[\]]PNG`}}}, []string{"*.]PNG"}},
		{FileFilters{{"", []string{`public.png`}}}, []string{"public.png"}},
	}
	for i, tt := range tests {
		tt.data.simplify()
		if got := tt.data[0].Patterns; !reflect.DeepEqual(got, tt.want) {
			t.Errorf("FileFilters.simplify[%d] = %q, want %q", i, got, tt.want)
		}
	}
}

func TestFileFilters_types(t *testing.T) {
	tests := []struct {
		data FileFilters
		want []string
	}{
		{FileFilters{{"", []string{``}}}, nil},
		{FileFilters{{"", []string{`*.png`}}}, []string{"", "png"}},
		{FileFilters{{"", []string{`*.pn?`}}}, nil},
		{FileFilters{{"", []string{`*.pn;`}}}, []string{"", "pn;"}},
		{FileFilters{{"", []string{`*.pn\?`}}}, []string{"", "pn?"}},
		{FileFilters{{"", []string{`*.[PpNnGg]`}}}, nil},
		{FileFilters{{"", []string{`*.[Pp][Nn][Gg]`}}}, []string{"", "PNG"}},
		{FileFilters{{"", []string{`*.[Pp][\Nn][G\g]`}}}, []string{"", "PNG"}},
		{FileFilters{{"", []string{`*.[PNG`}}}, []string{"", "[PNG"}},
		{FileFilters{{"", []string{`*.]PNG`}}}, []string{"", "]PNG"}},
		{FileFilters{{"", []string{`*.[[]PNG`}}}, []string{"", "[PNG"}},
		{FileFilters{{"", []string{`*.[]]PNG`}}}, []string{"", "]PNG"}},
		{FileFilters{{"", []string{`*.[\[]PNG`}}}, []string{"", "[PNG"}},
		{FileFilters{{"", []string{`*.[\]]PNG`}}}, []string{"", "]PNG"}},
		{FileFilters{{"", []string{`public.png`}}}, []string{"", "public.png"}},
	}
	for i, tt := range tests {
		if got := tt.data.types(); !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("FileFilters.types[%d] = %v, want %v", i, got, tt.want)
		}
	}
}
