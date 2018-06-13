provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "https://infoblox.alaskaair.com"
    sslverify = false
    usecookies = false
}


resource "infoblox_ip" "ip" {
  cidr = "10.80.100.0/22"
  count = 2
  hostname = "${format("seadvgk%02d", count.index + 1)}"
}

/* resource "infoblox_record_host" "test" { */
/*   name = "${format("seadvgk%02d", count.index + 1)}" */
/*   configure_for_dns = false */
/*   comment = "test comment" */
/*   view = "" */
/*   count = 3 */
/*   ipv4addr = { */
/*     cidr = "10.80.100.0/22" */
/*   } */
/* } */

