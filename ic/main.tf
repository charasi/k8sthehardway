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
}

module "target_nodes" {
  source       = "../ic/target_nodes"
  controller_0 = module.instances.controller_0_self_link
  controller_1 = module.instances.controller_1_self_link
  controller_2 = module.instances.controller_2_self_link
  ip_address   = module.vpc.static_ip_address
}

output "instance_outputs" {
  value     = module.instances
  sensitive = true
}

output "bucket_outputs" {
  value     = module.buckets
  sensitive = true
}