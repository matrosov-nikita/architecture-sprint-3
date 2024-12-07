resource "helm_release" "chart1" {
  name       = "chart1"
  namespace  = "default"
  chart      = "../charts/smart-device-service"
  dependency_update = true
}