resource "google_compute_target_pool" "k8_controllers_target_pool" {
  name    = "k8-controllers-target-pool"
  region  = "us-west1"  # specify your region
  health_checks = [google_compute_http_health_check.k8s_controllers_health_check.name]
  instances = [var.controller_0, var.controller_1, var.controller_2]
  session_affinity = "NONE"
}

resource "google_compute_http_health_check" "k8s_controllers_health_check" {
  name               = "k8s-controllers-health-check"
  request_path       = "/livez"  # Path for Kubernetes API server health check
  port               = 80        # Port for Kubernetes API server
  check_interval_sec = 10          # Interval between health checks
  timeout_sec        = 5           # Timeout for each health check
  unhealthy_threshold = 3          # Number of failures before considering unhealthy
  healthy_threshold  = 2           # Number of successes before considering healthy
}

resource "google_compute_forwarding_rule" "k8_api_server_forwarding_rule" {
  name        = "k8-api-server-forwarding-rule"
  region      = "us-west1"
  port_range  = "6443"
  target      = google_compute_target_pool.k8_controllers_target_pool.self_link
  ip_address  = var.ip_address
}

resource "google_compute_target_pool" "k8_workers_target_pool" {
  name    = "k8-workers-target-pool"
  region  = "us-west1"  # specify your region
  health_checks = [google_compute_http_health_check.k8s_workers_health_check.name]
  instances = [var.worker_0, var.worker_1, var.worker_2]
  session_affinity = "NONE"
}

resource "google_compute_http_health_check" "k8s_workers_health_check" {
  name               = "k8s-health-check"
  request_path       = "/livez"  # Path for Kubernetes API server health check
  port               = 80        # Port for Kubernetes API server
  check_interval_sec = 10          # Interval between health checks
  timeout_sec        = 5           # Timeout for each health check
  unhealthy_threshold = 3          # Number of failures before considering unhealthy
  healthy_threshold  = 2           # Number of successes before considering healthy
}

resource "google_compute_forwarding_rule" "k8_workers_forwarding_rule" {
  name        = "k8-workers-forwarding-rule"
  region      = "us-west1"
  port_range  = "443"
  target      = google_compute_target_pool.k8_workers_target_pool.self_link
  ip_address  = var.ip_address
}