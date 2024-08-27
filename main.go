package main

import (
	"canbus_temp/bitutils"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Record структура для хранения данных из CSV файла
type Record struct {
	Offset   float64 // смещение в секундах
	ID       string  // ID
	HexValue string  // HEX значение
}

// Config структура для хранения данных из YAML файла
type Config struct {
	DataFile   string    `yaml:"data_file"`
	LocalTime  string    `yaml:"localtime"`
	OutputFile string    `yaml:"output_file"`
	Messages   []Message `yaml:"messages"`
}

// Message структура для описания параметров сообщения
type Message struct {
	CanID      string  `yaml:"can_id"`
	StartBit   int     `yaml:"start_bit"`
	BitLength  int     `yaml:"bit_length"`
	Dlc        int     `yaml:"dlc"`
	Message    string  `yaml:"message"`
	Method     string  `yaml:"method"`  // LSB или MSB
	Scale      float64 `yaml:"scale"`
	Offset     float64 `yaml:"offset"`
}

// parseConfig читает YAML файл конфигурации
func parseConfig(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении файла конфигурации: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге конфигурации: %v", err)
	}

	return &config, nil
}

// parseCSV читает CSV файл данных
func parseCSV(fileName string) ([]Record, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // указали разделитель как точку с запятой

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var result []Record
	for _, record := range records {
		if len(record) < 6 {
			// Если строка имеет меньше чем 6 элементов, пропускаем её
			continue
		}

		offset, err := strconv.ParseFloat(strings.ReplaceAll(record[0], ",", "."), 64)
		if err != nil {
			fmt.Printf("Ошибка при парсинге Offset: %v\n", err)
			continue
		}

		id := strings.TrimSpace(record[3])
		hexValue := strings.TrimSpace(record[5])

		result = append(result, Record{
			Offset:   offset,
			ID:       id,
			HexValue: hexValue,
		})
	}

	return result, nil
}

// calculateValue вычисляет Value с учетом Scale и Offset
func calculateValue(dec int64, scale, offset float64) float64 {
	if scale != 0 || offset != 0 {
		return float64(dec)*scale + offset
	}
	return float64(dec)
}

// parseTimeWithCurrentDate добавляет текущую дату к времени
func parseTimeWithCurrentDate(timeStr string) (time.Time, error) {
	currentDate := time.Now().Format("2006-01-02")
	dateTimeStr := fmt.Sprintf("%sT%s", currentDate, timeStr)
	parsedTime, err := time.Parse("2006-01-02T15:04:05.999", dateTimeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка при парсинге времени: %v", err)
	}
	return parsedTime, nil
}

// processRecords обрабатывает записи и сохраняет результаты
func processRecords(records []Record, config *Config, outputFile string) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла %s: %v", outputFile, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';' // Устанавливаем разделитель как точку с запятой
	defer writer.Flush()

	// Запись заголовков в CSV файл
	headers := []string{"Время", "ID", "DLC", "StartBit", "Length", "HEX", "BIN", "BIN_Converted", "DEC", "Value", "Message"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("ошибка при записи заголовков в файл %s: %v", outputFile, err)
	}

	// Парсим локальное время из конфига
	localTime, err := parseTimeWithCurrentDate(config.LocalTime)
	if err != nil {
		return fmt.Errorf("ошибка при парсинге localtime: %v", err)
	}

	for _, message := range config.Messages {
		for _, record := range records {
			if record.ID == message.CanID {
				var binStr string
				switch message.Method {
				case "LSB":
					binStr, err = bitutils.HexToBinLSB(record.HexValue)
				case "MSB":
					binStr, err = bitutils.HexToBinMSB(record.HexValue)
				default:
					fmt.Printf("Неизвестный метод конвертации для сообщения с ID %s: %v\n", record.ID, message.Method)
					continue
				}
	
				if err != nil {
					fmt.Printf("Ошибка конвертации HEX в BIN для записи с ID %s: %v\n", record.ID, err)
					continue
				}
	
				// Учитываем DLC, проверяем длину бинарной строки
				maxBits := message.Dlc * 8
				if message.StartBit+message.BitLength > maxBits {
					fmt.Printf("Ошибка: диапазон битов выходит за пределы DLC для записи с ID %s\n", record.ID)
					continue
				}
	
				var processedBits string
				switch message.Method {
				case "MSB":
					processedBits, err = bitutils.ProcessBitsMSB(binStr, message.StartBit, message.BitLength)
				default:
					processedBits, err = bitutils.ProcessBits(binStr, message.StartBit, message.BitLength)
				}
				
				if err != nil {
					fmt.Printf("Ошибка обработки битов для записи с ID %s: %v\n", record.ID, err)
					continue
				}
	
				decValue, err := bitutils.BinToDec(processedBits)
				if err != nil {
					fmt.Printf("Ошибка конвертации BIN в DEC для записи с ID %s: %v\n", record.ID, err)
					continue
				}
	
				finalTime := localTime.Add(time.Duration(record.Offset) * time.Second)
				value := calculateValue(decValue, message.Scale, message.Offset)
	
				// Запись результата в CSV файл
				err = writer.Write([]string{
					finalTime.Format(time.RFC3339),
					record.ID,
					strconv.Itoa(message.Dlc),
					strconv.Itoa(message.StartBit),
					strconv.Itoa(message.BitLength),
					record.HexValue,
					binStr,
					processedBits,
					fmt.Sprintf("%d", decValue),
					fmt.Sprintf("%.6f", value),
					message.Message, // Добавлено текстовое сообщение
				})
				if err != nil {
					return fmt.Errorf("ошибка при записи в файл %s: %v", outputFile, err)
				}
			}
		}
	}

	return nil
}

func main() {
	configFile := flag.String("config", "config.yaml", "Путь к файлу конфигурации")
	flag.Parse()

	config, err := parseConfig(*configFile)
	if err != nil {
		fmt.Printf("Ошибка при чтении конфигурационного файла: %v\n", err)
		return
	}

	records, err := parseCSV(config.DataFile)
	if err != nil {
		fmt.Printf("Ошибка при чтении CSV файла: %v\n", err)
		return
	}

	err = processRecords(records, config, config.OutputFile)
	if err != nil {
		fmt.Printf("Ошибка при обработке записей: %v\n", err)
	}
}
