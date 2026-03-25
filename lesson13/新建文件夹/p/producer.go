package main

// 导入必要的包
// fmt: 用于打印输出
// log: 用于记录日志和错误信息
// github.com/IBM/sarama: Kafka 客户端库，用于与 Kafka 集群交互
import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	// =========================================
	// 步骤 1: 配置 Kafka 连接参数
	// =========================================
	// 创建一个新的配置对象，用于设置 Kafka 生产者的各种参数
	config := sarama.NewConfig()

	// 设置消息确认级别：WaitForAll 表示等待所有副本节点确认消息
	// 这样可以确保消息不会丢失，是生产环境的推荐配置
	config.Producer.RequiredAcks = sarama.WaitForAll

	// 设置分区策略：使用随机分区器，将消息随机发送到 Topic 的不同分区
	// 这样可以实现负载均衡，简化配置，无需手动指定分区
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	// 开启消息发送成功回调：这样 SendMessage 方法会返回消息的分区和偏移量
	// 便于我们确认消息是否发送成功，适合调试和入门学习
	config.Producer.Return.Successes = true

	// =========================================
	// 步骤 2: 创建 Kafka 生产者
	// =========================================
	// 连接到 Kafka 集群，创建同步生产者
	// 第一个参数是 Kafka 集群的地址列表，这里使用本地地址 127.0.0.1:9092
	// 第二个参数是上面创建的配置对象
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)

	// 错误处理：如果连接失败，打印错误信息并终止程序
	if err != nil {
		log.Fatal("Kafka 生产者连接失败：", err)
	}

	// 使用 defer 语句确保程序结束时关闭生产者连接
	// 这样可以避免资源泄露，是良好的编程习惯
	defer producer.Close()

	// =========================================
	// 步骤 3: 构造要发送的消息
	// =========================================
	// 创建一个 ProducerMessage 对象，包含要发送的消息信息
	msg := &sarama.ProducerMessage{
		// Topic: 指定消息要发送到的主题名称
		// 注意：这个主题需要已经在 Kafka 中创建好
		Topic: "go-test-topic",

		// Value: 消息的内容，使用 StringEncoder 将字符串转换为 Kafka 可接受的格式
		Value: sarama.StringEncoder("Hello Kafka 4.x! 这是 Go 代码发送的第一条消息"),
	}

	// =========================================
	// 步骤 4: 发送消息
	// =========================================
	// 调用 SendMessage 方法发送消息
	// 该方法会返回三个值：
	// 1. partition: 消息被发送到的分区编号
	// 2. offset: 消息在分区中的偏移量
	// 3. err: 发送过程中可能出现的错误
	partition, offset, err := producer.SendMessage(msg)

	// 错误处理：如果发送失败，打印错误信息并终止程序
	if err != nil {
		log.Fatal("消息发送失败：", err)
	}

	// =========================================
	// 步骤 5: 打印发送结果
	// =========================================
	// 打印消息发送成功的信息，包括分区和偏移量
	// 这样我们可以确认消息确实发送成功了
	fmt.Printf("消息发送成功！\n分区：%d\n偏移量：%d\n", partition, offset)
}
