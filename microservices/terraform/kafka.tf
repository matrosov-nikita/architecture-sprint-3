resource "helm_release" "kafka" {
    repository = "oci://registry-1.docker.io/bitnamicharts"
    name = "kafka"
    chart = "kafka"
    version = "30.0.0"
    namespace  = "default"

    values = [file("kafka_values.yaml")]
}

resource "kubernetes_job" "create_kafka_topics" {
    metadata {
        name      = "create-kafka-topics"
        namespace = "default"
    }

    spec {
        template {
            metadata {
                labels = {
                    app = "kafka"
                }
            }

            spec {
                init_container {
                    name  = "wait-for-kafka"
                    image = "busybox"
                    command = ["sh", "-c", "until nc -z kafka.default.svc.cluster.local 9092; do echo 'Waiting for Kafka to be ready...'; sleep 2; done"]
                }

                container {
                    image = "bitnami/kafka:latest"
                    name  = "create-topics"
                    command = [
                        "sh",
                        "-c",
                        <<EOT
kafka-topics.sh --create --topic device_statuses --bootstrap-server kafka.default.svc.cluster.local:9092 --partitions 1 --replication-factor 1 || true
kafka-topics.sh --create --topic device_commands --bootstrap-server kafka.default.svc.cluster.local:9092 --partitions 1 --replication-factor 1 || true
kafka-topics.sh --create --topic sensor_data --bootstrap-server kafka.default.svc.cluster.local:9092 --partitions 1 --replication-factor 1 || true
EOT
                    ]
                }

                restart_policy = "OnFailure"
            }
        }
    }
}
