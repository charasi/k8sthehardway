resource "google_compute_http_health_check" "k8s_health_check" {
  name               = "k8s-health-check"
  request_path       = "/livez"  # Path for Kubernetes API server health check
  port               = 80        # Port for Kubernetes API server
  check_interval_sec = 10          # Interval between health checks
  timeout_sec        = 5           # Timeout for each health check
  unhealthy_threshold = 3          # Number of failures before considering unhealthy
  healthy_threshold  = 2           # Number of successes before considering healthy
}

resource "google_compute_health_check" "k8-tcp-health-check" {
  name = "tcp-health-check"

  timeout_sec        = 1
  check_interval_sec = 1

  tcp_health_check {
    port = "6443"
  }
}

resource "google_compute_forwarding_rule" "kubernetes_forwarding_rule" {
  name        = "kubernetes-forwarding-rule"
  region      = "us-west1"
  port_range  = "6443"
  target      = google_compute_backend_service.k8s_api_server.self_link
  ip_address  = var.ip_address
}

resource "google_compute_forwarding_rule" "k8s_service_forwarding_rule" {
  name        = "k8-service-forwarding-rule"
  region      = "us-west1"
  port_range  = "443"
  target      = google_compute_backend_service.k8s_service.self_link
  ip_address  = var.ip_address
}

resource "google_compute_backend_service" "k8s_service" {
  name                  = "k8-https-backend-service"
  protocol              = "HTTPS"
  health_checks         = [google_compute_http_health_check.k8s_health_check.id]
  backend {
    group = google_compute_instance_group.bookservice.id
  }
}

resource "google_compute_backend_service" "k8s_api_server" {
  name                  = "k8s-api-server"
  protocol              = "TCP" # Kubernetes API uses TCP
  health_checks         = []
  backend {
    group = google_compute_instance_group.k8server.id
  }
}

resource "google_compute_instance_group" "bookservice" {
  name        = "bookservice-webservers"
  description = "Book service instance group"

  instances = [
    var.worker_0, var.worker_1, var.worker_2
  ]

  zone = "us-west1-b"
}

resource "google_compute_instance_group" "k8server" {
  name        = "k8-server"
  description = "K8 instance group"

  instances = [
    var.controller_0, var.controller_1, var.controller_2
  ]

  zone = "us-west1-b"
}