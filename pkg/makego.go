package makego

import (
	"fmt"
)

func sortTargets(target Target, alltargets map[string]Target) []string {
	dep_list := []string{}
	for _, dep := range target.Dependencies {
		fmt.Println("in:",target.Name,dep,dep_list)
		dep_list = append(dep_list, sortTargets(alltargets[dep], alltargets)...)
	}
	fmt.Println("outttt:",target.Name,dep_list)
	
	dep_list = append(dep_list, target.Name)
	return dep_list
}

func Make() error {
	targets, err := Parse()
	if err != nil {
		return err
	}

	fmt.Println(targets)
	for tar_name, target := range targets {
		dep_list := sortTargets(target, targets)
		fmt.Println("						out:",tar_name, dep_list)
	}

	


	return nil
}
