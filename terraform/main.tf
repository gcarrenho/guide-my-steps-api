provider "google" {
  project = "guide-my-steps"
  region  = "us-central1"
  zone    = "us-central1-c"
}

/*resource "google_compute_network" "vpc_network" {
  name                    = "guide-mystpes-network"
  auto_create_subnetworks = "true"
}*/

resource "google_compute_network" "vpc_network" {
  name                    = "guide-my-steps-network"
  auto_create_subnetworks = false
  mtu                     = 1460
}

resource "google_compute_subnetwork" "default" {
  name          = "my-custom-subnet"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-west1"
  network       = google_compute_network.vpc_network.id
}
//"vm_instance"
resource "google_compute_instance" "default" {
  name         = "staging-instance"
  machine_type = "f1-micro"
  tags         = ["ssh"]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    //network = google_compute_network.vpc_network.self_link
	subnetwork = google_compute_subnetwork.default.id
    access_config {
    }
  }

  resource "google_compute_firewall" "ssh" {
  name = "allow-ssh"
  allow {
    ports    = ["22"]
    protocol = "tcp"
  }
  direction     = "INGRESS"
  network       = google_compute_network.vpc_network.id
  priority      = 1000
  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["ssh"]
}

}