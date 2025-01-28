import os
import paramiko


def process_ssh_task(hostname: str, key_path: str, username: str, cmd: str):
    status = ""
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    key_path = os.path.expanduser(key_path)
    port = 22

    try:
        # Attempt to connect to the remote server with a timeout
        ssh.connect(hostname, port=port, username=username, key_filename=key_path, timeout=10)

        # Execute the command with a timeout for the execution
        stdin, stdout, stderr = ssh.exec_command(cmd, timeout=30)

        # Non-blocking reads to avoid deadlocks
        stdout_data = stdout.read()
        stderr_data = stderr.read()

        exit_status = stdout.channel.recv_exit_status()

        output = stdout_data.decode() if stdout_data else ""
        error = stderr_data.decode() if stderr_data else ""

        # Determine the status based on the output/error
        if exit_status == 0:
            status = "success" if output else "failure"
        else:
            status = f"failure (exit code {exit_status}): {error}"

    except Exception as e:
        status = f"Error: {str(e)}"

    finally:
        ssh.close()

    return status
