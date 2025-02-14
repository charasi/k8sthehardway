module "vpc" {
  source = "../ic/network"
}

module "firewall" {
  source        = "../ic/firewall"
  network_id    = module.vpc.network_id
  subnetwork_id = module.vpc.subnetwork_id
}

module "instances" {
  source                 = "../ic/nodes"
  network_name           = module.vpc.network_id
  subnetwork_name        = module.vpc.subnetwork_id
  master_node_ext_ip     = module.vpc.master_node_ip
  bucket_private_key     = module.buckets.kthw_misc_bucket
  kthw_private_key       = module.instances.private_key_pem
  kthw_private_agent_key = module.instances.private_agent_key_pem
  bucket_name            = module.buckets.kthw_misc_bucket
  static_ip_address      = module.vpc.static_ip_address
}

module "buckets" {
  source = "../ic/bucket"
  //sql_ip_address = module.seagram_database.seagram_ip_address
}


module "target_nodes" {
  source       = "../ic/target_nodes"
  controller_0 = module.instances.controller_0_id
  controller_1 = module.instances.controller_1_id
  controller_2 = module.instances.controller_2_id
  worker_0     = module.instances.worker_0_id
  worker_1     = module.instances.worker_1_id
  worker_2     = module.instances.worker_2_id
  ip_address   = module.vpc.static_ip_address
}


module "seagram_database" {
  source                  = "../ic/mysql"
  seagram_private_network = module.vpc.network_id
  //seagram_cidr_range_name      = module.vpc.seagram_private_subnet
}

output "instance_outputs" {
  value     = module.instances
  sensitive = true
}

output "bucket_outputs" {
  value     = module.buckets
  sensitive = true
}

output "vpc_network_id" {
  value = module.vpc.network_id
}

output "controller_0_id" {
  value = module.instances.controller_0_id
}

output "controller_1_id" {
  value = module.instances.controller_1_id
}

output "controller_2_id" {
  value = module.instances.controller_2_id
}

output "worker_0_id" {
  value = module.instances.worker_0_id
}

output "worker_1_id" {
  value = module.instances.worker_1_id
}

output "worker_2_id" {
  value = module.instances.worker_2_id
}

output "static_ip_address" {
  value = module.vpc.static_ip_address
}