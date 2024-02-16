package gcp

import "time"

type ObjectAttrs struct {
	path         string
	size         int64
	storageClass string
	lastModified time.Time
	checksum     uint32
}
