package makego

import (
	"fmt"
	"os"
	"strings"
)

type Target struct {
	Name         string
	Dependencies []string
	Commands     []string
}

func getMakefile() (string, error) {
	file, err := os.Open("Makefile")
	if err != nil {
		return "", err
	}
	defer file.Close()
	return file.Name(), nil

}

func targetArrayToMap(targets []Target) map[string]Target {
	target_map := make(map[string]Target)
	for _, target := range targets {
		old_target, ok := target_map[target.Name] 
		if ok {
			target_map[target.Name]= Target {
				Name: old_target.Name,
				Dependencies: append(old_target.Dependencies, target.Dependencies...),
				Commands: target.Commands,
			}
		} else {
			target_map[target.Name]= target
		}
	}
	return target_map
}

func Parse() (map[string]Target,error) {
	filename, err := getMakefile()
	if err != nil {
		return nil,err
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil,err
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
				return nil,fmt.Errorf("syntax error in line: %d", i+1)
			}
		}
	}

	targets = append(targets, target)
	target_map := targetArrayToMap(targets)
	return target_map,nil

}
