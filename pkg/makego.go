package makego

import (
	"flag"
	"fmt"
	"os/exec"
)

func sortDependencies(target Target, alltargets map[string]Target, visited map[string]bool) []string {
	dep_list := []string{}
	for _, dep := range target.Dependencies {
		if visited[dep] {
			continue
		}
		visited[dep] = true
		dep_list = append(dep_list, sortDependencies(alltargets[dep], alltargets, visited)...)
	}
	dep_list = append(dep_list, target.Name)
	return dep_list
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
	for tar_name, target := range targets {
		visited := make(map[string]bool)
		dep_list := sortDependencies(target, targets, visited)
		if detectCycle(dep_list) {
			return nil, fmt.Errorf("cycle detected in target %s", tar_name)
		}
		dependency_map[tar_name] = dep_list
	}
	return dependency_map, nil
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
				return err
			}
		}
	}
	return nil
}

func Make() error {
	targets, err := Parse()
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", targets)

	dep_map, err := getDependencies(targets)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", dep_map)

	targetName, err := getDefaultTarget(targets)
	if err != nil {
		return err
	}
	fmt.Println(targetName)

	err = executeTarget(targets, dep_map[targetName])
	if err != nil {
		return err
	}
	return nil
}
