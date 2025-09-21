package makego

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

func sortDependencies(target Target, alltargets map[string]Target, visited map[string]bool) ([]string, error) {
	dep_list := []string{}
	err := error(nil)
	for _, dep := range target.Dependencies {
		if visited[dep] {
			continue
		}
		visited[dep] = true
		_, ok := alltargets[dep]
		if !ok {
			err = fmt.Errorf("dependency '%s' not found", dep)
		}
		sorted_deps, _ := sortDependencies(alltargets[dep], alltargets, visited)

		dep_list = append(dep_list, sorted_deps...)
	}
	dep_list = append(dep_list, target.Name)
	return dep_list, err
}

func detectCycle(dep_list []string) bool {
	dep_map := make(map[string]bool)
	for _, dep := range dep_list {
		_, ok := dep_map[dep]
		if ok {
			return true
		}
		dep_map[dep] = true
	}
	return false

}

func getDependencies(targets map[string]Target) (map[string][]string, error) {
	dependency_map := make(map[string][]string)
	err := error(nil)
	for tar_name, target := range targets {
		visited := make(map[string]bool)
		dep_list, e := sortDependencies(target, targets, visited)
		if e != nil {
			err = e
		}
		if detectCycle(dep_list) {
			return nil, fmt.Errorf("cycle detected in target: %s", tar_name)
		}
		dependency_map[tar_name] = dep_list
	}
	return dependency_map, err
}

func getDefaultTarget(targets map[string]Target) (string, error) {
	flag.Parse()
	targetName := flag.Arg(0)
	if targetName == "" {
		targetName = "default"
	}
	_, ok := targets[targetName]
	if !ok {
		return "", fmt.Errorf("target '%s' not found", targetName)
	}
	return targetName, nil
}

func executeTarget(targets map[string]Target, dependencies []string) error {
	for _, dep := range dependencies {

		for _, command := range targets[dep].Commands {
			fmt.Println(command)
			cmd := exec.Command("sh", "-c", command)
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("error in command: '%s' in target: '%s' with message: %v", command,dep, err)
			}
		}
	}
	return nil
}

func Make() error {
	err := error(nil)
	targets, e := Parse()
	if e != nil {
		return e
	}
	fmt.Printf("%+v\n", targets)

	dep_map, e := getDependencies(targets)
	// ignore errors of incorrect dependencies
	if e != nil && !strings.HasPrefix(e.Error(), "dependency") {
		return e
	}
	if e != nil && strings.HasPrefix(e.Error(), "dependency") {
		err = e
	}
	fmt.Printf("%+v\n", dep_map)

	targetName, e := getDefaultTarget(targets)
	if e != nil {
		return e
	}
	fmt.Println(targetName)

	e = executeTarget(targets, dep_map[targetName])
	if e != nil {
		return e
	}
	return err
}
