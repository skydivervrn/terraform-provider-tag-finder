terraform {
  required_providers {
    tag-finder = {
      source  = "skydivervrn/tag-finder"
      version = "0.0.5"
    }
  }
}

provider "tag-finder" {}

locals {
  required_version = ">3.3.41"
}

variable "current_version" {
  default = "3.3.43"
}

data "version_validator" "example" {
  provider         = tag-finder
  current_version  = var.current_version
  required_version = local.required_version
}