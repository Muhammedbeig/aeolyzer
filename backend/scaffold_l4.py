import os
import re

base_dir = r"C:\Users\Muham\AEOlyzer\backend\layer_04_skills"

def ensure_dir(d):
    os.makedirs(os.path.join(base_dir, d), exist_ok=True)

directories = [
    "registry_views",
    "policies",
    "internal",
    "tests",
    "skills/keyword_research/references",
    "skills/keyword_research/assets",
    "skills/keyword_research/scripts",
    "skills/keyword_research/evals",
    # I'll just scaffold one skill fully, the others just basic dirs
    "skills/topic_discovery",
    "skills/content_brief_building",
    "skills/source_backed_research",
    "skills/seo_content_planning",
    "skills/page_content_analysis",
    "skills/article_planning",
    "skills/guarded_drafting",
    "skills/content_optimization",
    "skills/content_repurposing",
    "skills/tone_memory_guidance",
    "skills/citation_source_safety",
    "skills/content_quality_gates"
]

for d in directories:
    ensure_dir(d)

# I will parse the spec directly in Go or just write the Go files based on the structs in the spec.
print("Directories scaffolded.")
