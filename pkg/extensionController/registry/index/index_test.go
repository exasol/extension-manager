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

func TestGetExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		id       string
		expected Extension
	}{
		{"empty json", `{}`, "id", Extension{}},
		{"empty extension field", `{"extensions":[]}`, "id", Extension{}},
		{"single entry wrong id", `{"extensions":[{"id": "ext1"}]}`, "wrong-id", Extension{}},
		{"single entry correct id", `{"extensions":[{"id": "ext1", "url": "http://url"}]}`, "ext1", Extension{ID: "ext1", URL: "http://url"}},
		{"multiple entries wrong id", `{"extensions":[{"id": "ext1"},{"id": "ext2"},{"id": "ext3"}]}`, "wrong-id", Extension{}},
		{"multiple entries correct id", `{"extensions":[{"id": "ext1"},{"id": "ext2", "url": "ext2-url"},{"id": "ext3"}]}`, "ext2", Extension{ID: "ext2", URL: "ext2-url"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			index, err := Decode(strings.NewReader(test.input))
			if err != nil {
				t.Errorf("unexpected error %v", err)
			} else {
				actual, ok := index.GetExtension(test.id)
				if ok != (test.expected.URL != "") {
					t.Errorf("expected ok to be %v", test.expected.URL != "")
				}
				if !reflect.DeepEqual(test.expected, actual) {
					t.Errorf("expected %v but got %v for input %q", test.expected, actual, test.input)
				}
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
