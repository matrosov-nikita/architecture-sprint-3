resource "helm_release" "kafka" {
    repository = "oci://registry-1.docker.io/bitnamicharts"
    name = "kafka"
    chart = "kafka"
    version = "30.0.0"
    namespace  = "default"

    values = [file("kafka_values.yaml")]
}

resource "null_resource" "create_kafka_topic" {
    depends_on = [helm_release.kafka]

    provisioner "local-exec" {
        command = <<EOT
        echo "Ожидание, пока Kafka станет доступным..."
       until helm status kafka | grep 'deployed'; do
            sleep 1
        done

      echo "Kafka стал доступен. Начинаем создание топиков..."
      kubectl exec -n default -it $(kubectl get pods -l app.kubernetes.io/name=kafka -o jsonpath='{.items[0].metadata.name}') -- \
        kafka-topics.sh --create --topic device_statuses --bootstrap-server kafka.default.svc.cluster.local:9092 --partitions 1 --replication-factor 1 --if-not-exists

      kubectl exec -n default -it $(kubectl get pods -l app.kubernetes.io/name=kafka -o jsonpath='{.items[0].metadata.name}') -- \
        kafka-topics.sh --create --topic device_commands --bootstrap-server kafka.default.svc.cluster.local:9092 --partitions 1 --replication-factor 1  --if-not-exists

      kubectl exec -n default -it $(kubectl get pods -l app.kubernetes.io/name=kafka -o jsonpath='{.items[0].metadata.name}') -- \
        kafka-topics.sh --create --topic sensor_data --bootstrap-server kafka.default.svc.cluster.local:9092 --partitions 1 --replication-factor 1  --if-not-exists
    EOT
    }
}