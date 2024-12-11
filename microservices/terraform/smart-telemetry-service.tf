resource "helm_release" "chart2" {
  name       = "chart2"
  namespace  = "default"
  chart      = "../charts/smart-telemetry-service"
  dependency_update = true
}