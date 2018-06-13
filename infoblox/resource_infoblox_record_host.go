package infoblox

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	infoblox "github.com/defilan/go-infoblox"
	"github.com/hashicorp/terraform/helper/schema"
)

// hostIPv4Schema represents the schema for the host IPv4 sub-resource
func hostIPv4Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"configure_for_dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"host": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cidr": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

// hostIPv6Schema represents the schema for the host IPv4 sub-resource
func hostIPv6Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"configure_for_dhcp": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"host": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"mac": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func infobloxRecordHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceInfobloxHostRecordCreate,
		Read:   resourceInfobloxHostRecordRead,
		Update: resourceInfobloxHostRecordUpdate,
		Delete: resourceInfobloxHostRecordDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv4addr": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Resource{Schema: hostIPv4Schema()},
			},
			"ipv6addr": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Resource{Schema: hostIPv6Schema()},
			},
			"configure_for_dns": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"view": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ipv4sFromList(ipv4s []interface{}, d *schema.ResourceData, meta interface{}) []infoblox.HostIpv4Addr {
	result := make([]infoblox.HostIpv4Addr, 0, len(ipv4s))
	ip, _ := infobloxNextIP(d, meta)

	for _, v := range ipv4s {
		ipMap := v.(map[string]interface{})
		i := infoblox.HostIpv4Addr{}
		// if err != nil {
		// 	return nil
		// }

		i.Ipv4Addr = ip

		if val, ok := ipMap["configure_for_dhcp"]; ok {
			i.ConfigureForDHCP = val.(bool)
		}
		if val, ok := ipMap["host"]; ok {
			i.Host = val.(string)
		}
		if val, ok := ipMap["mac"]; ok {
			i.MAC = val.(string)
		}

		result = append(result, i)
	}
	d.Set("address", ip)
	d.Set("ipv4addr", result)
	return result
}

func ipv6sFromList(ipv6s []interface{}) []infoblox.HostIpv6Addr {
	result := make([]infoblox.HostIpv6Addr, 0, len(ipv6s))
	for _, v := range ipv6s {
		ip := v.(*schema.ResourceData)
		i := infoblox.HostIpv6Addr{}

		if attr, ok := ip.GetOk("address"); ok {
			i.Ipv6Addr = attr.(string)
		}
		if attr, ok := ip.GetOk("configure_for_dhcp"); ok {
			i.ConfigureForDHCP = attr.(bool)
		}
		if attr, ok := ip.GetOk("host"); ok {
			i.Host = attr.(string)
		}
		if attr, ok := ip.GetOk("mac"); ok {
			i.MAC = attr.(string)
		}
		result = append(result, i)
	}
	return result
}

func hostObjectFromAttributes(d *schema.ResourceData, meta interface{}) infoblox.RecordHostObject {
	hostObject := infoblox.RecordHostObject{}
	if attr, ok := d.GetOk("name"); ok {
		hostObject.Name = attr.(string)
	}
	if attr, ok := d.GetOk("configure_for_dns"); ok {
		log.Printf("[DEBUG] FOUND CONFIGURE_FOR_DNS")
		hostObject.ConfigureForDNS = attr.(bool)
	}
	if attr, ok := d.GetOk("comment"); ok {
		hostObject.Comment = attr.(string)
	}
	if attr, ok := d.GetOk("ttl"); ok {
		hostObject.Ttl = attr.(int)
	}
	if attr, ok := d.GetOk("view"); ok {
		hostObject.View = attr.(string)
	}
	if attr, ok := d.GetOk("ipv4addr"); ok {
		hostObject.Ipv4Addrs = ipv4sFromList(attr.([]interface{}), d, meta)
	}
	if attr, ok := d.GetOk("ipv6addr"); ok {
		hostObject.Ipv6Addrs = ipv6sFromList(attr.([]interface{}))
	}

	return hostObject
}

func resourceInfobloxHostRecordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record := url.Values{}
	hostObject := hostObjectFromAttributes(d, meta)

	log.Printf("[DEBUG] Creating Infoblox Host record with configuration: %#v", hostObject)
	opts := &infoblox.Options{
		ReturnFields: []string{"name", "ipv4addr", "ipv6addr", "configure_for_dns", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordHost().Create(record, opts, hostObject)
	if err != nil {
		return fmt.Errorf("error creating infoblox Host record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox Host record created with ID: %s", d.Id())
	return nil
}

func resourceInfobloxHostRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	record, err := client.GetRecordHost(d.Id(), nil)
	if err != nil {
		return handleReadError(d, "Host", err)
	}

	d.Set("name", record.Name)
	if &record.ConfigureForDNS != nil {
		d.Set("configure_for_dns", record.ConfigureForDNS)
	}
	if &record.Comment != nil {
		d.Set("comment", record.Comment)
	}
	if &record.Ttl != nil {
		d.Set("ttl", record.Ttl)
	}
	if &record.View != nil {
		d.Set("view", record.View)
	}
	if &record.Ipv4Addrs != nil {
		result := make([]map[string]interface{}, 1)
		for _, v := range record.Ipv4Addrs {
			i := make(map[string]interface{})

			i["address"] = v.Ipv4Addr
			if &v.ConfigureForDHCP != nil {
				i["configure_for_dhcp"] = v.ConfigureForDHCP
			}
			if &v.Host != nil {
				i["host"] = v.Host
			}
			if &v.MAC != nil {
				i["mac"] = v.MAC
			}

			result = append(result, i)
		}
		d.Set("ipv4addr", result)
	}
	if &record.Ipv6Addrs != nil {
		result := make([]map[string]interface{}, 1)
		for _, v := range record.Ipv6Addrs {
			i := make(map[string]interface{})

			i["address"] = v.Ipv6Addr
			if &v.ConfigureForDHCP != nil {
				i["configure_for_dhcp"] = v.ConfigureForDHCP
			}
			if &v.Host != nil {
				i["host"] = v.Host
			}
			if &v.MAC != nil {
				i["mac"] = v.MAC
			}

			result = append(result, i)
		}
		d.Set("ipv6addr", result)
	}

	return nil
}

func resourceInfobloxHostRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	_, err := client.GetRecordHost(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding infoblox Host record: %s", err.Error())
	}

	record := url.Values{}
	hostObject := hostObjectFromAttributes(d, meta)

	log.Printf("[DEBUG] Updating Infoblox Host record with configuration: %#v", hostObject)

	opts := &infoblox.Options{
		ReturnFields: []string{"name", "ipv4addr", "ipv6addr", "configure_for_dns", "comment", "ttl", "view"},
	}
	recordID, err := client.RecordHostObject(d.Id()).Update(record, opts, hostObject)
	if err != nil {
		return fmt.Errorf("error updating Infoblox Host record: %s", err.Error())
	}

	d.SetId(recordID)
	log.Printf("[INFO] Infoblox Host record updated with ID: %s", d.Id())

	return resourceInfobloxHostRecordRead(d, meta)
}

func resourceInfobloxHostRecordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*infoblox.Client)

	log.Printf("[DEBUG] Deleting Infoblox Host record: %s, %s", d.Get("name").(string), d.Id())
	_, err := client.GetRecordHost(d.Id(), nil)
	if err != nil {
		return fmt.Errorf("error finding Infoblox Host record: %s", err.Error())
	}

	err = client.RecordHostObject(d.Id()).Delete(nil)
	if err != nil {
		return fmt.Errorf("error deleting Infoblox Host record: %s", err.Error())
	}

	return nil
}

func infobloxNextIP(d *schema.ResourceData, meta interface{}) (string, error) {
	var (
		result string
		err    error
		err2   error
	)

	client := meta.(*infoblox.Client)
	excludedAddresses := buildExcludedAddressesArray(d)

	if name, ok := d.GetOk("name"); ok {
		result, err2 = getIPFromHostname(client, name.(string))
	}
	if err2 != nil {
		if cidr, ok := d.GetOk("cidr"); ok {
			result, err = getNextAvailableIPFromCIDR(client, cidr.(string), excludedAddresses)
		}
	}

	if err != nil {
		return "", fmt.Errorf("[ERROR] could not get next ip")
	}
	return result, err
}

func getIPFromHostname(client *infoblox.Client, hostname string) (string, error) {
	var (
		err    error
		result string
	)

	record, err := client.FindRecordHost(hostname)
	if len(record) <= 0 {
		return "", fmt.Errorf("[INFO] no host found")
	}

	if len(record) > 0 {
		for _, v := range record[0].Ipv4Addrs {
			result = v.Ipv4Addr
		}
	}

	return result, err
}

func getNextAvailableIPFromCIDR(client *infoblox.Client, cidr string, excludedAddresses []string) (string, error) {
	var (
		result string
		err    error
		ou     map[string]interface{}
	)

	network, err := getNetworks(client, cidr)

	if err != nil {
		if strings.Contains(err.Error(), "Authorization Required") {
			return "", fmt.Errorf("[ERROR] Authentication Error, Please check your username/password ")
		}
	}

	if len(network) == 0 {
		err = fmt.Errorf("[ERROR] Empty response from client.Network().find. Is %s a valid network?", cidr)
	}

	if err == nil {
		ou, err = client.NetworkObject(network[0]["_ref"].(string)).NextAvailableIP(1, excludedAddresses)
		result = getMapValueAsString(ou, "ips")
		if result == "" {
			err = fmt.Errorf("[ERROR] Unable to determine IP address from response")
		}
	}

	return result, err
}
