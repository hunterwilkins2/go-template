terraform {
  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5"
    }
  }
}

provider "google" {
  project = var.gcp_project_name
  region  = var.gcp_region
  zone    = var.gcp_zone
}

resource "google_compute_instance" "vm_instance" {
  name         = var.vm_name
  machine_type = var.vm_type

  boot_disk {
    initialize_params {
      image = var.vm_image
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.id
    access_config {
    }
  }

  metadata = {
    ssh-keys = "${var.SSH_USER}:${var.SSH_PUB}"
  }
}

resource "google_compute_network" "vpc_network" {
  name                    = var.vpc_network_name
  auto_create_subnetworks = "true"
}

provider "cloudflare" {
}

resource "cloudflare_dns_record" "a_record" {
  zone_id = var.cloudflare_zone
  name    = var.domain_name
  type    = "A"
  comment = "GCP VPC A record"
  content = google_compute_network.vpc_network.gateway_ipv4
  proxied = true
  ttl     = 1 # automatic
}

resource "cloudflare_dns_record" "www_record" {
  zone_id = var.cloudflare_zone
  name    = "www"
  type    = "CNAME"
  comment = "GCP VPC network www subdomain"
  content = var.domain_name
  proxied = true
  ttl     = 1 # automatic
}

resource "cloudflare_dns_record" "docs_a_record" {
  zone_id = var.cloudflare_zone
  name    = join(".", [var.docs_subdomain, var.domain_name])
  type    = "CNAME"
  comment = "Github Pages"
  content = var.github_pages_url
  proxied = true
  ttl     = 1 # automatic
}

