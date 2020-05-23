# psx-serial-number
This library extracts the serial from PSX bin files. The serials are normalized to use always `_` to separate string from numeric part of the serial and `.` for the last two digits of the numeric part. fe: `SLUS_015.52`

## How to use it
```golang
package main

import (
	"fmt"

	psxserialnumber "github.com/ncirocco/psx-serial-number"
)

func main() {
	serial, err := psxserialnumber.GetSerial("/path/to/game.bin")
	if err == nil {
		fmt.Printf("The serial of the game is %s\n", serial)
	}
}
```

This library has been tested with the whole `NTSC-U` library.
