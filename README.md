# CAN Message Decoder

This project is a tool for decoding CAN (Controller Area Network) messages from a CSV file using a YAML configuration file. It reads the CAN data from a CSV file, decodes the specified messages based on the configuration file, and outputs the results to another CSV file. The tool supports both Intel (LSB) and Motorola (MSB) bit ordering.

## Table of Contents
- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Output Format](#output-format)
- [License](#license)

## Features
- **Parsing CAN messages from CSV**: Reads and decodes hexadecimal CAN messages based on a specified configuration.
- **Customizable configuration**: Easily configurable using a YAML file.
- **Support for Intel and Motorola bit ordering**: Decodes CAN messages using either LSB (Intel) or MSB (Motorola) format.
- **Flexible output**: Outputs the decoded values to a CSV file for further analysis.
- **Local time adjustments**: Calculates the timestamp of each message based on a specified local time.

## Installation
1. **Clone the repository**:
```bash
git clone https://github.com/kochetovdv/canbus.git
```
2. **Navigate to the project directory**:
```bash
cd can-decoder
```
3. **Install dependencies and build the project**:
```bash
go mod tidy
go build -o can-decoder
```

## Configuration

The configuration is defined in a `config.yaml` file, which specifies the input data file, local time, output file, and the messages to decode. Below is an example configuration file:

```yaml
data_file: "data.csv"
localtime: "16:48:09.822"
output_file: "output.csv"
messages:
  - can_id: "107"
    start_bit: 24
    bit_length: 12
    dlc: 8
    message: "Engine RPM"
    method: "LSB"  # "LSB" for Intel, "MSB" for Motorola
    scale: 3.0
    offset: 0.0
```

## Configuration Parameters
- **data_file**: The path to the CSV file containing CAN messages.
- **localtime**: The starting local time (used to calculate absolute timestamps).
- **output_file**: The path to the output CSV file where decoded messages will be saved.
- **messages**: A list of messages to decode:
    - **can_id**: The CAN ID of the message to decode.
    - **start_bit**: The starting bit position of the signal within the message.
    - **bit_length**: The length of the signal in bits.
    - **dlc**: Data length code of the message.
    - **message**: A description of the signal (e.g., "Engine RPM").
    - **method**: The bit ordering method. Use "LSB" for Intel and "MSB" for Motorola.
    - **scale**: Scale factor to apply to the decoded value.
    - **offset**: Offset value to add to the decoded value.

## Usage
To run the program, use the following command:

```bash
./can-decoder -config config.yaml
```

## Command-Line Options
- -config: The path to the configuration YAML file (default is config.yaml).

## Output Format
The decoded messages are saved to the specified output CSV file (`output.csv` by default) with the following columns:

| Time                     | ID  | DLC | StartBit | Length | HEX        | BIN               | BIN_Converted   | DEC   | Value         | Message         |
|--------------------------|-----|-----|----------|--------|------------|-------------------|-----------------|-------|---------------|-----------------|
| Timestamp of the message | CAN ID | Data Length Code | Start bit of the signal | Length of the signal in bits | Original HEX value | Binary representation | Extracted binary bits | Decimal value | Final calculated value | Message description |

## Example Output
```css
2023-09-15T16:48:09.822+00:00;107;8;24;12;0x12AB;0001001010101011;001010101011;2731;8193.000000;Engine RPM
```

## License
This project is licensed under the MIT License. See the LICENSE file for more details.
