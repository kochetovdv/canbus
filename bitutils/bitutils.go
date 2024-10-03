package bitutils

import (
	"fmt"
	"strings"
)

// BytesToBin преобразует массив байтов в бинарную строку с учетом порядковости (LSB или MSB).
func BytesToBin(data []byte, isIntel bool) string {
	var binStr string
	if isIntel {
		// Для LSB (Intel) разворачиваем байты и биты в каждом байте
		for i := len(data) - 1; i >= 0; i-- {
			binStr += reverseString(fmt.Sprintf("%08b", data[i]))
		}
		binStr = reverseString(binStr)
	} else {
		// Для MSB (Motorola) оставляем как есть
		for _, b := range data {
			binStr += fmt.Sprintf("%08b", b)
		}
	}
	return binStr
}

// ExtractSignal извлекает значение сигнала и возвращает также извлеченные биты в виде строки.
func ExtractSignal(data []byte, startBit, length int, isIntel bool) (uint64, string, error) {
	var result uint64
	var bits []string

	if isIntel {
		// LSB (Intel)
		for i := 0; i < length; i++ {
			bitPosition := startBit + i
			byteIndex := bitPosition / 8
			bitIndex := bitPosition % 8
			if byteIndex >= len(data) {
				return 0, "", fmt.Errorf("byte index out of range")
			}
			bitValue := (data[byteIndex] >> bitIndex) & 0x01
			result |= uint64(bitValue) << i
			bits = append([]string{fmt.Sprintf("%d", bitValue)}, bits...)
		}
	} else {
		// MSB (Motorola)
		for i := 0; i < length; i++ {
			bitPosition := startBit - i
			byteIndex := bitPosition / 8
			bitIndex := bitPosition % 8
			if byteIndex >= len(data) || byteIndex < 0 {
				return 0, "", fmt.Errorf("byte index out of range")
			}
			bitValue := (data[byteIndex] >> (7 - bitIndex)) & 0x01
			result |= uint64(bitValue) << (length - 1 - i)
			bits = append([]string{fmt.Sprintf("%d", bitValue)}, bits...)
		}
	}
	extractedBits := strings.Join(bits, "")
	return result, extractedBits, nil
}

// reverseString разворачивает строку
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
