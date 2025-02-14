# Output the ID of the VPC network
output "network_id" {
  # This will display the ID of the VPC network created by the `google_compute_network` resource.
  # The `id` is a unique identifier for the VPC network within Google Cloud.
  # This output can be used in other Terraform configurations or simply for informational purposes.
  value = google_compute_network.kthw_network.id
}

# Output the ID of the subnetwork
output "subnetwork_id" {
  # This will display the ID of the subnetwork created by the `google_compute_subnetwork` resource.
  # The `id` is a unique identifier for the subnetwork within Google Cloud.
  # This output can be used in other Terraform configurations or simply for informational purposes.
  value = google_compute_subnetwork.kthw-network-us-west1-subnet.id
}

# Output the IP address of the 'master_node'
output "master_node_ip" {
  value = google_compute_address.master_node.address
}

output "static_ip_address" {
  value = google_compute_address.kubernetes_the_hard_way.address
}

/**
output "seagram_private_subnet" {
  value = google_compute_global_address.seagram_private_ip_range.name
}
 */


