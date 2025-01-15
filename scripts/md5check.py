import os
import hashlib

def calculate_md5(file_path):
    """Calculate the MD5 hash of a file."""
    hash_md5 = hashlib.md5()
    with open(file_path, 'rb') as f:
        for chunk in iter(lambda: f.read(4096), b""):
            hash_md5.update(chunk)
    return hash_md5.hexdigest()

def compare_folders(folder1, folder2):
    """Compare files in two folders by their MD5 hashes."""
    
    # Get a list of all files in both folders (including subdirectories)
    files_folder1 = {os.path.relpath(os.path.join(root, file), folder1): os.path.join(root, file)
                     for root, _, files in os.walk(folder1)
                     for file in files}

    files_folder2 = {os.path.relpath(os.path.join(root, file), folder2): os.path.join(root, file)
                     for root, _, files in os.walk(folder2)
                     for file in files}
    
    # Loop over files in folder1
    for relative_path, file1_path in files_folder1.items():
        file2_path = files_folder2.get(relative_path)
        
        if file2_path and os.path.exists(file2_path):
            # Compare MD5 hashes of the files in both folders
            md5_file1 = calculate_md5(file1_path)
            md5_file2 = calculate_md5(file2_path)
            
            if md5_file1 == md5_file2:
                print(f"Files match: {relative_path}")
            else:
                print(f"Files differ: {relative_path}")
        else:
            print(f"File missing in second folder: {relative_path}")

    # Check for files in folder2 that don't exist in folder1
    for relative_path, file2_path in files_folder2.items():
        if relative_path not in files_folder1:
            print(f"File missing in first folder: {relative_path}")

if __name__ == "__main__":
    folder1 = f"{os.getenv('SDISK_HOME')}/users"
    folder2 = f"{os.getenv('SDISK_HOME')}/client_root"
    compare_folders(folder1, folder2)
