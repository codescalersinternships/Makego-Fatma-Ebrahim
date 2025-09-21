package makego

import (
	"fmt"
	"os"
	"strings"
	"flag"
)

type Target struct {
	Name         string
	Dependencies []string
	Commands     []string
}

func getMakefile() (string, error) {
	var makefile string
	flag.StringVar(&makefile, "f", "Makefile", "custom makefile path")
	flag.Parse()

	fmt.Println(makefile)
	file, err := os.Open(makefile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	return file.Name(), nil

}

func merge(old []string, new []string) []string {
	merged := make(map[string]bool)
	for _, old_dep := range old {
		merged[old_dep] = true
	}

	for _, new_dep := range new {
		merged[new_dep] = true
	}

	merged_list := make([]string, 0)
	for dep := range merged {
		merged_list = append(merged_list, dep)
	}
	return merged_list
}

func targetArrayToMap(targets []Target) map[string]Target {
	target_map := make(map[string]Target)
	for _, target := range targets {
		old_target, ok := target_map[target.Name]
		if ok {
			target_map[target.Name] = Target{
				Name:         old_target.Name,
				Dependencies: merge(old_target.Dependencies, target.Dependencies),
				Commands:     target.Commands,
			}
		} else {
			target_map[target.Name] = target
		}
	}
	_, ok := target_map["default"]
	if !ok {
		target_map["default"] = targets[0]
	}
	return target_map
}

func Parse() (map[string]Target, error) {
	filename, err := getMakefile()
	if err != nil {
		return nil, err
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(file), "\n")

	targets := []Target{}
	targetDetected := false
	var target Target
	for i, line := range lines {
		if line != "" {
			if !strings.HasPrefix(line, "\t") && strings.Contains(line, ":") {
				// fmt.Println("is target:", line)
				targetDetected = true
				targetName := strings.TrimSpace(line[:strings.Index(line, ":")])
				dependencies := strings.TrimSpace(line[strings.Index(line, ":")+1:])
				depArray := strings.Fields(dependencies)
				if target.Name != "" {
					targets = append(targets, target)
				}
				target = Target{Name: targetName, Dependencies: depArray, Commands: []string{}}
				continue
			}
			if strings.HasPrefix(line, "\t") && targetDetected {
				// fmt.Println("is command:", line)
				target.Commands = append(target.Commands, line[1:])
				continue
			}
			if !strings.HasPrefix(line, "\t") && !strings.Contains(line, ":") {
				// fmt.Println("is error:", line)
				return nil, fmt.Errorf("missing correct separator in line: %d", i+1)
			}
		}
	}

	targets = append(targets, target)
	target_map := targetArrayToMap(targets)
	return target_map, nil

}
