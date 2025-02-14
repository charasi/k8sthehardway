# Reference the state of the first Terraform configuration
data "terraform_remote_state" "vpc_compute" {
  backend = "local"
  config = {
    path = "../ic/terraform.tfstate"
  }
}

# Create a Proxy-Only Subnet
resource "google_compute_subnetwork" "k8s_https_lb_target_proxy" {
  name          = "k8s-https-lb-target-proxy"
  region        = "us-west1"
  network       = data.terraform_remote_state.vpc_compute.outputs.vpc_network_id
  ip_cidr_range = "10.240.30.0/24"
  #gateway_address = "10.240.30.1"  # Gateway address for routing traffic

  # The "private_ip_google_access" is optional but useful for proxy-based services
  # private_ip_google_access = true
  role          = "ACTIVE"
  # Set the "purpose" to "PURPOSE_REGIONAL_MANAGED_PROXY" for load balancing
  purpose = "REGIONAL_MANAGED_PROXY"
}

# Health Check
resource "google_compute_region_health_check" "worker_health_check" {
  name                = "worker-health-check"
  timeout_sec         = 1
  check_interval_sec  = 1
  healthy_threshold   = 4
  unhealthy_threshold = 5

  http_health_check {
    port_name    = "worker-health-check-port"
    port         = 30637
    request_path = "/healthz"
  }
}

# Instance Group
resource "google_compute_instance_group" "worker_instance_group" {
  name = "worker-instance-group"
  zone = "us-west1-b"
  instances = [
    data.terraform_remote_state.vpc_compute.outputs.worker_0_id,
    data.terraform_remote_state.vpc_compute.outputs.worker_1_id,
    data.terraform_remote_state.vpc_compute.outputs.worker_2_id
  ]

  /**
  named_port {
    name = "ingress"
    port = 30637
  }

   */

  named_port {
    name = "ingress-http"
    port = 30637
  }
}

/**
# SSL Certificate
resource "google_compute_region_ssl_certificate" "k8_ssl_cert" {
  region      = "us-west1"
  name        = "k8-ssl-cert"
  private_key = file("../k8_lb/load-balancer-key.pem") # Your private key
  certificate = file("../k8_lb/load-balancer.pem")     # Your certificate
}

 */

# Regional Backend Service
resource "google_compute_region_backend_service" "worker_backend_service" {
  name                  = "worker-backend-service"
  protocol              = "HTTP"
  port_name             = "http"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  region                = "us-west1"
  backend {
    group           = google_compute_instance_group.worker_instance_group.id
    balancing_mode  = "UTILIZATION"
    max_utilization = 0.8
    capacity_scaler = 1.0
  }

  health_checks = [google_compute_region_health_check.worker_health_check.id]
}

# Regional URL Map
resource "google_compute_region_url_map" "k8_url_map" {
  name            = "k8-url-map"
  default_service = google_compute_region_backend_service.worker_backend_service.id
}

/**
# Regional Target HTTPS Proxy
resource "google_compute_target_https_proxy" "k8_target_proxy" {
  name = "k8-target-proxy"
  //region           = "us-west1"  # Specify your region here
  url_map          = google_compute_region_url_map.k8_url_map.id
  ssl_certificates = [google_compute_region_ssl_certificate.k8_ssl_cert.id]
}

 */

# Regional Target HTTPS Proxy
resource "google_compute_region_target_http_proxy" "k8_target_proxy" {
  name             = "k8-target-proxy"
  region           = "us-west1"
  url_map          = google_compute_region_url_map.k8_url_map.id
  //ssl_certificates = [google_compute_region_ssl_certificate.k8_ssl_cert.id]

  /**
  depends_on = [
    google_compute_region_url_map.k8_url_map  # Ensure URL map is created before the proxy
  ]

   */
}

# Regional Forwarding Rule
resource "google_compute_forwarding_rule" "worker_forwarding_rule" {
  name                  = "worker-https-forwarding-rule"
  network               = data.terraform_remote_state.vpc_compute.outputs.vpc_network_id
  region                = "us-west1" # Specify your region here
  ip_address            = data.terraform_remote_state.vpc_compute.outputs.static_ip_address
  target                = google_compute_region_target_http_proxy.k8_target_proxy.self_link
  port_range            = "80"
  load_balancing_scheme = "EXTERNAL_MANAGED"
  network_tier          = "PREMIUM"
}


