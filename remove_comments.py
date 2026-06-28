import os

def remove_comments(directory):
    extensions = {'.go', '.ts', '.tsx', '.js', '.jsx'}
    exclude_dirs = {'node_modules', '.git'}

    comment_header = """// ==========================================
// AEOlyzer Source Code
// This file is part of the AEOlyzer project.
// Please refer to the relevant layer specs and AGENTS.md for detailed documentation.
// ==========================================
"""

    for root, dirs, files in os.walk(directory):
        dirs[:] = [d for d in dirs if d not in exclude_dirs]
        for f in files:
            ext = os.path.splitext(f)[1]
            if ext in extensions:
                filepath = os.path.join(root, f)
                try:
                    with open(filepath, 'r', encoding='utf-8') as file:
                        content = file.read()
                    
                    if content.startswith(comment_header):
                        with open(filepath, 'w', encoding='utf-8') as file:
                            file.write(content[len(comment_header):])
                except Exception as e:
                    pass

remove_comments(r'C:\Users\Muham\AEOlyzer\backend')
remove_comments(r'C:\Users\Muham\AEOlyzer\frontend')
print("Headers removed.")
