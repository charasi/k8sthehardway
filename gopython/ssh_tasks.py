import os
import paramiko


def process_ssh_task(hostname: str, key_path: str, username: str, cmd: str):
    status = ""
    ssh = paramiko.SSHClient()
    ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())

    key_path = os.path.expanduser(key_path)
    port = 22
    try:
        ssh.connect(hostname, port=port, username=username, key_filename=key_path)
        stdin, stdout, stderr = ssh.exec_command(cmd)
        # Wait for the command to complete
        stdout.channel.recv_exit_status()

        output = stdout.read().decode()
        error = stderr.read().decode()

        if output:
            status = "success"
        if error:
            status = "failure"

    except Exception as e:
        status = str(e)

    finally:
        ssh.close()

    return status
