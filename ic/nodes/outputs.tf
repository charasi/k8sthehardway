output "private_key_pem" {
  value     = tls_private_key.kthw_ssh.private_key_pem
  sensitive = true
}

/**
output "public_key_openssh" {
  value = tls_private_key.kthw_ssh.public_key_openssh
}
 */

output "instance_external_ip" {
  value = google_compute_instance.master.network_interface[0].access_config[0].nat_ip
}

output "private_agent_key_pem" {
  value     = tls_private_key.kthw_ssh_agent.private_key_pem
  sensitive = true
}

/**
output "public_agent_key_openssh" {
  value = tls_private_key.kthw_ssh_agent.public_key_openssh
}
 */