import os
import shutil

base_dir = r"C:\Users\Muham\AEOlyzer\backend"
layer4_dir = os.path.join(base_dir, "layer_04_skills")
internal_skills_dir = os.path.join(base_dir, "internal", "skills")

if os.path.exists(layer4_dir):
    for item in os.listdir(layer4_dir):
        src = os.path.join(layer4_dir, item)
        dst = os.path.join(internal_skills_dir, item)
        if os.path.exists(dst):
            if os.path.isdir(src):
                # merge dirs
                for sub_item in os.listdir(src):
                    shutil.move(os.path.join(src, sub_item), os.path.join(dst, sub_item))
                shutil.rmtree(src)
            else:
                shutil.move(src, dst)
        else:
            shutil.move(src, dst)
    
    # Remove old dir
    if os.path.exists(layer4_dir) and not os.listdir(layer4_dir):
        shutil.rmtree(layer4_dir)

print("Layer 4 migrated entirely into internal/skills")
