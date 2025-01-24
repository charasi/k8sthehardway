resource "google_compute_target_pool" "k8_target_pool" {
  name    = "k8-target-pool"
  region  = "us-west1"  # specify your region
  health_checks = [google_compute_http_health_check.k8s_health_check.name]
  session_affinity = "NONE"
}

resource "google_compute_http_health_check" "k8s_health_check" {
  name               = "k8s-health-check"
  request_path       = "/livez"  # Path for Kubernetes API server health check
  port               = 80        # Port for Kubernetes API server
  check_interval_sec = 10          # Interval between health checks
  timeout_sec        = 5           # Timeout for each health check
  unhealthy_threshold = 3          # Number of failures before considering unhealthy
  healthy_threshold  = 2           # Number of successes before considering healthy
}


resource "google_compute_instance_group" "k8s_instance_group" {
  name        = "k8s-instance-group"
  zone               = "us-west1-b"
  instances   = [var.controller_0, var.controller_1, var.controller_2]
  named_port {
    name = "https"
    port = 443  # For HTTPS traffic
  }

  named_port {
    name = "kubernetes-api"
    port = 6443  # For Kubernetes API server
  }

  named_port {
    name = "http"
    port = 80  # For HTTP traffic (e.g., ingress)
  }
}

resource "google_compute_forwarding_rule" "kubernetes_forwarding_rule" {
  name        = "kubernetes-forwarding-rule"
  region      = "us-west1"
  port_range  = "6443"
  target      = google_compute_target_pool.k8_target_pool.self_link
  ip_address  = var.ip_address
}