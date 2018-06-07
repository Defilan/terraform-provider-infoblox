provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "https://infoblox.alaskaair.com"
    sslverify = false
    usecookies = false
}

resource "infoblox_record_host" "test" {
  name = "${format("seadvgk%02d", count.index + 1)}"
  ipv4addr = {
    address = "10.80.100.88"
    configure_for_dhcp = false
  }
  configure_for_dns = false
  comment = "test comment"
  view = ""
  count = 3
  cidr = "10.80.100.0/22"
}

