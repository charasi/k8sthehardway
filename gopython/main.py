import os

import requests
import xml.etree.ElementTree as ET
from jenkinsapi.jenkins import Jenkins

from gopython import jenkins_tasks, gcp_tasks, read_tasks, ssh_tasks


def main():
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/private_key.pem", ".")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/external_ip.txt", ".")

    ip_addr = read_tasks.get_ip_address('external_ip.txt')

    key_path = os.path.expanduser('~/.ssh/kthw_key')
    username = 'wisccourant'

    command = 'sudo gsutil cp gs://kthw-misc/private_key.pem ~/.ssh/private_key.pem'
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)
    command = 'chmod 600 ~/.ssh/private_key.pem'
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)
    command = 'ssh-keyscan -H 10.240.0.60 >> ~/.ssh/known_hosts'
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    # Connect to the Jenkins server
    jenkins_url = 'http://' + ip_addr + ':8080/'

    jenkins_user = 'kube'
    password = '11d466ed80ba909682b652797477def2e0'
    jenkins = Jenkins(jenkins_url, username=jenkins_user, password=password)

    # Define the node parameters
    node_name = "kthw-agent"
    agent = jenkins_tasks.create_nodes(jenkins, node_name)

    api_url = f"{jenkins_url}/computer/{node_name}/slave-agent.jnlp"

    # Make the request to Jenkins using basic auth
    response = requests.get(api_url, auth=(jenkins_user, password))
    root = ET.fromstring(response.text)

    # Extract the first <argument> (which contains the secret)
    secret = root.find(".//argument").text

    agent_work_dir = "/home/wisccourant/jenkins/agent"

    gsutil_command = 'sudo gsutil cp gs://kthw-misc/private_agent_key.pem ~/.ssh/private_agent_key.pem'
    host_b = "10.240.0.60"
    key_path_b = os.path.expanduser('/home/wisccourant/.ssh/private_key.pem')

    command = f"ssh -i {key_path_b} {username}@{host_b} {gsutil_command}"
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    ip_addresses = [
        "10.240.0.10",
        "10.240.0.11",
        "10.240.0.12",
        "10.240.0.20",
        "10.240.0.21",
        "10.240.0.22"
    ]

    # Define the command to run on the remote machine to add Machine B's SSH key to known_hosts
    gsutil_command = 'ssh-keyscan -H {0} >> ~/.ssh/known_hosts'
    # Correct SSH command to execute on Machine A, ensuring proper quoting
    # For each IP, replace the placeholder in gsutil_command
    for ip_addr in ip_addresses:
        # Format the gsutil_command with the current IP address
        command = f"ssh -i {key_path_b} {username}@{host_b} '{gsutil_command.format(ip_addr)}'"
        # Execute the command for the current IP
        status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    curl_command = f"curl -sO {jenkins_url}jnlpJars/agent.jar"
    java_command = f"java -jar agent.jar -url {jenkins_url} -secret {secret} -name {node_name} -webSocket -workDir \"{agent_work_dir}\""
    command = f"ssh -i {key_path_b} {username}@{host_b} {curl_command}"
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)
    command = f"ssh -i {key_path_b} {username}@{host_b} {java_command}"
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)
    print("hi")


if __name__ == "__main__":
    main()
