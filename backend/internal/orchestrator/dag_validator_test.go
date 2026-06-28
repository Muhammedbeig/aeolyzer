package orchestrator

import "testing"

func TestValidateDAGPlanRejectsCycleAndUnreachableNode(t *testing.T) {
	valid := DAGPlan{
		WorkflowID: "workflow-1",
		Version:    "1.0.0",
		EntryTask:  "start",
		Nodes: []TaskNode{
			{ID: "start", ActionClass: "read_context"},
			{ID: "finish", ActionClass: "present_result", Dependencies: []TaskID{"start"}},
		},
	}
	if err := ValidateDAGPlan(valid); err != nil {
		t.Fatalf("ValidateDAGPlan(valid) failed: %v", err)
	}
	cyclic := valid
	cyclic.Nodes = []TaskNode{
		{ID: "start", ActionClass: "read", Dependencies: []TaskID{"finish"}},
		{ID: "finish", ActionClass: "write", Dependencies: []TaskID{"start"}},
	}
	if err := ValidateDAGPlan(cyclic); err == nil {
		t.Fatal("ValidateDAGPlan() accepted cycle")
	}
	unreachable := valid
	unreachable.Nodes = append(unreachable.Nodes, TaskNode{
		ID:          "orphan",
		ActionClass: "unexpected",
	})
	if err := ValidateDAGPlan(unreachable); err == nil {
		t.Fatal("ValidateDAGPlan() accepted unreachable node")
	}
}
