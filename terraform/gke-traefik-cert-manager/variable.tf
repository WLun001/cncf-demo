variable "project" {
  type = string
}

variable "credentials_file" {}

variable "google_dns_sa_file" {}

variable "region" {
  default = "asia-southeast1"
}

variable "zone" {
  default = "asia-southeast1-a"
}
