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
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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
		URL:       d.Get("url").(string),
		UserAgent: "Terraform",
	}

	if _, ok := d.GetOk("username"); ok {
		config.Credentials = &client.Credentials{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
			Method: client.UserAuthentication,
		}
	}

	conn, err := client.New(config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
