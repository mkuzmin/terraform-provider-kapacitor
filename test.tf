provider "kapacitor" {
  url = "http://localhost:9092"
//  username = "test"
//  password = "secret"
}

resource "kapacitor_task" "test" {
  name = "test"
  type = "stream"
  tick_script = "${file("test.tick")}"
  database = "test"
}
