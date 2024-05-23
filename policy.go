package fastimage

type Policy struct {
	imageTypes            []ImageType
	maxImageFileBytesSize int
	minImageFileBytesSize int
	fetchFullImage        bool
}

func NewPolicy() *Policy {
	return &Policy{}
}

func (p *Policy) WithImageTypes(types ...ImageType) *Policy {
	p.imageTypes = types
	return p
}

func (p *Policy) WithMaxImageFileSizeInBytes(size int) *Policy {
	p.maxImageFileBytesSize = size
	return p
}

func (p *Policy) WithMinImageFileSizeInBytes(size int) *Policy {
	p.minImageFileBytesSize = size
	return p
}

func (p *Policy) WithFetchFullImage(fetchFullImage bool) *Policy {
	p.fetchFullImage = fetchFullImage
	return p
}
