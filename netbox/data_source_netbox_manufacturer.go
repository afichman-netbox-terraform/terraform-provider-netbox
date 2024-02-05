package netbox

import (
	"errors"
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/dcim"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetboxManufacturer() *schema.Resource {
	return &schema.Resource{
		Read:        dataSourceNetboxManufacturerRead,
		Description: `:meta:subcategory:Data Center Inventory Management (DCIM)`,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"slug": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 30),
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetboxManufacturerRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)

	params := dcim.NewDcimManufacturersListParams()

	if filter, ok := d.GetOk("filter"); ok {
		var filterParams = filter.(*schema.Set)
		for _, f := range filterParams.List() {
			id := f.(map[string]interface{})["id"]
			if id != nil {
				vID := id.(int)
				if vID != 0 {
					vIDString := strconv.Itoa(vID)
					params.ID = &vIDString
				}
			}
			name := f.(map[string]interface{})["name"]
			if name != nil {
				vName := name.(string)
				params.Name = &vName
			}
			slug := f.(map[string]interface{})["slug"]
			if slug != nil {
				vSlug := slug.(string)
				params.Slug = &vSlug
			}
		}
	}

	res, err := api.Dcim.DcimManufacturersList(params, nil)
	if err != nil {
		return err
	}

	if *res.GetPayload().Count > int64(1) {
		return errors.New("more than one manufacturer returned, specify a more narrow filter")
	}
	if *res.GetPayload().Count == int64(0) {
		return errors.New("no manufacturer found matching filter")
	}
	result := res.GetPayload().Results[0]
	d.SetId(strconv.FormatInt(result.ID, 10))
	d.Set("name", result.Name)
	d.Set("slug", result.Slug)
	return nil
}
