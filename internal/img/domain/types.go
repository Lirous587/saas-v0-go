package domain

type TenantID string

func (t TenantID) IsZero() bool {
	return t == ""
}

func (t TenantID) String() string {
	return string(t)
}

type ImgID string

func (i ImgID) IsZero() bool {
	return i == ""
}

func (i ImgID) String() string {
	return string(i)
}

type CategoryID string

func (c CategoryID) IsZero() bool {
	return c == ""
}

func (c CategoryID) String() string {
	return string(c)
}
