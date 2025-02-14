# Define a Google Cloud firewall rule for allowing internal traffic in the network
resource "google_compute_firewall" "kthw_network_us_west1_subnet_firewall_allow_internal" {
  
  # The name of the firewall rule. This should be descriptive, reflecting the region and the purpose (internal traffic allowed).
  name    = "kthw-network-us-west1-subnet-firewall-allow-internal"
  
  # The VPC network where the firewall rule will be applied.
  # The `network_id` variable should reference the ID of the VPC network.
  network = var.network_id

  # Allow ICMP traffic (ping) for network diagnostics and basic connectivity checks.
  allow {
    protocol = "icmp"  # Allow ICMP traffic, which is typically used for network diagnostics (ping)
  }

  # Allow all TCP traffic within the source IP ranges. This means all TCP ports are open for internal communication.
  allow {
    protocol = "tcp"  # Allow all TCP traffic
  }

  # Allow all UDP traffic within the source IP ranges. This opens all UDP ports for internal communication.
  allow {
    protocol = "udp"  # Allow all UDP traffic
  }

  # Define the allowed source IP ranges from which traffic can access the network.
  # In this case, only IPs from `10.240.0.0/24` and `10.200.0.0/16` subnets are allowed to send traffic.
  source_ranges = ["10.240.0.0/24", "10.200.0.0/16"]
}

# Define a Google Cloud firewall rule for allowing external access to the specified network and subnetwork
resource "google_compute_firewall" "kthw_network_us_west1_subnet_firewall_allow_external" {
  
  # The name of the firewall rule. The name should reflect the purpose (allow external access) and region (us-west1).
  # This rule is intended to allow external access to the network on specific ports.
  name    = "kthw-network-us-west1-subnet-firewall-allow-external"
  
  # The VPC network the firewall rule applies to. `network_id` is a variable that should be the ID of your VPC network.
  network = var.network_id

  # Define the allowed protocols and ports for the firewall rule.

  # Allow ICMP traffic (ping), useful for basic network diagnostics.
  allow {
    protocol = "icmp"  # Allow ICMP (ping)
  }

  /**
  # Allow incoming TCP traffic on port 22 (SSH) from any source IP range.
  # This is typically used for allowing SSH access to instances for management.
  allow {
    protocol = "tcp"
    ports    = ["22"]  # Allow SSH traffic on port 22
  }

   */

  # Allow incoming TCP traffic on port 6443 (typically used by Kubernetes API server).
  allow {
    protocol = "tcp"
    ports    = ["22", "80", "8080", "443", "3306", "30637", "31420", "6443", "8443"]  # Allow Kubernetes API server traffic (port 6443)
  }

  # Define the source IP ranges that are allowed to access the resources within the network.
  # `0.0.0.0/0` means that this rule allows traffic from any IP address in the world.
  # For production environments, it's strongly recommended to restrict this to a smaller range if possible.
  source_ranges = ["0.0.0.0/0"]
}

/**
resource "google_compute_firewall" "kthw-allow-egress" {
  name    = "kthw-fault-allow-egress"
  network = var.network_id

  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }

  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }

  direction = "EGRESS"
  priority  = 1000
}
*/



