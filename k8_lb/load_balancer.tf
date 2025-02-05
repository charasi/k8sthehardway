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

resource "google_compute_health_check" "worker_health_check" {
  name               = "worker-health-check"
  timeout_sec         = 1
  check_interval_sec  = 1
  healthy_threshold   = 4
  unhealthy_threshold = 5

  http_health_check {
    port_name          = "worker-health-check-port"
    port = 80
    request_path       = "/livez"
  }
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

resource "google_compute_region_backend_service" "worker_backend_service" {
  name        = "worker-backend-service"
  protocol    = "HTTPS"
  backend {
    group = google_compute_instance_group.worker_instance_group.id
  }
  health_checks = [google_compute_health_check.worker_health_check.id]
}

# URL Map
resource "google_compute_region_url_map" "k8_url_map" {
  name = "k8-url-map"

  default_service = google_compute_region_backend_service.worker_backend_service.id
}


# Target HTTPS Proxy
resource "google_compute_region_target_https_proxy" "k8_target_proxy" {
  name      = "k8-target-proxy"
  url_map   = google_compute_region_url_map.k8_url_map.id
  ssl_certificates = [google_compute_ssl_certificate.k8_ssl_cert.id]
}

resource "google_compute_ssl_certificate" "k8_ssl_cert" {
  name        = "k8-ssl-cert"
  private_key = file("../k8_lb/load-balancer-key.pem")  # Your private key
  certificate = file("../k8_lb/load-balancer.pem")  # Your certificate
}

resource "google_compute_forwarding_rule" "worker_forwarding_rule" {
  name       = "worker-https-forwarding-rule"
  ip_address = data.terraform_remote_state.vpc_compute.outputs.static_ip_address
  target     = google_compute_region_target_https_proxy.k8_target_proxy.self_link
  port_range = "443"
  region     = "us-west1"
}


