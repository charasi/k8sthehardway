import jenkinsapi
from jenkinsapi.jenkins import Jenkins


def create_nodes(jenkins: Jenkins, name):
    remote_fs = "/home/jenkins/agent"  # The file system root on the agent machine
    num_executors = 1  # Number of executors for the agent
    labels = "kthw"  # Labels for the node
    agent = jenkins.create_node(name=name, remote_fs=remote_fs, num_executors=num_executors, labels=labels)
    return agent
