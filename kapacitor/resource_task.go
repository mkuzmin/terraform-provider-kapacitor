package kapacitor

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/influxdata/kapacitor/client/v1"
	"errors"
)

func taskResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
			"database": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
		},

		Create: taskResourceCreare,
		Read:   taskResourceRead,
		//Update: taskResourceUpdate,
		Delete: taskResourceDelete,
	}
}

func taskResourceCreare(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.Client)

	name := d.Get("name").(string)

	var task_type client.TaskType
	switch d.Get("type").(string) {
	case "stream":
		task_type = client.StreamTask
	case "batch":
		task_type = client.BatchTask
	default:
		return errors.New("Unknown task type")
	}

	tick_script := d.Get("tick_script").(string)
	database := d.Get("database").(string)
	retention_policy := d.Get("retention_policy").(string)

	var status client.TaskStatus
	switch d.Get("enabled").(bool) {
	case true:
		status = client.Enabled
	case false:
		status = client.Disabled
	}

	task , err := conn.CreateTask(client.CreateTaskOptions{
		ID:         name,
		Type:       task_type,
		TICKscript: tick_script,
		DBRPs: []client.DBRP{{
			Database:        database,
			RetentionPolicy: retention_policy,
		}},
		Status: status,
	})
	if err != nil {
		return err
	}

	d.SetId(task.ID)

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
