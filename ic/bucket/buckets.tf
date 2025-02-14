resource "google_storage_bucket" "kthw_misc" {
  location = "us-west1"
  name     = "kthw-misc"
  # Set force_destroy to true to automatically delete objects in the bucket
  force_destroy = true
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "admin_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "admin-csr.json"
  source = "../certificates/admin-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "ca_config" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "ca-config.json"
  source = "../certificates/ca-config.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "ca_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "ca-csr.json"
  source = "../certificates/ca-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "kube_proxy_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "kube-proxy-csr.json"
  source = "../certificates/kube-proxy-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "kubernetes_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "kubernetes-csr.json"
  source = "../certificates/kubernetes-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "worker_0_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "worker-0-csr.json"
  source = "../certificates/worker-0-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "worker_1_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "worker-1-csr.json"
  source = "../certificates/worker-1-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "worker_2_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "worker-2-csr.json"
  source = "../certificates/worker-2-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "encrp_cfg" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "encryption-config.yaml"
  source = "../certificates/encryption-config.yaml"
  content_type  = "application/x-yaml"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "controller_manager_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "controller-manager-csr.json"
  source = "../certificates/controller-manager-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "scheduler_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "scheduler-csr.json"
  source = "../certificates/scheduler-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "etcd_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "etcd-csr.json"
  source = "../certificates/etcd-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "etcd_client_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "etcd-client-csr.json"
  source = "../certificates/etcd-client-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "service_account_csr" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "service-account-csr.json"
  source = "../certificates/service-account-csr.json"
  content_type = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "reverse_proxy" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "nginx.conf"
  source = "../certificates/nginx.conf"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "worker_0" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "worker-0.yml"
  source = "../certificates/worker-0.yml"
  content_type  = "application/x-yaml"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "worker_1" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "worker-1.yml"
  source = "../certificates/worker-1.yml"
  content_type  = "application/x-yaml"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "worker_2" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "worker-2.yml"
  source = "../certificates/worker-2.yml"
  content_type  = "application/x-yaml"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "load_balancer" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "load-balancer-csr.json"
  source = "../certificates/load-balancer-csr.json"
  content_type  = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "load_balancer_config" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "ca-lb-config.json"
  source = "../certificates/ca-lb-config.json"
  content_type  = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "load_balancer_ca" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "ca-lb-csr.json"
  source = "../certificates/ca-lb-csr.json"
  content_type  = "application/json"
}

# Upload a file to the GCS bucket
resource "google_storage_bucket_object" "seagram_database" {
  bucket = google_storage_bucket.kthw_misc.name
  name   = "ca-lb-csr.json"
  source = "../certificates/ca-lb-csr.json"
  content_type  = "application/json"
}

/**
resource "google_storage_bucket_object" "ip_object" {
  name   = "seagram-database.txt"         # The object name in the bucket
  bucket = google_storage_bucket.kthw_misc.name
  content = <<EOF
    ${var.sql_ip_address}
  EOF
}
   */

