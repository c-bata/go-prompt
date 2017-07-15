package prompt

import (
	"reflect"
	"testing"
)

func TestFormatCompletion(t *testing.T) {
	in := []string{
		"select",
		"from",
		"insert",
		"where",
	}
	ex := []string{
		"select",
		"from  ",
		"insert",
		"where ",
	}

	ac, width := formatCompletions(in)
	if !reflect.DeepEqual(ac, ex) {
		t.Errorf("Should be %#v, but got %#v", ex, ac)
	}
	if width != 6 {
		t.Errorf("Should be %#v, but got %#v", 4, width)
	}
}
