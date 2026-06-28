package runtime

import (
	_ "embed"
	"fmt"

	"aeolyzer/internal/schemavalidation"
)

var (
	//go:embed runtime-execution.schema.json
	runtimeExecutionSchema []byte
	//go:embed sandbox-lease.schema.json
	sandboxLeaseSchema []byte
	//go:embed runtime-result.schema.json
	runtimeResultSchema []byte
	//go:embed jit-token.schema.json
	jitTokenSchema []byte
	//go:embed quarantine-command.schema.json
	quarantineCommandSchema []byte
	//go:embed dependency-policy.schema.json
	dependencyPolicySchema []byte
	//go:embed filesystem-policy.schema.json
	filesystemPolicySchema []byte
	//go:embed egress-policy.schema.json
	egressPolicySchema []byte
)

// CompileSchemas compiles every Layer 6 JSON Schema.
func CompileSchemas() error {
	schemas := map[string][]byte{
		"runtime execution":  runtimeExecutionSchema,
		"sandbox lease":      sandboxLeaseSchema,
		"runtime result":     runtimeResultSchema,
		"jit token":          jitTokenSchema,
		"quarantine command": quarantineCommandSchema,
		"dependency policy":  dependencyPolicySchema,
		"filesystem policy":  filesystemPolicySchema,
		"egress policy":      egressPolicySchema,
	}
	for name, schema := range schemas {
		if _, err := schemavalidation.Compile(schema); err != nil {
			return fmt.Errorf("compile %s schema: %w", name, err)
		}
	}
	return nil
}
