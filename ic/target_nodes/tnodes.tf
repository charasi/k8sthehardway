# Health Check Configuration (TCP)
resource "google_compute_region_health_check" "k8_health_check" {
  name               = "k8-tcp-health-check"
  timeout_sec         = 1
  check_interval_sec  = 1
  healthy_threshold   = 2
  unhealthy_threshold = 3

  http_health_check {
    port = 80  # Kubernetes API port
    request_path = "/livez"  # Adjust based on your application health check path
  }
}

# Controller Instance Group (with named ports)
resource "google_compute_instance_group" "controller_instance_group" {
  name        = "controller-instance-group"
  zone        = "us-west1-b"
  instances = [
    var.controller_0,
    var.controller_1,
    var.controller_2
  ]

  named_port {
    name = "k8-api-server"
    port = 6443
  }

}

# Regional Backend Service for Controller Instances (Regional Load Balancer)
resource "google_compute_region_backend_service" "controller_backend_service" {
  name        = "controller-backend-service"
  protocol    = "TCP"
  load_balancing_scheme = "EXTERNAL"

  backend {
    group = google_compute_instance_group.controller_instance_group.id
    balancing_mode = "CONNECTION"
  }

  health_checks = [google_compute_region_health_check.k8_health_check.id]  # Using the health check for port 6443
}

# Global Forwarding Rule for Controller API (6443)
resource "google_compute_forwarding_rule" "controller_forwarding_rule" {
  name       = "controller-tcp-forwarding-rule"
  ip_address = var.ip_address  # Static IP address (ensure it's created or specified)
  backend_service = google_compute_region_backend_service.controller_backend_service.id
  port_range = "6443"
  region     = "us-west1"  # Specify the region
  ip_protocol = "TCP"
}
