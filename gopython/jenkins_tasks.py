import urllib

import jenkinsapi
from jenkinsapi.jenkins import Jenkins
import requests
from requests.auth import HTTPBasicAuth
import json


def create_nodes(jenkins: Jenkins, name):
    remote_fs = "/home/wisccourant/jenkins/"  # The file system root on the agent machine
    num_executors = 1  # Number of executors for the agent
    labels = "kthw"  # Labels for the node
    agent = jenkins.create_node(name=name, remote_fs=remote_fs, num_executors=num_executors, labels=labels)
    return agent


def install_plugins(jenkins: Jenkins, plugins_list):
    jenkins.install_plugins(plugin_list=plugins_list, restart=True)


def create_ssh_credential(jenkins_url, username, api_token, private_key, credential_id, ssh_username, description=""):
    with open(private_key, 'r') as key_file:
        private_key_content = key_file.read()

    # API URL for creating new credentials
    cred_api_url = f"{jenkins_url}credentials/store/system/domain/_/createCredentials"

    """
    # JSON Payload to create SSH Username and Private Key Credential
    # Form the payload as a dictionary to be URL-encoded
    # Prepare the form data, including the private key and other fields
    form_data = {
        "": "0",  # Placeholder for the API payload, if necessary
        "credentials.scope": "GLOBAL",  # Scope can be GLOBAL or SYSTEM
        "credentials.id": credential_id,  # Unique ID for the credential
        "credentials.username": ssh_username,  # SSH username
        "credentials.privateKeySource.stapler-class": "hudson.plugins.sshslaves.impl.BasicSSHUserPrivateKey$DirectEntryPrivateKeySource",
        "credentials.privateKeySource.privateKey": private_key_content,  # The SSH private key content
        "credentials.description": description,  # Optional description
    }

    # Headers for the form submission
    headers = {
        "Content-Type": "application/x-www-form-urlencoded"
    }

    headers = {
        "Content-Type": "application/xml"
    }
    """


    # Create the form data (the payload to be sent as form submission)
    form_data = {
        "json": json.dumps({
            "credentials": {
                "scope": "GLOBAL",
                "id": credential_id,
                "username": ssh_username,
                "privateKeySource": {
                    "class": "hudson.plugins.sshslaves.impl.BasicSSHUserPrivateKey$DirectEntryPrivateKeySource",
                    "privateKey": private_key
                },
                "description": description
            }
        })
    }

    # Define headers for the form submission
    headers = {
        "Content-Type": "application/x-www-form-urlencoded"
    }

    # Send the POST request to create the SSH credential
    response = requests.post(cred_api_url, data=form_data, auth=HTTPBasicAuth(username, api_token), headers=headers)


    # Return the response object for further handling or error checking
    return response


def create_jobs(jenkins: Jenkins, job_name: str, job_cfg):
    job = jenkins.create_job(job_name, job_cfg)
    return job
