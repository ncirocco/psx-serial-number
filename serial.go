package psxserialnumber

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const cdromStr string = "cdrom:\\"
const serialCodeLength int = 11

// GetSerial returns the serial for the given PSX bin file
func GetSerial(filepath string) (string, error) {
	const BufferSize = 100
	serial := ""
	file, err := os.Open(filepath)
	if err != nil {
		return serial, err
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)
	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				return serial, err
			}

			break
		}

		position := strings.Index(string(buffer[:bytesread]), cdromStr)
		if position >= 0 {
			serial = string(buffer[:bytesread])[position+len(cdromStr) : position+len(cdromStr)+serialCodeLength]

			break
		}
	}

	if serial == "" {
		return serial, fmt.Errorf("No serial found for the file %s", filepath)
	}

	return serial, nil
}
