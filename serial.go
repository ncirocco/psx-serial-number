package psxserialnumber

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

const serialRegex string = "(SLPS|SLES|SLUS|SCPS|SCUS|SCES|SIPS|SLPM|SLEH|SLED|SCED|ESPM|PBPX)_\\d{3}.\\d{2}"
const serialCodeLength int = 11

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

			break
		}
	}

	if serial == "" {
		return serial, fmt.Errorf("No serial found for the file %s", filepath)
	}

	return serial, nil
}
