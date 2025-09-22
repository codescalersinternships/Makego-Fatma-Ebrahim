package makego

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"
	"sync"
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
			return nil, fmt.Errorf("circular dependency detected")
		}
		dependency_map[tar_name] = dep_list
	}
	return dependency_map, err
}

func getDefaultTarget(targets map[string]Target) ([]string, error) {
	cliTargets := make([]string, 0)
	flag.Parse()
	for arg := range flag.Args() {
		targetName := flag.Arg(arg)
		_, ok := targets[targetName]
		if !ok {
			return nil, fmt.Errorf("target '%s' not found", targetName)
		}
		cliTargets = append(cliTargets, targetName)
	}
	if len(cliTargets) == 0 {
		cliTargets = append(cliTargets, "default")
	}
	return cliTargets, nil
}

func executeTarget(targets map[string]Target, dependencies []string, executed map[string]bool, mu *sync.Mutex) error {
	for _, dep := range dependencies {
		mu.Lock()
		_, ok := executed[dep]
		if ok {
			mu.Unlock()
			continue
		}
		executed[dep] = true
		for _, command := range targets[dep].Commands {
			fmt.Println(command)
			cmd := exec.Command("sh", "-c", command)
			err := cmd.Run()
			if err != nil {
				return fmt.Errorf("error in command: '%s' in target: '%s' with message: %v", command, dep, err)
			}
		}
		mu.Unlock()

	}
	return nil
}

func executeTargetsConcurrently(targets map[string]Target, cliTargets []string, dependencies map[string][]string) error {
	wg := sync.WaitGroup{}
	executed := make(map[string]bool)
	mu := &sync.Mutex{}
	for _, cliTarget := range cliTargets {
		wg.Go(func() { executeTarget(targets, dependencies[cliTarget], executed, mu) })
	}
	wg.Wait()
	return nil
}

func makeHandler(cliTargets []string, targets map[string]Target) error {
	err := error(nil)

	dep_map, e := getDependencies(targets)
	// ignore errors of incorrect dependencies
	if e != nil && !strings.HasPrefix(e.Error(), "dependency") {
		return e
	}
	if e != nil && strings.HasPrefix(e.Error(), "dependency") {
		err = e
	}

	e = executeTargetsConcurrently(targets, cliTargets, dep_map)
	if e != nil {
		return e
	}
	return err
}

// Make is a function that parse a makefile specified by the command line flag -f 
// or the makefile in the current directory 
// then execute the targets specified by the command line arguments
// or execute the default target
func Make() error {
	targets, e := parse()
	if e != nil {
		return e
	}

	cliTargets, e := getDefaultTarget(targets)
	if e != nil {
		return e
	}

	return makeHandler(cliTargets, targets)
}
