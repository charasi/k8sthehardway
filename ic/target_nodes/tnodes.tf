# Health check for the Load Balancer (on port 80)
resource "google_compute_http_health_check" "k8s_health_check" {
  name               = "k8s-http-health-check"
  request_path       = "/livez"  # Health check path for the reverse proxy
  port               = 80         # HTTP port for the health check
  check_interval_sec = 10         # Interval between health checks
  timeout_sec        = 5          # Timeout for each health check
  unhealthy_threshold = 3         # Number of failures before considering unhealthy
  healthy_threshold  = 2          # Number of successes before considering healthy
}

# Instance group containing the reverse proxy (NGINX or Go program)
resource "google_compute_instance_group" "proxy_instance_group" {
  name        = "proxy-instance-group"
  zone        = "us-west1-b"
  instances   = [var.controller_0, var.controller_1, var.controller_2]

  named_port {
    name = "http"
    port = 80  # Reverse proxy HTTP port
  }
}

# Backend service for the Kubernetes cluster (uses the reverse proxy)
resource "google_compute_backend_service" "k8_backend_service" {
  name        = "k8-backend-service"
  protocol    = "HTTPS"
  health_checks = [google_compute_http_health_check.k8s_health_check.self_link]

  backend {
    group = google_compute_instance_group.proxy_instance_group.self_link
  }
}

# Forwarding rule to direct traffic to the backend service (port 6443)
resource "google_compute_forwarding_rule" "kubernetes_forwarding_rule" {
  name        = "kubernetes-forwarding-rule"
  region      = "us-west1"
  port_range  = "6443"
  target      = google_compute_backend_service.k8_backend_service.self_link
  ip_address  = var.ip_address
}
