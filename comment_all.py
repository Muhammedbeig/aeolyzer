import os
import sys

def comment_files(directory):
    extensions = {'.go', '.ts', '.tsx', '.js', '.jsx'}
    exclude_dirs = {'node_modules', '.git'}

    comment_header = """// ==========================================
// AEOlyzer Source Code
// This file is part of the AEOlyzer project.
// Please refer to the relevant layer specs and AGENTS.md for detailed documentation.
// ==========================================
"""

    count = 0
    for root, dirs, files in os.walk(directory):
        # modify dirs in place to skip excluded directories
        dirs[:] = [d for d in dirs if d not in exclude_dirs]
        for f in files:
            ext = os.path.splitext(f)[1]
            if ext in extensions:
                filepath = os.path.join(root, f)
                with open(filepath, 'r', encoding='utf-8') as file:
                    content = file.read()
                
                # Check if it already has our comment to avoid double commenting
                if not content.startswith("// =========================================="):
                    with open(filepath, 'w', encoding='utf-8') as file:
                        file.write(comment_header + content)
                    count += 1
    return count

backend_count = comment_files(r'C:\Users\Muham\AEOlyzer\backend')
frontend_count = comment_files(r'C:\Users\Muham\AEOlyzer\frontend')

print(f"Added comments to {backend_count} backend files and {frontend_count} frontend files.")
