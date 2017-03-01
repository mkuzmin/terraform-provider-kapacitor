package kapacitor

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/influxdata/kapacitor/client/v1"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},

		ConfigureFunc: configure,

		ResourcesMap: map[string]*schema.Resource{
			"kapacitor_task": taskResource(),
		},
	}
}

func configure(d *schema.ResourceData) (interface{}, error) {
	config := client.Config{
		URL: d.Get("url").(string),
	}
	conn, err := client.New(config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
