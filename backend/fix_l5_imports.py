import os

base = r"C:\Users\Muham\AEOlyzer\backend"

def replace_in_file(filepath, old, new):
    if not os.path.exists(filepath):
        return
    with open(filepath, 'r') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(filepath, 'w') as f:
            f.write(content)

# 1. Fix imports in httpapi and cmd
replace_in_file(os.path.join(base, "internal", "httpapi", "handler.go"), 
                "aeolyzer/layer_05_extensions", "aeolyzer/internal/extensions")
replace_in_file(os.path.join(base, "internal", "httpapi", "handler_test.go"), 
                "aeolyzer/layer_05_extensions", "aeolyzer/internal/extensions")
replace_in_file(os.path.join(base, "cmd", "api", "main.go"), 
                "aeolyzer/layer_05_extensions", "aeolyzer/internal/extensions")
replace_in_file(os.path.join(base, "internal", "orchestrator", "onboarding.go"), 
                "aeolyzer/layer_05_extensions", "aeolyzer/internal/extensions")
replace_in_file(os.path.join(base, "internal", "orchestrator", "onboarding_test.go"), 
                "aeolyzer/layer_05_extensions", "aeolyzer/internal/extensions")

# 2. Fix the package issues for Layer 5 Go files
# Since types.go is in internal/extensions (package extensions) and presentation_intent_validator is in surface_router
# Let's move presentation_intent_validator.go and types.go to the correct places and fix their packages.
os.rename(
    os.path.join(base, "internal", "extensions", "surface_router", "presentation_intent_validator.go"),
    os.path.join(base, "internal", "extensions", "presentation_intent_validator.go")
)

# And the test should be in internal/extensions/tests. Let's move it to internal/extensions/ for now, or just leave it in tests and import aeolyzer/internal/extensions.
# Wait, tests/presentation_schema_test.go imports `aeolyzer/internal/extensions` and calls extensions.ValidatePresentationIntent. That is correct!
# The only issue was that presentation_intent_validator.go was in the surface_router folder but declared as package extensions, causing it to not be found in extensions package. 
# By moving it to internal/extensions/, it will be part of the extensions package.

print("Imports and packages fixed.")
