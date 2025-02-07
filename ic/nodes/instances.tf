resource "google_compute_instance" "controllers" {
    count = 3
    name = "controller-${count.index}"
    network_interface {
      network = var.network_name
      subnetwork = var.subnetwork_name
      network_ip = "10.240.0.1${count.index}"
    }
    boot_disk {
      initialize_params {
        size = 20
        type = var.boot_disk_type
        image = var.boot_disk_image
      }
    }

    metadata = {
      startup-script = file("startup-script.sh")
      ssh-keys = <<EOF
        wisccourant:${tls_private_key.kthw_ssh.public_key_openssh}
        wisccourant:${tls_private_key.kthw_ssh_agent.public_key_openssh}
      EOF
    }
    
    service_account {
      scopes = var.scopes
    }
    
    machine_type = var.machine_type
    can_ip_forward = true
    zone = var.zone
}

resource "google_compute_instance" "workers" {
    count = 3
    name = "worker-${count.index}"
    network_interface {
      network = var.network_name
      subnetwork = var.subnetwork_name
      network_ip = "10.240.0.2${count.index}"
    }
    boot_disk {
      initialize_params {
        size = 20
        type = var.boot_disk_type
        image = var.boot_disk_image
      }
    }

    
  metadata = {
    pod-cidr = "10.200.${count.index}.0/24"
    startup-script = file("worker_script.sh")
    ssh-keys = <<EOF
      wisccourant:${tls_private_key.kthw_ssh.public_key_openssh}
      wisccourant:${tls_private_key.kthw_ssh_agent.public_key_openssh}
    EOF
  }
    
    
    service_account {
      scopes = var.scopes
    }
    
    machine_type = var.machine_type
    can_ip_forward = true
    zone = var.zone
}

resource "google_compute_instance" "master" {
    name = "master"
    depends_on = [google_storage_bucket_object.private_key_object]
    network_interface {
      network = var.network_name
      subnetwork = var.subnetwork_name
      network_ip = "10.240.0.50"
      # Assign external (public) IP address
      access_config {
        # This block assigns a public external IP (NAT IP)
        nat_ip = var.master_node_ext_ip  # Associate the static IP
      }
    }
    
    boot_disk {
      initialize_params {
        size = 20
        type = var.boot_disk_type
        image = var.boot_disk_image
      }
    }
    
    service_account {
      scopes = var.scopes
    }

    metadata = {
      startup-script = file("master-script.sh")
      ssh-keys = "wisccourant:${file("../ic/public_kthw_key.pub")}"
    }
    
    machine_type = var.machine_type
    can_ip_forward = true
    zone = var.zone
}

resource "google_compute_instance" "agent" {
  name = "agent"
  network_interface {
    network = var.network_name
    subnetwork = var.subnetwork_name
    network_ip = "10.240.0.60"
  }

  boot_disk {
    initialize_params {
      size = 20
      type = var.boot_disk_type
      image = var.boot_disk_image
    }
  }

  service_account {
    scopes = var.scopes
  }

  metadata = {
    startup-script = file("startup-script.sh")
    ssh-keys = <<EOF
      wisccourant:${tls_private_key.kthw_ssh.public_key_openssh}
      wisccourant:${tls_private_key.kthw_ssh_agent.public_key_openssh}
    EOF
  }

  machine_type = var.machine_type
  can_ip_forward = true
  zone = var.zone
}

resource "google_compute_instance" "k8main" {
  name = "k8main"
  network_interface {
    network = var.network_name
    subnetwork = var.subnetwork_name
    network_ip = "10.240.0.70"
  }

  boot_disk {
    initialize_params {
      size = 20
      type = var.boot_disk_type
      image = var.boot_disk_image
    }
  }

  service_account {
    scopes = var.scopes
  }

  metadata = {
    startup-script = file("startup-script.sh")
    ssh-keys = <<EOF
      wisccourant:${tls_private_key.kthw_ssh.public_key_openssh}
      wisccourant:${tls_private_key.kthw_ssh_agent.public_key_openssh}
    EOF
  }

  machine_type = var.machine_type
  can_ip_forward = true
  zone = var.zone
}

# Generate SSH key pair
resource "tls_private_key" "kthw_ssh" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

# Generate SSH key pair
resource "tls_private_key" "kthw_ssh_agent" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "google_storage_bucket_object" "private_key_object" {
  name   = "private_key.pem"
  bucket = var.bucket_name
  content = var.kthw_private_key
  content_type = "application/x-pem-file"
}

resource "google_storage_bucket_object" "ip_object" {
  name   = "external_ip.txt"         # The object name in the bucket
  bucket = var.bucket_name
  content = <<EOF
    ${google_compute_instance.master.network_interface[0].access_config[0].nat_ip}
  EOF
}

resource "google_storage_bucket_object" "private_agent_key_object" {
  name   = "private_agent_key.pem"
  bucket = var.bucket_name
  content = var.kthw_private_agent_key
  content_type = "application/x-pem-file"
}

resource "google_storage_bucket_object" "static_ip_address" {
  bucket = var.bucket_name
  name   = "static_ip.txt"
  content = var.static_ip_address
}