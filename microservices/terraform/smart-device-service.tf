resource "helm_release" "chart1" {
  name       = "chart1"
  namespace  = "default"
  chart      = "../smart-device-service"
}