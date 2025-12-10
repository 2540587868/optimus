package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerator_Run_GoldenFile(t *testing.T) {
	wd, _ := os.Getwd()
	testDataDir := filepath.Join(wd, "testdata")
	inputFile := filepath.Join(testDataDir, "simple.go")
	goldenFile := filepath.Join(testDataDir, "simple_golden.txt")

	generatedFileName := "simpleconfig_options.go"
	generatedFilePath := filepath.Join(wd, generatedFileName)

	defer os.Remove(generatedFilePath)

	g := Generator{
		StructName: "SimpleConfig",
		FileName:   inputFile,
		Package:    "testdata",
	}

	err := g.Run()
	if err != nil {
		t.Fatalf("Generator.Run() failed: %v", err)
	}

	actualBytes, err := os.ReadFile(generatedFilePath)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	expectedBytes, err := os.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("Failed to read golden file: %v", err)
	}

	actual := strings.TrimSpace(string(actualBytes))
	expected := strings.TrimSpace(string(expectedBytes))

	if actual != expected {
		t.Errorf("Generated code does not match golden file.\n")
		t.Errorf("EXPECTED:\n%s\n", expected)
		t.Errorf("--------------------------------------------------\n")
		t.Errorf("ACTUAL:\n%s\n", actual)

		// os.WriteFile(goldenFile, actualBytes, 0644)
	}
}
