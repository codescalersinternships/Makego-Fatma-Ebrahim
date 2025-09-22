package makego

import (
	"os"
	"reflect"
	"testing"
)

func createMakefile(filename string, bytes []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func deleteMakefile(filename string) error {
	os.Remove("hello")
	os.Remove("clean")
	os.Remove("hello.txt")
	return os.Remove(filename)
}

func TestParser(t *testing.T) {
	t.Run("test parse with no makefile", func(t *testing.T) {
		_, err := parseFile("Makefile")
		expected := "open Makefile: no such file or directory"
		if err.Error() != expected {
			t.Errorf("expected error: %v, got: %v", expected, err)
		}
	})

	t.Run("test parse makefile with syntax error", func(t *testing.T) {
		makeExample := []byte(`run:
 echo "Hello World" > hello.txt
clean: 
	rm -rf hello.txt`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		_, err := parseFile(filename)
		if err == nil {
			t.Errorf("expected error, got: no error")
		}

		expected := "missing correct separator in line: 2"
		if err.Error() != expected {
			t.Errorf("expected error: %v, got: %v", expected, err)
		}

		deleteMakefile(filename)
	})

	t.Run("test parse makefile with targets of no dependencies", func(t *testing.T) {
		makeExample := []byte(`
run:
	echo "Hello World" > hello.txt
clean: 
	rm -rf hello.txt`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		targets, err := parseFile(filename)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := map[string]Target{
			"default": {
				Name:         "default",
				Dependencies: []string{"run"},
				Commands:     []string{},
			},
			"run": {
				Name:         "run",
				Dependencies: []string{},
				Commands:     []string{"echo \"Hello World\" > hello.txt"},
			},
			"clean": {
				Name:         "clean",
				Dependencies: []string{},
				Commands:     []string{"rm -rf hello.txt"},
			},
		}

		if len(targets) != 3 {
			t.Errorf("expected 3 targets, got: %v", len(targets))
		}

		for tar_name, target := range targets {
			expected_target := expected[tar_name]
			if !reflect.DeepEqual(target, expected_target) {
				t.Errorf("expected target: %v, got: %v", expected_target, target)
			}
		}

		deleteMakefile(filename)
	})

	t.Run("test parse makefile with targets of dependencies", func(t *testing.T) {
		makeExample := []byte(`
run: clean
	echo "Hello World" > hello.txt
clean: 
	rm -rf hello.txt`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		targets, err := parseFile(filename)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		expected := map[string]Target{
			"default": {
				Name:         "default",
				Dependencies: []string{"run"},
				Commands:     []string{},
			},
			"run": {
				Name:         "run",
				Dependencies: []string{"clean"},
				Commands:     []string{"echo \"Hello World\" > hello.txt"},
			},
			"clean": {
				Name:         "clean",
				Dependencies: []string{},
				Commands:     []string{"rm -rf hello.txt"},
			},
		}

		if len(targets) != 3 {
			t.Errorf("expected 3 targets, got: %v", len(targets))
		}

		for tar_name, target := range targets {
			expected_target := expected[tar_name]
			if !reflect.DeepEqual(target, expected_target) {
				t.Errorf("expected target: %v, got: %v", expected_target, target)
			}
		}

		deleteMakefile(filename)
	})

}

func TestDependencies(t *testing.T) {
	t.Run("test makefile with targets of circular dependencies", func(t *testing.T) {
		makeExample := []byte(`
run: clean
	echo "Hello World" > hello.txt
clean: run
	rm -rf hello.txt`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		targets, err := parseFile(filename)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		cliTargets := []string{"default"}

		err = makeHandler(cliTargets, targets)
		if err == nil {
			t.Errorf("expected error, got: no error")
		}

		expected := "circular dependency detected"
		if err.Error() != expected {
			t.Errorf("expected error: %v, got: %v", expected, err)
		}

		deleteMakefile(filename)
	})

	t.Run("test makefile with targets of wrong dependencies and default target", func(t *testing.T) {
		makeExample := []byte(`
run: notclean
	echo "Hello World" > hello.txt
clean:
	rm -rf hello.txt`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		targets, err := parseFile(filename)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		cliTargets := []string{"default"}

		err = makeHandler(cliTargets, targets)
		if err == nil {
			t.Errorf("expected error, got: no error")
		}

		expected := "dependency 'notclean' not found"
		if err.Error() != expected {
			t.Errorf("expected error: %v, got: %v", expected, err)
		}

		bytes, err := os.ReadFile("hello.txt")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		fileContent := string(bytes)
		expectedFileContent := "Hello World\n"
		if fileContent != expectedFileContent {
			t.Errorf("expected file content: %s, got: %s", expectedFileContent, fileContent)
		}

		deleteMakefile(filename)
	})

	t.Run("test makefile with targets of dependencies and one target", func(t *testing.T) {
		makeExample := []byte(`
run: clean
	echo "Hello World" >> hello.txt
clean:
	echo "Clean should appear first" > hello.txt`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		targets, err := parseFile(filename)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		cliTargets := []string{"run"}

		err = makeHandler(cliTargets, targets)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		bytes, err := os.ReadFile("hello.txt")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		fileContent := string(bytes)
		expectedFileContent := "Clean should appear first\nHello World\n"
		if fileContent != expectedFileContent {
			t.Errorf("expected file content: %s, got: %s", expectedFileContent, fileContent)
		}

		deleteMakefile(filename)
	})

	t.Run("test makefile with targets of dependencies and multiple dependent targets", func(t *testing.T) {
		makeExample := []byte(`
run: clean
	echo "Hello World" >> hello.txt
clean:
	echo "Clean should appear first" > hello.txt`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		targets, err := parseFile(filename)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		cliTargets := []string{"run", "clean"}

		err = makeHandler(cliTargets, targets)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		bytes, err := os.ReadFile("hello.txt")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		fileContent := string(bytes)
		expectedFileContent := "Clean should appear first\nHello World\n"
		if fileContent != expectedFileContent {
			t.Errorf("expected file content: %s, got: %s", expectedFileContent, fileContent)
		}

		deleteMakefile(filename)
	})

	t.Run("test makefile with targets of dependencies and multiple independent target", func(t *testing.T) {
		makeExample := []byte(`
run:
	echo "Hello World" > hello
clean:
	echo "Hello World" > clean`)
		filename := "Makefile"
		createMakefile(filename, makeExample)
		targets, err := parseFile(filename)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		cliTargets := []string{"run", "clean"}

		err = makeHandler(cliTargets, targets)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		bytes, err := os.ReadFile("hello")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		fileContent := string(bytes)
		expectedFileContent := "Hello World\n"
		if fileContent != expectedFileContent {
			t.Errorf("expected file content: %s, got: %s", expectedFileContent, fileContent)
		}

		bytes, err = os.ReadFile("clean")
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}

		fileContent = string(bytes)
		expectedFileContent = "Hello World\n"
		if fileContent != expectedFileContent {
			t.Errorf("expected file content: %s, got: %s", expectedFileContent, fileContent)
		}

		deleteMakefile(filename)
	})

}
