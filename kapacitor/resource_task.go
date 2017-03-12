package kapacitor

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/influxdata/kapacitor/client/v1"
	"errors"
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform/helper/hashcode"
)

func taskResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tick_script": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dbrp": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"retention_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "autogen",
							ForceNew: true,
						},
					},
				},
				Set: dbrpHash,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},

		Create: taskResourceCreare,
		Read:   taskResourceRead,
		Update: taskResourceUpdate,
		Delete: taskResourceDelete,
	}
}

func taskResourceCreare(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	var opts client.CreateTaskOptions

	if name, ok := d.GetOk("name"); ok {
		opts.ID = name.(string)
	}

	switch d.Get("type").(string) {
	case "stream":
		opts.Type = client.StreamTask
	case "batch":
		opts.Type = client.BatchTask
	default:
		return errors.New("Unknown task type")
	}

	opts.TICKscript = d.Get("tick_script").(string)

	for _, v := range d.Get("dbrp").(*schema.Set).List() {
		v1 := v.(map[string]interface{})
		d := &client.DBRP{
			Database:        v1["database"].(string),
			RetentionPolicy: v1["retention_policy"].(string),
		}
		opts.DBRPs = append(opts.DBRPs, *d)
	}

	switch d.Get("enabled").(bool) {
	case true:
		opts.Status = client.Enabled
	case false:
		opts.Status = client.Disabled
	}

	task, err := conn.CreateTask(opts)
	if err != nil {
		return err
	}

	d.SetId(task.ID)
	d.Set("name", task.ID)

	return nil
}

func taskResourceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)
	id := d.Id()

	task, err := conn.Task(conn.TaskLink(id), &client.TaskOptions{ScriptFormat: "raw"})
	if err != nil {
		return err
	}

	d.Set("name", id)

	switch task.Type {
	case client.StreamTask:
		d.Set("type", "stream")
	case client.BatchTask:
		d.Set("type", "batch")
	default:
		return errors.New("Unknown task type")
	}

	d.Set("tick_script", task.TICKscript)
	// TODO: multiple connections
	d.Set("database", task.DBRPs[0].Database)
	d.Set("retention_policy", task.DBRPs[0].RetentionPolicy)
	d.Set("enabled", task.Status == client.Enabled)

	return nil
}

func taskResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)
	id := d.Id()

	var opts client.UpdateTaskOptions
	if d.HasChange("enabled") {
		switch d.Get("enabled").(bool) {
		case true:
			opts.Status = client.Enabled
		case false:
			opts.Status = client.Disabled
		}
	}
	_, err := conn.UpdateTask(conn.TaskLink(id), opts)
	if err != nil {
		return err
	}

	return nil
}

func taskResourceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)
	id := d.Id()

	err := conn.DeleteTask(conn.TaskLink(id))
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func dbrpHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s.%s", m["database"].(string), m["retention_policy"].(string)))
	return hashcode.String(buf.String())
}
