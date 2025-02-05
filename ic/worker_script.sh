#!/bin/bash

# Ensure the script stops on error
set -e

# Update system package list
sudo apt-get update

# Install all required packages in a single command
sudo apt-get install -y maven openjdk-17-jre-headless python3 python3-pip software-properties-common

sudo apt install net-tools

# Add Ansible repository and install Ansible
sudo add-apt-repository --yes --update ppa:ansible/ansible
sudo apt-get install -y ansible-core

# install nginx
# sudo apt install nginx -y

# Script completion message
echo "All required packages (Maven, Java, Python3, Ansible) have been installed successfully."