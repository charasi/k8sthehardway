# Step 2: Set up the Private Services Access connection for Cloud SQL
/**
resource "google_service_networking_connection" "seagram_database_subnetwork" {
  network                 = var.seagram_private_network
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [var.seagram_cidr_range_name]
  update_on_creation_fail = true
}
 */

resource "google_sql_database_instance" "mysql_seagram" {
  name = "mysql-seagram"
  database_version = "MYSQL_8_0_37"
  region = "us-west1"
  deletion_protection = false
  settings {
    tier = "db-custom-2-8192"
    availability_type = "ZONAL"
    activation_policy = "ALWAYS"
    deletion_protection_enabled = false
    disk_size = "10"
    disk_type = "PD_SSD"
    /**
    ip_configuration {
      private_network = var.seagram_private_network
      allocated_ip_range = var.seagram_cidr_range_name
      ipv4_enabled          = false
    }
     */
    location_preference {
      zone = "us-west1-b"
    }
  }

  depends_on = []
  root_password = "recycle"
}

resource "google_sql_database" "bookstore" {
  name     = "bookstore"
  instance = google_sql_database_instance.mysql_seagram.name
}

resource "google_sql_database" "customers" {
  name     = "customers"
  instance = google_sql_database_instance.mysql_seagram.name
}

resource "google_sql_user" "user" {
  instance = google_sql_database_instance.mysql_seagram.name
  name     = "remedy"
  password = "skincream"
}