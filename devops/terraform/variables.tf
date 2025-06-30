variable "domain_name" {
  description = "Domain name"
  type        = string
  default     = "hivestack.dev"
}

variable "docs_subdomain" {
  description = "Domain name for documentation"
  type        = string
  default     = "docs"
}

variable "SSH_PUB" {
  description = "SSH Public Key to access VM"
  type        = string
}

variable "SSH_USER" {
  description = "SSH User"
  type        = string
  default     = "orion"
}

variable "gcp_project_name" {
  description = "GCP Project Name"
  type        = string
  default     = "go-template"
}

variable "gcp_region" {
  description = "GCP Region"
  type        = string
  default     = "us-east1"
}

variable "gcp_zone" {
  description = "GCP Zone"
  type        = string
  default     = "us-east1-b"
}

variable "vm_name" {
  description = "Virtual Machine Name"
  type        = string
  default     = "go-template-instance"
}

variable "vm_type" {
  description = "Virtual Machine Type"
  type        = string
  default     = "e2-micro"
}

variable "vm_image" {
  description = "Virtual Machine Image"
  type        = string
  default     = "ubuntu-os-cloud/ubuntu-2404-lts-amd64"
}

variable "vpc_network_name" {
  description = "VPC Network Name"
  type        = string
  default     = "go-template-network"
}

variable "cloudflare_zone" {
  description = "Cloudflare Zone ID"
  type        = string
  default     = "31c9a5dc2c98865baddb71c4539aefb8"
}

variable "github_pages_url" {
  description = "Github Pages URL"
  type        = string
  default     = "hunterwilkins.github.io"
}
