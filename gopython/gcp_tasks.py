import subprocess

from google.cloud import storage
from google.oauth2 import service_account


def gcp_cp_tasks(gsutil_path, gsutil_cmd, source, destination):
    # Run a gsutil command
    command = [gsutil_path, gsutil_cmd, source, destination]
    result = subprocess.run(command, capture_output=True, text=True)

    # Check the result and print output
    if result.returncode == 0:
        print("Files in the bucket:")
        print(result.stdout)
    else:
        print(f"Error: {result.stderr}")
