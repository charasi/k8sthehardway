output "seagram_ip_address" {
  value = google_sql_database_instance.mysql_seagram.ip_address.0.ip_address
}