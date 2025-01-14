
output "health_check_self_link" {
  value = google_compute_http_health_check.k8s_health_check.self_link
}