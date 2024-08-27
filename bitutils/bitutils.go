package bitutils

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// HexToBinMSB конвертирует HEX строку в бинарную строку в формате MSB
func HexToBinMSB(hexStr string) (string, error) {
	hexStr = strings.ReplaceAll(hexStr, " ", "")
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("ошибка конвертации HEX в BIN (MSB): %v", err)
	}

	var binStr string
	for _, b := range bytes {
		// Форматирование каждого байта в 8-битное бинарное представление с ведущими нулями
		binStr += fmt.Sprintf("%08b", b)
	}
	return binStr, nil
}


// HexToBinLSB конвертирует HEX строку в бинарную строку в формате LSB
func HexToBinLSB(hexStr string) (string, error) {
	hexStr = strings.ReplaceAll(hexStr, " ", "")
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", fmt.Errorf("ошибка конвертации HEX в BIN (LSB): %v", err)
	}

	var binStr string
	for i := len(bytes) - 1; i >= 0; i-- {
		binStr += fmt.Sprintf("%08b", bytes[i])
	}
	return binStr, nil
}

// ProcessBits извлекает битовую подстроку из бинарной строки
func ProcessBits(binStr string, startBit, bitLength int) (string, error) {
	endBit := startBit + bitLength
	if endBit > len(binStr) {
		return "", fmt.Errorf("конечный бит выходит за пределы бинарной строки")
	}
	return binStr[startBit:endBit], nil
}

// ProcessBitsMSB извлекает битовую подстроку из бинарной строки для MSB
func ProcessBitsMSB(binStr string, startBit int, bitLength int) (result string, err error) {
	// Последний бит
	endBit := startBit + bitLength
	tempLength := bitLength

	if endBit > len(binStr) {
		return "", fmt.Errorf("конечный бит выходит за пределы бинарной строки")
	}

	// Индекс первого байта
	firstByteIndex := startBit / 8
	// Индекс последнего байта
	lastByteIndex := (endBit - 1) / 8

	resultBits := ""

	// Если нужные биты в одном байте
	if firstByteIndex == lastByteIndex {
		byteBits := binStr[startBit:endBit]
		return byteBits, nil
	}

	for i := lastByteIndex; i >= firstByteIndex; i-- {
		bitStart := i * 8
		bitEnd := bitStart + 8

		if i == firstByteIndex {
			// Обработка первого байта: отбрасываем биты до startBit
			bitsToKeep := binStr[bitStart : bitStart+bitEnd-startBit]
			resultBits += bitsToKeep
			tempLength = tempLength - (bitStart - startBit)
		} else if i == lastByteIndex {
			// Обработка последнего байта: оставляем только нужные биты
			bitsToKeep := binStr[bitStart:endBit]
			resultBits += bitsToKeep
			tempLength = tempLength - (endBit - bitStart)
		} else {
			// Обработка промежуточных байтов: берем их полностью
			resultBits = resultBits + binStr[bitStart:bitEnd]
			tempLength = tempLength - 8
		}

		// Проверка на ошибки
		if tempLength == 0 {
			panic("tempLength == 0")
		}
	}

		// Добавляем ведущие нули, если нужно
		if len(resultBits) < bitLength {
			resultBits = strings.Repeat("0", bitLength-len(resultBits)) + resultBits
		}
	return resultBits, nil
}

// BitField структура для определения битового поля в смешанном кодировании
type BitField struct {
	Name      string
	StartBit  int
	BitLength int
	Format    string // "MSB" или "LSB"
}

// ExtractBitsFromMixedEncoding извлекает битовые поля с учетом смешанного кодирования
func ExtractBitsFromMixedEncoding(binStr string, fields []BitField) (map[string]string, error) {
	result := make(map[string]string)

	for _, field := range fields {
		var bits string
		var err error

		if field.Format == "LSB" {
			bits, err = ProcessBits(reverseBits(binStr), field.StartBit, field.BitLength)
		} else {
			bits, err = ProcessBits(binStr, field.StartBit, field.BitLength)
		}

		if err != nil {
			return nil, fmt.Errorf("ошибка при извлечении битов для %s: %v", field.Name, err)
		}

		result[field.Name] = bits
	}

	return result, nil
}

// ExtractFlags извлекает флаги из строки с битами
func ExtractFlags(binStr string, flags map[int]string) map[string]bool {
	result := make(map[string]bool)
	for position, name := range flags {
		if position < len(binStr) && binStr[position] == '1' {
			result[name] = true
		} else {
			result[name] = false
		}
	}
	return result
}

// reverseBits разворачивает битовую строку
func reverseBits(bits string) string {
	runes := []rune(bits)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// BinToDec конвертирует бинарную строку в десятичное число
func BinToDec(binStr string) (int64, error) {
	decValue, err := strconv.ParseInt(binStr, 2, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка конвертации BIN в DEC: %v", err)
	}
	return decValue, nil
}
