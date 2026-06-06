# HU MQ SDK 接入指南

HU MQ 基于 Apache Kafka，**所有标准 Kafka 客户端均可直接接入**。

## Go

```go
package main

import (
    "github.com/IBM/sarama"
)

func main() {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true

    // 生产者
    producer, _ := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
    defer producer.Close()

    msg := &sarama.ProducerMessage{
        Topic: "my-topic",
        Value: sarama.StringEncoder("Hello HU MQ!"),
    }
    producer.SendMessage(msg)

    // 消费者
    consumer, _ := sarama.NewConsumer([]string{"localhost:9092"}, config)
    defer consumer.Close()

    pc, _ := consumer.ConsumePartition("my-topic", 0, sarama.OffsetOldest)
    defer pc.Close()

    for msg := range pc.Messages() {
        println(string(msg.Value))
    }
}
```

## Java

```java
import org.apache.kafka.clients.producer.*;
import java.util.Properties;

public class Producer {
    public static void main(String[] args) {
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:9092");
        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");

        KafkaProducer<String, String> producer = new KafkaProducer<>(props);
        producer.send(new ProducerRecord<>("my-topic", "Hello HU MQ!"));
        producer.close();
    }
}
```

## Python

```python
from kafka import KafkaProducer, KafkaConsumer

# 生产者
producer = KafkaProducer(bootstrap_servers='localhost:9092')
producer.send('my-topic', b'Hello HU MQ!')
producer.close()

# 消费者
consumer = KafkaConsumer(
    'my-topic',
    bootstrap_servers='localhost:9092',
    auto_offset_reset='earliest'
)
for msg in consumer:
    print(msg.value.decode())
```

## 注意事项

1. **连接地址**：使用 HU MQ 集群的 Kafka 原生端口
2. **序列化**：使用标准序列化器，与 Kafka 完全一致
3. **消费者组**：配置 `group.id` 以使用消费组功能
4. **消息确认**：生产者可配置 `acks=all` 保证可靠性