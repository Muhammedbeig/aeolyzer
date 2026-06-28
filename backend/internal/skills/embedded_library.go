package skills

import "embed"

// embeddedLibrary is the immutable Layer 4 library deployed with the process.
//
//go:embed skill-registry.yaml skills
var embeddedLibrary embed.FS
