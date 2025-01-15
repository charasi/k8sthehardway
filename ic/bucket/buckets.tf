resource "google_storage_bucket" "kthw_misc" {
  location = "us-west1"
  name     = "kthw-misc"
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