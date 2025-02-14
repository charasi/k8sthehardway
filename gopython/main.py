import os
import time

import requests
import xml.etree.ElementTree as ET
from jenkinsapi.jenkins import Jenkins

from gopython import jenkins_tasks, gcp_tasks, read_tasks, ssh_tasks


# Main function to orchestrate tasks
def main():
    """
    This script automates the setup of a Jenkins agent on a remote machine by performing
    several steps such as downloading necessary files from GCP, setting up SSH keys,
    configuring Jenkins plugins, creating a Jenkins agent node, and starting the agent.
    """

    # Step 1: Download necessary files from GCP Storage using gsutil
    # The gsutil command copies files from a Google Cloud Storage bucket to the local directory.
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/private_key.pem", ".")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/external_ip.txt", ".")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/private_agent_key.pem",
                           ".")

    # Step 2: Read the external IP address from the file `external_ip.txt`
    ip_addr = read_tasks.get_ip_address('external_ip.txt')

    # Step 3: Define paths and SSH credentials
    # The following are variables to define the SSH key path, username, and commands to run on the remote machine.
    key_path = os.path.expanduser('~/.ssh/kthw_key')  # Path to the private SSH key
    username = 'wisccourant'  # SSH username on the remote machine

    # Step 4: Set up the SSH environment on the remote machine
    # These commands are executed on the remote machine using SSH. They handle file copying,
    # changing file permissions, and adding the remote machine's IP to known_hosts.
    command = 'sudo gsutil cp gs://kthw-misc/private_key.pem ~/.ssh/private_key.pem'
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    command = 'chmod 600 ~/.ssh/private_key.pem'  # Set secure permissions for the private key
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    command = 'ssh-keyscan -H 10.240.0.60 >> ~/.ssh/known_hosts'  # Add machine's IP to known_hosts
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    # Step 5: Connect to the Jenkins server
    # Construct the Jenkins server URL using the IP address and port 8080
    jenkins_url = 'http://' + ip_addr + ':8080/'

    # Jenkins user credentials for authentication
    jenkins_user = 'kube'
    password = '11913f3d311d73ba9918026c3ad6dd3fad'

    # Initialize the Jenkins API client
    jenkins = Jenkins(jenkins_url, username=jenkins_user, password=password)

    # Step 6: Install Jenkins plugins
    # Define a list of plugins to install on the Jenkins server
    plugins_list = ['ansible@500.v7564a_db_8feec', 'blueocean@1.27.16']
    jenkins_tasks.install_plugins(jenkins, plugins_list)  # Install the plugins via Jenkins API

    # Step 7: Create a new Jenkins agent node
    node_name = "kthw-agent"  # Define the name of the agent node
    agent = jenkins_tasks.create_nodes(jenkins, node_name)  # Create the agent on Jenkins

    # Step 8: Fetch the Jenkins agent secret
    # This request is used to get the secret needed to authenticate the Jenkins agent
    api_url = f"{jenkins_url}/computer/{node_name}/slave-agent.jnlp"
    response = requests.get(api_url, auth=(jenkins_user, password))
    root = ET.fromstring(response.text)  # Parse the XML response from Jenkins
    secret = root.find(".//argument").text  # Extract the agent secret

    # Step 9: Define the working directory for the agent on the remote machine
    agent_work_dir = "/home/wisccourant/jenkins/"

    # Step 10: Transfer the private agent key to the remote machine
    gsutil_command = 'sudo gsutil cp gs://kthw-misc/private_agent_key.pem ~/.ssh/private_agent_key.pem'
    host_b = "10.240.0.60"  # IP address of the remote machine
    key_path_b = os.path.expanduser('/home/wisccourant/.ssh/private_key.pem')  # Path to the SSH key for host_b

    # Execute the SSH command to transfer the agent key
    command = f"ssh -i {key_path_b} {username}@{host_b} '{gsutil_command}'"
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    # Step 11: Add multiple machine IPs to known_hosts
    # This block iterates through a list of IP addresses and adds them to the SSH known_hosts file
    ip_addresses = [
        "10.240.0.10", "10.240.0.11", "10.240.0.12",
        "10.240.0.20", "10.240.0.21", "10.240.0.22", "10.240.0.70"
    ]

    # Define the command to add an IP address to the known_hosts file
    gsutil_command = 'ssh-keyscan -H {0} >> ~/.ssh/known_hosts'
    for addr in ip_addresses:
        command = f"ssh -i {key_path_b} {username}@{host_b} '{gsutil_command.format(addr)}'"
        status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    #jenkins_tasks.create_ssh_credential(jenkins_url, jenkins_user, password, "../gopython/private_key.pem", "ad-agent", username, "")

    # Step 12: Download the Jenkins agent JAR file
    curl_command = f"curl -sO {jenkins_url}jnlpJars/agent.jar"  # Command to download the agent JAR file
    java_command = f"java -jar agent.jar -url {jenkins_url} -secret {secret} -name {node_name} -webSocket -workDir \"{agent_work_dir}\""  # Command to start the Jenkins agent

    # Execute the curl command to download the agent
    command = f"ssh -i {key_path_b} {username}@{host_b} '{curl_command}'"
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    # Execute the java command to start the agent on the remote machine
    command = f"ssh -i {key_path_b} {username}@{host_b} '{java_command}'"
    status = ssh_tasks.process_ssh_task(ip_addr, key_path, username, command)

    # Read the config.xml file
    with open("./config.xml", 'r') as file:
        config_xml = file.read()

    job = jenkins_tasks.create_jobs(jenkins, "install-k8", config_xml)

    jenkins.build_job(job.name)

    while job.is_running():
        time.sleep(10)

    if job.get_last_build().get_status() != 'SUCCESS':
        exit(1)
    
    # Read the config.xml file
    with open("./certificates.xml", 'r') as file:
        config_xml = file.read()

    job = jenkins_tasks.create_jobs(jenkins, "create-certificates", config_xml)

    jenkins.build_job(job.name)

    while job.is_running():
        time.sleep(10)

    if job.get_last_build().get_status() != 'SUCCESS':
        exit(1)
    
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/admin-key.pem", "/home/charasi/cmu/devops/k8sthehardway/kubelet")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/admin.pem", "/home/charasi/cmu/devops/k8sthehardway/kubelet")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/ca-key.pem", "/home/charasi/cmu/devops/k8sthehardway/kubelet")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/ca.pem", "/home/charasi/cmu/devops/k8sthehardway/kubelet")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/load-balancer-key.pem", "/home/charasi/cmu/devops/k8sthehardway/k8_lb")
    gcp_tasks.gcp_cp_tasks('/home/charasi/google-cloud-sdk/bin/gsutil', "cp", "gs://kthw-misc/load-balancer.pem", "/home/charasi/cmu/devops/k8sthehardway/k8_lb")

    # Read the config.xml file
    with open("./etcd.xml", 'r') as file:
        config_xml = file.read()

    job = jenkins_tasks.create_jobs(jenkins, "create-etcd", config_xml)

    jenkins.build_job(job.name)

    while job.is_running():
        time.sleep(10)

    if job.get_last_build().get_status() != 'SUCCESS':
        exit(1)

    # Read the config.xml file
    with open("./controllers.xml", 'r') as file:
        config_xml = file.read()

    job = jenkins_tasks.create_jobs(jenkins, "create-controllers", config_xml)

    jenkins.build_job(job.name)

    while job.is_running():
        time.sleep(10)

    if job.get_last_build().get_status() != 'SUCCESS':
        exit(1)

    # Read the config.xml file
    with open("./rbac.xml", 'r') as file:
        config_xml = file.read()

    job = jenkins_tasks.create_jobs(jenkins, "create-rbac", config_xml)

    jenkins.build_job(job.name)

    while job.is_running():
        time.sleep(10)

    if job.get_last_build().get_status() != 'SUCCESS':
        exit(1)

    # Read the config.xml file
    with open("./workers.xml", 'r') as file:
        config_xml = file.read()

    job = jenkins_tasks.create_jobs(jenkins, "create-workers", config_xml)

    jenkins.build_job(job.name)

    while job.is_running():
        time.sleep(10)

    if job.get_last_build().get_status() != 'SUCCESS':
        exit(1)


# Entry point for the script
if __name__ == "__main__":
    main()
