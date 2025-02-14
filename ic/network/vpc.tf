# Define a Google Cloud VPC Network
resource "google_compute_network" "kthw_network" {
  
  # The name of the network
  name = "kthw-network"
  
  # Disable automatic subnet creation, so we can define subnets manually
  # Setting this to `false` will prevent Google Cloud from creating subnets automatically in each region
  auto_create_subnetworks = false
  
  # Specify the routing mode of the network.
  # `REGIONAL` means routes are scoped to the region where the subnetwork exists
  routing_mode = "REGIONAL"
  
  # Delete default routes created when the network is first created
  # This can help to ensure that unwanted routes (e.g., for internet access) do not get created by default
  delete_default_routes_on_create = true
}

# Define a Google Cloud subnetwork within the specified VPC network
resource "google_compute_subnetwork" "kthw-network-us-west1-subnet" {

  # The name of the subnetwork.
  # It is a good practice to follow a naming convention that reflects both the VPC and region for clarity.
  # Example naming convention: <vpc-name>-<region>-subnet
  # Here, the subnetwork name reflects the VPC name (`kthw-network`) and the region (`us-west1`).
  name = "kthw-network-us-west1-subnet"  # The name assigned to the subnetwork

  # The network to which this subnetwork belongs.
  # You must reference the ID of the VPC network that the subnetwork will be created in.
  network = google_compute_network.kthw_network.id  # Link to the VPC network by ID

  # The region where this subnetwork will be created.
  # Google Cloud subnetworks are region-specific, so choose a region that is geographically closest
  # to where your resources will be deployed. This can help reduce latency.
  # In this case, the subnetwork is created in the `us-west1` region (Western United States).
  region = "us-west1"  # Specify the region for the subnetwork (ensure the region is correct for your needs)

  # The CIDR block for the subnetwork.
  # The CIDR block defines the IP address range available to resources in this subnetwork.
  # `10.240.0.0/24` provides 256 IP addresses (from `10.240.0.0` to `10.240.0.255`).
  # Make sure the CIDR range does not overlap with other subnetworks within the same VPC.
  ip_cidr_range = "10.240.0.0/24"  # The IP address range assigned to the subnetwork
}


# Define a static IP address resource for use in Google Cloud
resource "google_compute_address" "kubernetes_the_hard_way" {
  
  # The name of the static IP address resource. It should be descriptive and unique within the region.
  name    = "kubernetes-the-hard-way"
  
  # The region in which the static IP will be allocated. The IP will be specific to this region.
  region  = "us-west1"
}

# Define a static IP address resource for use in Google Cloud
resource "google_compute_address" "master_node" {
  
  # The name of the static IP address resource. It should be descriptive and unique within the region.
  name    = "master-node"
  
  # The region in which the static IP will be allocated. The IP will be specific to this region.
  region  = "us-west1"
}

resource "google_compute_router" "kthw_router" {
  name    = "kthw-router"
  region  = google_compute_subnetwork.kthw-network-us-west1-subnet.region
  network = google_compute_network.kthw_network.id

  bgp {
    asn = 64514
  }
}

resource "google_compute_route" "default_route" {
  name       = "default-route"
  network    = google_compute_network.kthw_network.id
  dest_range = "0.0.0.0/0"  # Default route for internet access
  //nex = google_compute_router.kthw_router.self_link
  //next_hop_gateway = google_compute_router.kthw_router.self_link
  //next_hop_ip = google_compute_router_nat.kthw_nat.nat_ip_allocate_option
  next_hop_gateway = "default-internet-gateway"
}

resource "google_compute_router_nat" "kthw_nat" {
  name                               = "kthw-nat"
  router                             = google_compute_router.kthw_router.name
  region                             = google_compute_router.kthw_router.region
  nat_ip_allocate_option             = "AUTO_ONLY"
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    enable = true
    filter = "ERRORS_ONLY"
  }
}

resource "google_compute_route" "worker_route_0" {
  name               = "kubernetes-route-10-200-0-0-24"
  network            = google_compute_network.kthw_network.name
  dest_range  = "10.200.0.0/24"
  next_hop_ip   = "10.240.0.20"  # Internal IP of worker-0
  depends_on = [google_compute_subnetwork.kthw-network-us-west1-subnet]
}

resource "google_compute_route" "worker_route_1" {
  name               = "kubernetes-route-10-200-1-0-24"
  network            = google_compute_network.kthw_network.name
  dest_range  = "10.200.1.0/24"
  next_hop_ip   = "10.240.0.21"  # Internal IP of worker-1
  depends_on = [google_compute_subnetwork.kthw-network-us-west1-subnet]
}

resource "google_compute_route" "worker_route_2" {
  name               = "kubernetes-route-10-200-2-0-24"
  network            = google_compute_network.kthw_network.name
  dest_range  = "10.200.2.0/24"
  next_hop_ip   = "10.240.0.22"  # Internal IP of worker-2
  depends_on = [google_compute_subnetwork.kthw-network-us-west1-subnet]
}

/**
# Step 1: Reserve a Private IP Range in the VPC
resource "google_compute_global_address" "seagram_private_ip_range" {
  name          = "seagram-private-ip-range"
  address_type  = "INTERNAL"
  network       = google_compute_network.kthw_network.id
  purpose       = "VPC_PEERING"
  prefix_length = 24
  address = "10.240.40.0"
}


# Step 2: Set up the Private Services Access connection for Cloud SQL
resource "google_service_networking_connection" "seagram_database_subnetwork" {
  network                 = google_compute_network.kthw_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.seagram_private_ip_range.name]
  depends_on              = [google_compute_global_address.seagram_private_ip_range]
  update_on_creation_fail = true
}

 */
