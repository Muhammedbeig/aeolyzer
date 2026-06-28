import os
import re

spec_path = r'C:\Users\Muham\AEOlyzer\backend\specs\Layer 3 specs.md'
base_dir = r'C:\Users\Muham\AEOlyzer\backend'

with open(spec_path, 'r', encoding='utf-8') as f:
    content = f.read()

# Create directories
os.makedirs(os.path.join(base_dir, 'config', 'capability_profiles'), exist_ok=True)
os.makedirs(os.path.join(base_dir, 'workflows'), exist_ok=True)

# Regex to find sections and their yaml blocks
# E.g. ### 5.1 `content_collaborator.yaml`
# ```yaml ... ```
pattern = re.compile(r'### [^\n`]+`([^`]+\.yaml)`.*?```yaml\n(.*?)\n```', re.DOTALL)
matches = pattern.findall(content)

for filename, yaml_content in matches:
    if filename in ['content_collaborator.yaml', 'content_execution_guard.yaml']:
        out_path = os.path.join(base_dir, 'config', 'capability_profiles', filename)
    else:
        out_path = os.path.join(base_dir, 'workflows', filename)
    
    with open(out_path, 'w', encoding='utf-8') as f:
        f.write(yaml_content)
        
print(f"Extracted and wrote {len(matches)} YAML files from the spec.")
