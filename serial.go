package psxserialnumber

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const serialRegex string = "((SLPS|SLES|SLUS|SCPS|SCUS|SCES|SIPS|SLPM|SLEH|SLED|SCED|ESPM|PBPX|LSP)[_P\\-])|(LSP9|907127)"
const serialCodeDotPosition int = 8
const serialCodeLength int = 11
const bufferSize int = 1024 * 1024

var serialExceptions = map[string]string{
	"SLUSP":  "SLUS_",
	"LSP9":   "LSP_9",
	"907127": "LSP_907127",
}

var m = regexp.MustCompile(serialRegex)

// GetSerial returns the serial for the given PSX bin file
func GetSerial(filepath string) (serial string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return serial, err
	}
	defer file.Close()

	buffer := make([]byte, bufferSize)

	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return serial, err
			}

			break
		}

		serial = findSerial(string(buffer[:bytesread]))
		if serial != "" {
			fmt.Println(serial)
			serial = normalizeSerial(serial)
			break
		}
	}

	if serial == "" {
		return serial, fmt.Errorf("No serial found for the file %s", filepath)
	}

	return serial, nil
}

func findSerial(s string) string {
	serialPosition := m.FindStringIndex(s)
	if len(serialPosition) > 0 {
		return s[serialPosition[0] : serialPosition[0]+serialCodeLength]
	}

	return ""
}

func normalizeSerial(s string) string {
	serial := strings.Replace(s, ".", "", 1)
	serial = strings.Replace(serial, "-", "_", 1)
	serial = strings.Replace(serial, "-", "", 1)

	for key, value := range serialExceptions {
		if strings.Contains(serial, key) {
			serial = strings.Replace(serial, key, value, 1)
		}
	}

	return serial[:serialCodeDotPosition] + "." + serial[serialCodeDotPosition:serialCodeLength-1]
}
