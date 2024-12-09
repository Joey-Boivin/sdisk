package models

type Disk struct {
	totalSize uint64
}

func NewDisk(size uint64) *Disk {
	return &Disk{
		totalSize: size,
	}
}

func (d *Disk) GetSpaceLeft() uint64 {
	return d.totalSize
}
