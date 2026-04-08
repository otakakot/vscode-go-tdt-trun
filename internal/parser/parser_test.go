package parser_test

import (
	"testing"

	"github.com/otakakot/vscode-go-tdt-trun/internal/parser"
)

func TestExtractSubtests_SliceKeyValue(t *testing.T) {
	subTests, err := parser.ExtractSubTests("testdata/slice_kv_test.go")
	if err != nil {
		t.Fatal(err)
	}

	expected := []struct {
		name     string
		funcName string
	}{
		{name: "both positive", funcName: "TestSliceKV"},
		{name: "positive and negative", funcName: "TestSliceKV"},
		{name: "both negative", funcName: "TestSliceKV"},
		{name: "both zero", funcName: "TestSliceKV"},
	}

	if len(subTests) != len(expected) {
		t.Fatalf("got %d subtests, want %d", len(subTests), len(expected))
	}

	for i, sub := range subTests {
		if sub.Name != expected[i].name {
			t.Errorf("subtests[%d].Name = %q, want %q", i, sub.Name, expected[i].name)
		}

		if sub.Func != expected[i].funcName {
			t.Errorf("subtests[%d].Func = %q, want %q", i, sub.Func, expected[i].funcName)
		}
	}
}

func TestExtractSubtests_SlicePositional(t *testing.T) {
	subTests, err := parser.ExtractSubTests("testdata/slice_positional_test.go")
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"first case", "second case"}
	if len(subTests) != len(expected) {
		t.Fatalf("got %d subtests, want %d", len(subTests), len(expected))
	}

	for i, sub := range subTests {
		if sub.Name != expected[i] {
			t.Errorf("subTests[%d].Name = %q, want %q", i, sub.Name, expected[i])
		}
	}
}

func TestExtractSubtests_Map(t *testing.T) {
	subTests, err := parser.ExtractSubTests("testdata/map_test.go")
	if err != nil {
		t.Fatal(err)
	}

	expectedNames := map[string]bool{
		"empty string": true,
		"single char":  true,
	}

	if len(subTests) != len(expectedNames) {
		t.Fatalf("got %d subtests, want %d", len(subTests), len(expectedNames))
	}

	for _, sub := range subTests {
		if !expectedNames[sub.Name] {
			t.Errorf("unexpected subtest name: %q", sub.Name)
		}

		if sub.Func != "TestMap" {
			t.Errorf("Func = %q, want %q", sub.Func, "TestMap")
		}
	}
}

func TestExtractSubtests_Literal(t *testing.T) {
	subTests, err := parser.ExtractSubTests("testdata/literal_test.go")
	if err != nil {
		t.Fatal(err)
	}

	if len(subTests) != 0 {
		t.Fatalf("got %d subtests, want 0 (literal t.Run should be skipped)", len(subTests))
	}
}
