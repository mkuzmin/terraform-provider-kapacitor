provider "kapacitor" {
  url = "http://localhost:9092"
}

resource "kapacitor_task" "test" {
  name = "test"
  type = "stream"
  tick_script = "${file("test.tick")}"
  database = "test"
}
