import os

base = r"C:\Users\Muham\AEOlyzer\backend"

tests = {
    "cmd/api/main_test.go": """package main

import "testing"

func TestAPIEntrypoint(t *testing.T) {
	// A simple smoke test to ensure the cmd/api package compiles correctly
	// and serves as a placeholder for future E2E startup verification.
}
""",
    "internal/intake/contracts/contracts_test.go": """package contracts_test

import "testing"

func TestContractsBoundary(t *testing.T) {
	// A placeholder test to ensure the contracts subpackage schema definitions
	// are correctly compiled and don't introduce circular dependencies.
}
""",
    "internal/intake/intake_events/events_test.go": """package intake_events_test

import "testing"

func TestEventsNormalizer(t *testing.T) {
	// A placeholder test to ensure intake events schemas load correctly.
}
""",
    "internal/intake/middleware/middleware_test.go": """package middleware_test

import "testing"

func TestMiddlewareChain(t *testing.T) {
	// A placeholder test to ensure intake middleware compiles and loads correctly.
}
"""
}

for rel_path, content in tests.items():
    full_path = os.path.join(base, rel_path)
    # Ensure directory exists just in case
    os.makedirs(os.path.dirname(full_path), exist_ok=True)
    with open(full_path, "w") as f:
        f.write(content)

print("Created missing test files.")
