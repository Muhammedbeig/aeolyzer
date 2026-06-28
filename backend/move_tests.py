import os
import shutil

base = r"C:\Users\Muham\AEOlyzer\backend\internal\orchestrator"

moves = {
    "content_workflow_test.go": ("handoff", "handoff_test"),
    "output_surface_test.go": ("handoff", "handoff_test"),
    "memory_update_test.go": ("requests", "requests_test"),
    "mode_gate_test.go": ("state", "state_test"),
}

for filename, (subpkg, pkgname) in moves.items():
    src = os.path.join(base, filename)
    if os.path.exists(src):
        dst = os.path.join(base, subpkg, filename)
        
        with open(src, "r") as f:
            content = f.read()
            
        content = content.replace("package orchestrator_test", f"package {pkgname}")
        
        with open(dst, "w") as f:
            f.write(content)
            
        os.remove(src)

print("Tests moved successfully.")
