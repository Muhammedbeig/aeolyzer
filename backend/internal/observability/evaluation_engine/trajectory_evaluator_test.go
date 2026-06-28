package evaluationengine

import "testing"

func TestEvaluateTrajectoryModes(t *testing.T) {
	actual := []ActionEvent{
		{Class: "brief_read"},
		{Class: "source_read"},
		{Class: "draft_section"},
		{Class: "quality_gate"},
	}
	tests := map[string]struct {
		mode     TrajectoryMode
		required []ActionSpec
		want     bool
	}{
		"exact pass": {
			mode: TrajectoryExact,
			required: []ActionSpec{
				{Class: "brief_read"},
				{Class: "source_read"},
				{Class: "draft_section"},
				{Class: "quality_gate"},
			},
			want: true,
		},
		"exact rejects extra action": {
			mode: TrajectoryExact,
			required: []ActionSpec{
				{Class: "brief_read"},
				{Class: "draft_section"},
				{Class: "quality_gate"},
			},
		},
		"in order pass": {
			mode: TrajectoryInOrder,
			required: []ActionSpec{
				{Class: "brief_read"},
				{Class: "draft_section"},
				{Class: "quality_gate"},
			},
			want: true,
		},
		"in order rejects reversal": {
			mode: TrajectoryInOrder,
			required: []ActionSpec{
				{Class: "draft_section"},
				{Class: "brief_read"},
			},
		},
		"any order pass": {
			mode: TrajectoryAnyOrder,
			required: []ActionSpec{
				{Class: "quality_gate"},
				{Class: "brief_read"},
			},
			want: true,
		},
		"any order requires duplicates": {
			mode: TrajectoryAnyOrder,
			required: []ActionSpec{
				{Class: "brief_read"},
				{Class: "brief_read"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			score, err := EvaluateTrajectory(AgentTrace{
				TraceID: "trace-1",
				Actions: actual,
			}, TrajectorySpec{
				Mode:     test.mode,
				Required: test.required,
			})
			if err != nil {
				t.Fatalf("EvaluateTrajectory() failed: %v", err)
			}
			if score.Pass != test.want {
				t.Fatalf("EvaluateTrajectory().Pass = %t, want %t", score.Pass, test.want)
			}
		})
	}
}

func TestEvaluateTrajectoryRejectsForbiddenAction(t *testing.T) {
	score, err := EvaluateTrajectory(AgentTrace{
		TraceID: "trace-1",
		Actions: []ActionEvent{{Class: "direct_canvas_write"}},
	}, TrajectorySpec{
		Mode:      TrajectoryForbiddenAbsent,
		Forbidden: []ActionSpec{{Class: "direct_canvas_write"}},
	})
	if err != nil {
		t.Fatalf("EvaluateTrajectory() failed: %v", err)
	}
	if score.Pass {
		t.Fatal("EvaluateTrajectory().Pass = true, want false")
	}
	if score.ForbiddenAbsent {
		t.Fatal("EvaluateTrajectory().ForbiddenAbsent = true, want false")
	}
}

func TestEvaluateTrajectoryRejectsInvalidInput(t *testing.T) {
	tests := map[string]struct {
		trace AgentTrace
		spec  TrajectorySpec
	}{
		"missing trace": {
			spec: TrajectorySpec{
				Mode:     TrajectoryExact,
				Required: []ActionSpec{{Class: "read"}},
			},
		},
		"unknown mode": {
			trace: AgentTrace{TraceID: "trace-1"},
			spec:  TrajectorySpec{Mode: "SOMETIMES"},
		},
		"empty exact": {
			trace: AgentTrace{TraceID: "trace-1"},
			spec:  TrajectorySpec{Mode: TrajectoryExact},
		},
		"empty observed action": {
			trace: AgentTrace{
				TraceID: "trace-1",
				Actions: []ActionEvent{{}},
			},
			spec: TrajectorySpec{
				Mode:     TrajectoryExact,
				Required: []ActionSpec{{Class: "read"}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := EvaluateTrajectory(test.trace, test.spec); err == nil {
				t.Fatal("EvaluateTrajectory() returned nil error")
			}
		})
	}
}
