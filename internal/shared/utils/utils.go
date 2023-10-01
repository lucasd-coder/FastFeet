package utils

import "strings"

func ExtractQueueName(queueURL string) string {
	index := 2
	parts := strings.SplitN(queueURL, "rabbit://", index)

	queueName := parts[1]

	return queueName
}
