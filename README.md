# Terraform Provider for Kapacitor

```
provider "kapacitor" {
  url = "http://kapacitor:9092"
  
  // username = "test"
  // password = "secret"
}

resource "kapacitor_task" "test" {
  name = "test"
  type = "stream"
  tick_script = "${file("test.tick")}"
  database = "test"
  
  // retention_policy = "autogen"
  // enabled = true
}
```
