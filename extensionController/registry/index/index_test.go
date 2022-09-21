package index

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetExtensionIDs(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedIds []string
	}{
		{"empty json", `{}`, []string{}},
		{"empty extension field", `{"extensions":[]}`, []string{}},
		{"single entry", `{"extensions":[{"id": "ext1"}]}`, []string{"ext1"}},
		{"ignore unexpected fields", `{"unexpectedField": true, "extensions":[{"id": "ext1"}]}`, []string{"ext1"}},
		{"multiple entries", `{"extensions":[{"id": "ext1"},{"id": "ext2"},{"id": "ext3"}]}`, []string{"ext1", "ext2", "ext3"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index, err := Decode(strings.NewReader(test.input))
			if err != nil {
				t.Errorf("unexpected error %v", err)
			} else if actual := index.GetExtensionIDs(); !reflect.DeepEqual(test.expectedIds, actual) {
				t.Errorf("expected %v but got %v for input %q", test.expectedIds, actual, test.input)
			}
		})
	}
}

func TestDecodeFails(t *testing.T) {
	_, err := Decode(strings.NewReader(""))
	if err == nil {
		t.Error("expected decode to fail")
	} else if err.Error() != "failed to decode registry content: EOF" {
		t.Errorf("got wrong error: %q", err.Error())
	}
}
