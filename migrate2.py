import os
import shutil
import glob

base_dir = r"C:\Users\Muham\AEOlyzer\backend"

def move_and_update(src_name, dst_name, old_import, new_import, old_pkg, new_pkg):
    src_dir = os.path.join(base_dir, src_name)
    dst_dir = os.path.join(base_dir, dst_name)
    
    # Create dst_dir parent if needed
    os.makedirs(os.path.dirname(dst_dir), exist_ok=True)
    
    if os.path.exists(src_dir):
        if os.path.exists(dst_dir):
            for item in os.listdir(src_dir):
                shutil.move(os.path.join(src_dir, item), os.path.join(dst_dir, item))
            shutil.rmtree(src_dir)
        else:
            shutil.move(src_dir, dst_dir)
    
    # Update all .go files in backend for imports and type usages
    go_files = glob.glob(os.path.join(base_dir, "**", "*.go"), recursive=True)
    
    for file_path in go_files:
        with open(file_path, "r", encoding="utf-8") as f:
            content = f.read()
            
        new_content = content.replace(f'"{old_import}"', f'"{new_import}"')
        new_content = new_content.replace(f'"{old_import}/', f'"{new_import}/')
        
        # We also need to update the package usage, e.g. layer_02_intake. to intake.
        old_pkg_prefix = old_pkg + "."
        new_pkg_prefix = new_pkg + "."
        new_content = new_content.replace(old_pkg_prefix, new_pkg_prefix)
        
        # If this file is inside the moved directory, update its package declaration
        if file_path.startswith(dst_dir):
            new_content = new_content.replace(f"package {old_pkg}\n", f"package {new_pkg}\n")
            new_content = new_content.replace(f"package {old_pkg}_test\n", f"package {new_pkg}_test\n")

        if new_content != content:
            with open(file_path, "w", encoding="utf-8") as f:
                f.write(new_content)

# Move layer 3 back
move_and_update(
    src_name="layer_03_orchestration",
    dst_name=r"internal\orchestrator",
    old_import="aeolyzer/layer_03_orchestration",
    new_import="aeolyzer/internal/orchestrator",
    old_pkg="orchestration",
    new_pkg="orchestrator"
)

# Move layer 2
move_and_update(
    src_name="layer_02_intake",
    dst_name=r"internal\intake",
    old_import="aeolyzer/layer_02_intake",
    new_import="aeolyzer/internal/intake",
    old_pkg="layer_02_intake",
    new_pkg="intake"
)

print("Migration to internal/* completed.")
