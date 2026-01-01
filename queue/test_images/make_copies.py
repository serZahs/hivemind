import shutil

def make_copies(image_path, n):
    base_name = "test"
    for i in range(1, n + 1):
        new_image_path = f"{base_name}{i}.jpg"
        shutil.copy(image_path, new_image_path)
    print(f"Created {n} copies of {image_path}.")

# Usage example:
make_copies("test0.jpg", 30)