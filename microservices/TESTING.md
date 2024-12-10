## Проверяем работоспособность
Чтобы проверять события из кафки можно в отдельной вкладке поднять под с клиентом кафки и подключиться к нему. Далее в этой вкладке будем консьюмить и паблишить событий в нужные топики.
```bash
kubectl run kafka-client --restart='Never' --image docker.io/bitnami/kafka:3.8.0-debian-12-r0 --namespace default --command -- sleep infinity
kubectl exec --tty -i kafka-client --namespace default -- bash
```

HTTP-запросы из этого файла направлятся в Kusk Gateway (порт 8080).

### Smart Device Service:

#### Получить информацию об устройстве:
```bash
curl -X GET localhost:8080/devices/1
```

#### Обновить информацию об устройстве и убедиться, что статус обновился:
```bash
curl -X PUT  http://localhost:8080/devices/1/status -d '{"status": "turn_on"}'
curl -X GET localhost:8080/devices/1
```
Также можно проверить, что при обновлении устройства улетело событие в kafka в топик device_statuses (выполнять из вкладки с клиентом кафки):
```bash
kafka-console-consumer.sh --bootstrap-server kafka.default.svc.cluster.local:9092 --topic device_statuses --from-beginning
```

#### Отправить команду на устройство (async):
Инициировать отправку:
```bash
kafka-console-producer.sh --broker-list kafka-controller-0.kafka-controller-headless.default.svc.cluster.local:9092,kafka-controller-1.kafka-controller-headless.default.svc.cluster.local:9092,kafka-controller-2.kafka-controller-headless.default.svc.cluster.local:9092 --topic device_commands
```
Вставить команду в stdin:
```bash
 {"device_id": 1, "user_id": 1, "command": "turn_on"}
```
Для отправки команды на устройство в сервисе поставлена заглушка, которая просто логируется факт отправки. В логах должно быть что-то такое:
`sending command: turn_on by user 1 to device id: 1`
```bash
Найти под:
  kubectl get pods
Посмотреть логи пода:
  kubectl logs chart1-smart-device-service-76d66886d7-cvmkj
```

### Telemetry Service:

#### Получить телеметрию по устройству:
```bash
curl -X GET localhost:8080/telemetry/devices/1
```
В самом начале должен в ответе приходить в пустой массив.

#### Отправка события в телеметрию (async):
Инициировать отправку:
```bash
 kafka-console-producer.sh --broker-list kafka-controller-0.kafka-controller-headless.default.svc.cluster.local:9092,kafka-controller-1.kafka-controller-headless.default.svc.cluster.local:9092,kafka-controller-2.kafka-controller-headless.default.svc.cluster.local:9092 --topic sensor_data
```
Вставить событие в stdin:
```bash
 {"device_id": 1, "temperature": 26.4, "type": "temperature"}
```

#### Получить телеметрию по устройству:
```bash
curl -X GET localhost:8080/telemetry/devices/1
```
Должен в ответе приходить массив с событием.