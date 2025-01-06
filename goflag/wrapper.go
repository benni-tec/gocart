package goflag

type flag[TInfo any] struct {
	info TInfo
}

func (f *flag[TInfo]) Info() *TInfo {
	return &f.info
}
