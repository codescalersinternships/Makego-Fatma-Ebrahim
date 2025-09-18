package makego

import (
	"fmt"
)

func sortTargets(target Target, alltargets map[string]Target, visited map[string]bool) []string {
	dep_list := []string{}
	for _, dep := range target.Dependencies {
		if visited[dep] {
			continue
		}
		visited[dep] = true
		dep_list = append(dep_list, sortTargets(alltargets[dep], alltargets, visited)...)
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

func Make() error {
	targets, err := Parse()
	if err != nil {
		return err
	}
	fmt.Println(targets)
	dependency_map := make(map[string][]string)
	for tar_name, target := range targets {
		visited := make(map[string]bool)
		dep_list := sortTargets(target, targets, visited)
		if detectCycle(dep_list) {
			return fmt.Errorf("cycle detected in target %s", tar_name)
		}
		dependency_map[tar_name] = dep_list
	}

	fmt.Printf("%+v\n", dependency_map)
	return nil
}
