variable "machine_type" {
 description = "The machine type to create"
 type        = string
 default     = "e2-medium"
}

variable "boot_disk_image" {
 description = "The boot disk for the instance"
 type        = string
 default     = "ubuntu-os-cloud/ubuntu-2204-lts"
}

variable "boot_disk_type" {
 description = "The GCE disk type"
 type        = string
 default     = "pd-balanced"
}

variable "network_name" {
 description = "name of network"
 type        = string
}

variable "subnetwork_name" {
 description = "name of subnetwork"
 type        = string
}

variable "network_ip" {
 description = "name of subnetwork"
 type        = string
 default = "10.240.0.1"
}

variable "zone" {
 description = "The zone that the machine should be created in"
 type        = string
 default     = "us-west1-b"
}

variable "master_node_ext_ip" {
  type = string
}

variable "scopes" {
  description = "List of OAuth scopes for the service account"
  type        = list(string)
  default     = [
    "https://www.googleapis.com/auth/compute",              # compute-rw
    "https://www.googleapis.com/auth/devstorage.read_only",  # storage-ro
    "https://www.googleapis.com/auth/service.management",    # service-management
    "https://www.googleapis.com/auth/servicecontrol",        # service-control
    "https://www.googleapis.com/auth/logging.write",         # logging-write
    "https://www.googleapis.com/auth/monitoring",             # monitoring
    "https://www.googleapis.com/auth/cloud-platform"
  ]
}

variable "bucket_private_key" {
  type = string
}

variable "bucket_name" {
  type = string
}

variable "kthw_private_key" {
  type = string
  sensitive = true
}

variable "kthw_private_agent_key" {
  type = string
  sensitive = true
}

variable "static_ip_address" {
  type = string
}