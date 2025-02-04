# Reference the state of the first Terraform configuration
data "terraform_remote_state" "vpc_compute" {
  backend = "local"  # Or you could use a remote backend like S3, GCS, etc.

  config = {
    path = "../ic/terraform.tfstate"  # Path to your first Terraform state file
  }
}

# Create a Proxy-Only Subnet
resource "google_compute_subnetwork" "proxy_only_subnet" {
  name          = "proxy-only-subnet"
  region        = "us-west1" # Change to your region
  network       = data.terraform_remote_state.vpc_compute.outputs.vpc_network_id
  ip_cidr_range = "10.240.30.0/24"  # Adjust CIDR block as per your requirements
  private_ip_google_access = true # For Private Google Access
}

# Health Checks
resource "google_compute_http_health_check" "controller_health_check" {
  name               = "controller-health-check"
  request_path       = "/livez"
  port               = 80
  check_interval_sec = 5
  timeout_sec        = 5
  healthy_threshold  = 2
  unhealthy_threshold = 2
}

resource "google_compute_http_health_check" "worker_health_check" {
  name               = "worker-health-check"
  request_path       = "/livez"
  port               = 80
  check_interval_sec = 5
  timeout_sec        = 5
  healthy_threshold  = 2
  unhealthy_threshold = 2
}

resource "google_compute_instance_group" "controller_instance_group" {
  name        = "controller-instance-group"
  zone        = "us-west1-b"  # Adjust the zone where your controller instances reside
  instances = [
    data.terraform_remote_state.vpc_compute.outputs.controller_0_id,
    data.terraform_remote_state.vpc_compute.outputs.controller_1_id,
    data.terraform_remote_state.vpc_compute.outputs.controller_2_id
  ]
}

resource "google_compute_instance_group" "worker_instance_group" {
  name        = "worker-instance-group"
  zone        = "us-west1-b"  # Adjust the zone where your controller instances reside
  instances = [
    data.terraform_remote_state.vpc_compute.outputs.worker_0_id,
    data.terraform_remote_state.vpc_compute.outputs.worker_1_id,
    data.terraform_remote_state.vpc_compute.outputs.worker_2_id
  ]
}

# Backend Services
resource "google_compute_backend_service" "controller_backend_service" {
  name        = "controller-backend-service"
  protocol    = "HTTPS"
  backend {
    group = google_compute_instance_group.controller_instance_group.id
  }
  health_checks = [google_compute_http_health_check.controller_health_check.id]
}

resource "google_compute_backend_service" "worker_backend_service" {
  name        = "worker-backend-service"
  protocol    = "HTTPS"
  backend {
    group = google_compute_instance_group.worker_instance_group.id
  }
  health_checks = [google_compute_http_health_check.worker_health_check.id]
}

# URL Map
resource "google_compute_url_map" "k8_url_map" {
  name = "k8-url-map"

  default_service = google_compute_backend_service.controller_backend_service.id

  host_rule {
    hosts = ["*"]
    path_matcher = "k8"
  }

  path_matcher {
    name = "k8"
    path_rule {
      paths   = ["/"]
      service = google_compute_backend_service.worker_backend_service.id
    }
  }
}


# Target HTTPS Proxy
resource "google_compute_target_https_proxy" "k8_target_proxy" {
  name      = "k8-target-proxy"
  url_map   = google_compute_url_map.k8_url_map.id
}

resource "google_compute_ssl_certificate" "k8_ssl_cert" {
  name        = "k8-ssl-cert"
  private_key = file("../k8_lb/load-balancer-key.pem")  # Your private key
  certificate = file("../k8_lb/load-balancer.pem")  # Your certificate
}


# Forwarding Rules
resource "google_compute_global_forwarding_rule" "controller_forwarding_rule" {
  name       = "controller-https-forwarding-rule"
  ip_address = data.terraform_remote_state.vpc_compute.outputs.static_ip_address
  target     = google_compute_target_https_proxy.k8_target_proxy.self_link
  port_range = "6443"
}

resource "google_compute_global_forwarding_rule" "worker_forwarding_rule" {
  name       = "worker-https-forwarding-rule"
  ip_address = data.terraform_remote_state.vpc_compute.outputs.static_ip_address
  target     = google_compute_target_https_proxy.k8_target_proxy.self_link
  port_range = "443"
}


