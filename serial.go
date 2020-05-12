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

var serialExceptions = map[string]string{
	"SLUSP":  "SLUS_",
	"LSP9":   "LSP_9",
	"907127": "LSP_907127",
}

// GetSerial returns the serial for the given PSX bin file
func GetSerial(filepath string) (string, error) {
	const BufferSize = 1024 * 1024
	serial := ""
	file, err := os.Open(filepath)
	if err != nil {
		return serial, err
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)
	m := regexp.MustCompile(serialRegex)
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return serial, err
			}

			break
		}

		positions := m.FindStringIndex(string(buffer[:bytesread]))
		if len(positions) > 0 {
			serial = string(buffer[:bytesread])[positions[0] : positions[0]+serialCodeLength]
			serial = strings.Replace(serial, ".", "", 1)
			serial = strings.Replace(serial, "-", "_", 1)
			serial = strings.Replace(serial, "-", "", 1)

			for key, value := range serialExceptions {
				if strings.Contains(serial, key) {
					serial = strings.Replace(serial, key, value, 1)
				}
			}

			serial = serial[:serialCodeDotPosition] + "." + serial[serialCodeDotPosition:serialCodeLength-1]

			break
		}
	}

	if serial == "" {
		return serial, fmt.Errorf("No serial found for the file %s", filepath)
	}

	return serial, nil
}
