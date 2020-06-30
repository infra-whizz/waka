package waka

type WakaPartitioner interface {
	Flush()
	Create(name string, mbsize int, typeName string)
}
