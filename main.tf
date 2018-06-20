provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "https://infoblox.alaskaair.com"
    sslverify = false
    usecookies = false
}

resource "infoblox_record_host" "host" {
  count = 2 
    name = "${format("seadvgk%02d", count.index + 1)}"
    comment = "Bozo test"
    ipv4addr {
      address = "func:nextavailableip:10.80.100.0/22"
  }
  lifecycle {
    ignore_changes = ["ipv4addr", "ipv6addr", "comment", "view"]
  }
    configure_for_dns = false
}

output "ip" {
  depends_on = ["infoblox_record_host.host"]
  value = "${infoblox_record_host.host.*.returnaddress}"
}

