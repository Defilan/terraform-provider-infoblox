/* resource "infoblox_record_host" "test" { */
/*   name = "test_name" */
/*   ipv4addr { */
/*     address = "1.1.1.1" */
/*   } */
/*   configure_for_dns = false */
/*   comment = "test comment" */
/* } */

provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "https://infoblox.alaskaair.com"
    sslverify = false
}

resource "infoblox_record_host" "host" {
    name = "seadvmaherinfotest007"
    comment = "Bozo test"
    configure_for_dns = false
    ipv4addr {
    address = "10.80.102.48"
    configure_for_dhcp = false
  }
}
/* resource "infoblox_ip" "ip" { */
/*   cidr = "10.80.100.0/22" */
/* } */

/* output "ip" { */
/*  value = "${infoblox_ip.ip.ipaddress}" */
/* } */
