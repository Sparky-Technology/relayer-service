package utils

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func ConvertHexStringToBytes(hexStr string) ([]byte, error) {
	// 去掉可能的 "0x" 前缀
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	// 将十六进制字符串解码为字节数组
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func ConvertHexStringToByte32(hexStr string) ([32]byte, error) {
	// 去掉可能的 "0x" 前缀
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	// 将十六进制字符串解码为字节数组
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return [32]byte{}, err
	}
	// 创建一个 [32]byte 数组
	var byteArray [32]byte

	// 将解码后的字节数组复制到 [32]byte 数组中
	copy(byteArray[:], bytes)
	return byteArray, nil
}

func ConvertHexStringToUint8(hexStr string) (uint8, error) {
	// 将十六进制字符串 V (如"0x1c") 解析为 uint64
	val, err := strconv.ParseUint(hexStr, 0, 8)
	if err != nil {
		return 0, err
	}
	// 将 uint64 转换为 uint8
	return uint8(val), err
}

func ConvertStringToBigInt(str string) (*big.Int, error) {
	var bigInt big.Int
	_, success := bigInt.SetString(str, 0) // 自动检测进制
	if !success {
		return nil, fmt.Errorf("无法转换为大整数")
	}
	return &bigInt, nil
}

func ConvertStringToBigFloat(str string) (*big.Float, error) {
	var bigFloat big.Float
	_, success := bigFloat.SetString(str)
	if !success {
		return nil, fmt.Errorf("无法转换为大浮点数")
	}
	return &bigFloat, nil
}

func DivideBigIntsToBigFloat(dividend, divisor *big.Int) *big.Float {
	numerator := new(big.Float).SetInt(dividend)
	denominator := new(big.Float).SetInt(divisor)
	numerator.Quo(numerator, denominator)
	return numerator
}
