package ogitransformer

import (
	"encoding/json"
	"fmt"

	"github.com/abhishekkr/gol/golenv"
	"github.com/abhishekkr/gol/golerror"

	ogiproducer "github.com/gojekfarm/ogi/producer"
)

type KubernetesKafkaLog struct {
	Message    string     `json:"message"`
	Stream     string     `json:"stream"`
	LogLine    string     `json:"log"`
	Docker     Docker     `json:"docker"`
	Kubernetes Kubernetes `json:"kubernetes"`
	MessageKey string     `json:"message_key"`
}

type Docker struct {
	ContainerId string `json:"container_id"`
}

type Kubernetes struct {
	ContainerName string            `json:"container_name"`
	NamespaceName string            `json:"namespace_name"`
	PodName       string            `json:"pod_name"`
	PodId         string            `json:"pod_id"`
	Labels        map[string]string `json:"labels"`
	Host          string            `json:"host"`
	MasterUrl     string            `json:"master_url"`
}

var (
	KubernetesTopicLabel = golenv.OverrideIfEnv("PRODUCER_KUBERNETES_TOPIC_LABEL", "app")
)

func (kafkaLog *KubernetesKafkaLog) Transform(msgBytes []byte) (err error) {
	if err = json.Unmarshal(msgBytes, &kafkaLog); err != nil {
		err = golerror.Error(123, "failed to parse")
		return
	}

	if kafkaLog.Kubernetes.Labels[KubernetesTopicLabel] == "" {
		err = golerror.Error(123, fmt.Sprintf("correct target topic id '%s' is missing", KubernetesTopicLabel))
		return
	}

	kafkaLog.MessageKey = kafkaLog.Kubernetes.PodName

	msgWithKey, err := json.Marshal(kafkaLog)

	ogiproducer.Produce(kafkaLog.Kubernetes.Labels[KubernetesTopicLabel],
		msgWithKey,
		kafkaLog.MessageKey)
	return
}

func NewKubernetesKafkaLog() Transformer {
	return &KubernetesKafkaLog{}
}
