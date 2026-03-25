package main

// 导入必要的包
// fmt: 用于打印输出
// log: 用于记录日志和错误信息
// time: 用于处理时间相关操作，这里用于设置超时
// github.com/IBM/sarama: Kafka 客户端库，用于与 Kafka 集群交互
import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// =========================================
	// 步骤 1: 配置 Kafka 消费者参数
	// =========================================
	// 创建一个新的配置对象，用于设置 Kafka 消费者的各种参数
	config := sarama.NewConfig()

	// 设置初始偏移量：OffsetNewest 表示从最新的消息开始消费
	// 这样可以避免消费历史消息，适合首次测试
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// 关闭自动提交偏移量：这样可以手动控制偏移量的提交
	// 对于入门阶段，手动控制可以更好地理解消费过程
	config.Consumer.Offsets.AutoCommit.Enable = false

	// =========================================
	// 步骤 2: 创建 Kafka 消费者
	// =========================================
	// 连接到 Kafka 集群，创建消费者
	// 第一个参数是 Kafka 集群的地址列表，这里使用本地地址 127.0.0.1:9092
	// 第二个参数是上面创建的配置对象
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, config)

	// 错误处理：如果连接失败，打印错误信息并终止程序
	if err != nil {
		log.Fatal("Kafka 消费者连接失败：", err)
	}

	// 使用 defer 语句确保程序结束时关闭消费者连接
	// 这样可以避免资源泄露，是良好的编程习惯
	defer consumer.Close()

	// =========================================
	// 步骤 3: 订阅指定 Topic 的分区
	// =========================================
	// 订阅 "go-test-topic" 主题的 0 号分区
	// 三个参数分别是：
	// 1. 主题名称
	// 2. 分区编号，这里使用 0 号分区（简化配置）
	// 3. 初始偏移量，这里使用 OffsetNewest 表示从最新消息开始
	partitionConsumer, err := consumer.ConsumePartition("go-test-topic", 0, sarama.OffsetNewest)

	// 错误处理：如果订阅失败，打印错误信息并终止程序
	if err != nil {
		log.Fatal("订阅 Topic 失败：", err)
	}

	// 使用 defer 语句确保程序结束时关闭分区消费者
	defer partitionConsumer.Close()

	// 打印提示信息，表明消费者已经开始监听消息
	fmt.Println("========= 开始监听 Kafka 4.x 消息 =========")

	// =========================================
	// 步骤 4: 循环监听消息
	// =========================================
	// 使用无限循环持续接收 Kafka 中的消息
	for {
		// 使用 select 语句监听多个通道
		select {
		// 监听消息通道：当有消息到达时，会从 partitionConsumer.Messages() 通道中接收到消息
		case msg := <-partitionConsumer.Messages():
			// 打印接收到的消息详情
			// 包括：消息所在的分区、偏移量和消息内容
			fmt.Printf("收到消息：\n分区：%d\n偏移量：%d\n内容：%s\n\n", msg.Partition, msg.Offset, string(msg.Value))

		// 监听超时通道：当 5 秒内没有消息时，会从 time.After() 通道中接收到信号
		case <-time.After(5 * time.Second):
			// 打印提示信息，避免终端长时间无输出，让用户知道程序还在运行
			fmt.Println("等待消息中...")

		// 监听错误通道：当消费过程中出现错误时，会从 partitionConsumer.Errors() 通道中接收到错误信息
		case err := <-partitionConsumer.Errors():
			// 打印错误信息，但不终止程序，让程序继续运行
			log.Printf("消费消息出错：%v\n", err)
		}
	}
}
