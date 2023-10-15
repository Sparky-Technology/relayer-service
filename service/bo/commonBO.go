package bo

import (
	"go.uber.org/zap"
	"strings"
)

func printRequestParam(logTag string, data string) {
	data = strings.ReplaceAll(data, "\n", "")
	data = strings.ReplaceAll(data, " ", "")

	zap.L().Info(logTag, zap.String("request", data))
}
