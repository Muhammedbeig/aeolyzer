import os
import shutil
import glob

src_dir = r"C:\Users\Muham\AEOlyzer\backend\internal\orchestrator"
dst_dir = r"C:\Users\Muham\AEOlyzer\backend\layer_03_orchestration"

# Move all files and folders
for item in os.listdir(src_dir):
    s = os.path.join(src_dir, item)
    d = os.path.join(dst_dir, item)
    if os.path.exists(d):
        if os.path.isdir(s):
            for sub_item in os.listdir(s):
                shutil.move(os.path.join(s, sub_item), os.path.join(d, sub_item))
        else:
            shutil.move(s, d)
    else:
        shutil.move(s, d)

# Clean up src_dir
shutil.rmtree(src_dir)

# Update imports and package names
go_files = glob.glob(os.path.join(dst_dir, "**", "*.go"), recursive=True)

for file_path in go_files:
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    # Update packages
    content = content.replace("package orchestrator\n", "package orchestration\n")
    content = content.replace("package orchestrator_test\n", "package orchestration_test\n")

    # Update imports
    content = content.replace('"aeolyzer/internal/orchestrator"', '"aeolyzer/layer_03_orchestration"')
    content = content.replace('"aeolyzer/internal/orchestrator/state"', '"aeolyzer/layer_03_orchestration/state"')
    content = content.replace('"aeolyzer/internal/orchestrator/handoff"', '"aeolyzer/layer_03_orchestration/handoff"')
    content = content.replace('"aeolyzer/internal/orchestrator/requests"', '"aeolyzer/layer_03_orchestration/requests"')

    # Update type references
    content = content.replace("orchestrator.", "orchestration.")

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)

print("Files moved and namespaces updated successfully.")
