package models

type Disk struct {
	totalSize uint64
}

func NewDisk(size uint64) *Disk {
	return &Disk{
		totalSize: size,
	}
}
