# resource provider
terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 6.14.1"
    }
  }
}


provider "google" {
  project     = "kthw-448206"
  region      = "us-west1"
  credentials = "kthw.json"
}

provider "tls" {}