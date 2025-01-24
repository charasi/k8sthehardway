#!/bin/bash

# Ensure the script stops on error
set -e

# Update the system package list
sudo apt-get update

# Install required packages in a single command
sudo apt-get install -y maven openjdk-17-jre-headless python3 python3-pip software-properties-common

sudo apt install net-tools

# Add Jenkins repository and key
sudo wget -O /usr/share/keyrings/jenkins-keyring.asc \
  https://pkg.jenkins.io/debian-stable/jenkins.io-2023.key

# Add Jenkins to the sources list
echo "deb [signed-by=/usr/share/keyrings/jenkins-keyring.asc] \
  https://pkg.jenkins.io/debian-stable binary/" | sudo tee \
  /etc/apt/sources.list.d/jenkins.list > /dev/null

# Update package list and install Jenkins
sudo apt-get update
sudo apt-get install -y jenkins

# Add Ansible repository and install it
sudo add-apt-repository --yes --update ppa:ansible/ansible
sudo apt-get install -y ansible

# install nginx
sudo apt install nginx -y

# Ensure the .ssh directory exists
mkdir -p ~/.ssh

# Download the private key from GCS bucket
sudo gsutil cp gs://kthw-misc/private_key.pem ~/.ssh/id_rsa

# Fix permissions for the private key
sudo chmod 600 ~/.ssh/id_rsa