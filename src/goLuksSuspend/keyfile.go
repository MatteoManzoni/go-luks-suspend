package goLuksSuspend

import (
	"os"
	"strconv"
	"strings"
)

type Keyfile struct {
	Path   string
	Offset int
	Size   int
}

func parseKeyfileFromCrypttabEntry(line string) (name string, key Keyfile) {
	fields := strings.Fields(line)

	// fields: name, device, keyfile, options
	//
	// crypttab(5):
	// The third field specifies the encryption password. If the field is
	// not present or the password is set to "none" or "-", the password
	// has to be manually entered during system boot.
	if len(fields) < 3 || fields[2] == "-" || fields[2] == "none" {
		return "", Keyfile{}
	}

	k := Keyfile{Path: fields[2]}

	if len(fields) >= 4 {
		opts := strings.Split(fields[3], ",")
		for i := range opts {
			kv := strings.SplitN(opts[i], "=", 2)
			if len(kv) < 2 {
				continue
			} else if kv[0] == "keyfile-offset" {
				n, err := strconv.Atoi(kv[1])
				if err != nil {
					continue
				}
				k.Offset = n
			} else if kv[0] == "keyfile-size" {
				n, err := strconv.Atoi(kv[1])
				if err != nil {
					continue
				}
				k.Size = n
			}
		}
	}

	return fields[0], k
}

func (k *Keyfile) Defined() bool {
	return len(k.Path) > 0
}

func (k *Keyfile) Available() bool {
	if !k.Defined() {
		return false
	}
	_, err := os.Stat(k.Path)
	return !os.IsNotExist(err)
}
