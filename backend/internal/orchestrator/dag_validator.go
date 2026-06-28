package orchestrator

import (
	"errors"
	"fmt"
)

// ValidateDAGPlan verifies identity, references, a single reachable acyclic
// topology, and non-empty action classes.
func ValidateDAGPlan(plan DAGPlan) error {
	if plan.WorkflowID == "" ||
		plan.Version == "" ||
		plan.EntryTask == "" ||
		len(plan.Nodes) == 0 ||
		len(plan.Nodes) > 500 {
		return errors.New("dag plan is incomplete")
	}
	nodes := make(map[TaskID]TaskNode, len(plan.Nodes))
	for _, node := range plan.Nodes {
		if node.ID == "" || node.ActionClass == "" {
			return errors.New("dag node id and action class are required")
		}
		if _, duplicate := nodes[node.ID]; duplicate {
			return fmt.Errorf("duplicate dag node %q", node.ID)
		}
		nodes[node.ID] = node
	}
	if _, found := nodes[plan.EntryTask]; !found {
		return errors.New("dag entry task is missing")
	}
	dependents := make(map[TaskID][]TaskID, len(nodes))
	for _, node := range nodes {
		for _, dependency := range node.Dependencies {
			if dependency == node.ID {
				return fmt.Errorf("dag node %q depends on itself", node.ID)
			}
			if _, found := nodes[dependency]; !found {
				return fmt.Errorf("dag dependency %q is missing", dependency)
			}
			dependents[dependency] = append(dependents[dependency], node.ID)
		}
	}
	state := make(map[TaskID]uint8, len(nodes))
	var visit func(TaskID) error
	visit = func(id TaskID) error {
		switch state[id] {
		case 1:
			return errors.New("dag contains a cycle")
		case 2:
			return nil
		}
		state[id] = 1
		for _, next := range dependents[id] {
			if err := visit(next); err != nil {
				return err
			}
		}
		state[id] = 2
		return nil
	}
	if err := visit(plan.EntryTask); err != nil {
		return err
	}
	if len(state) != len(nodes) {
		return errors.New("dag contains nodes unreachable from entry task")
	}
	return nil
}
