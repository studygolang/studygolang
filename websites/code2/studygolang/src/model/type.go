package model

import "time"

type OftenTime time.Time

func (self OftenTime) String() string {
	return time.Time(self).Format("2006-01-02 15:04:05")
}

func (self OftenTime) MarshalBinary() ([]byte, error) {
	return time.Time(self).MarshalBinary()
}

func (self OftenTime) MarshalJSON() ([]byte, error) {
	return time.Time(self).MarshalJSON()
}

func (self OftenTime) MarshalText() ([]byte, error) {
	return time.Time(self).MarshalText()
}

func (this *OftenTime) UnmarshalBinary(data []byte) error {
	t := time.Time(*this)
	return t.UnmarshalBinary(data)
}

func (this *OftenTime) UnmarshalJSON(data []byte) (err error) {
	t := time.Time(*this)
	return t.UnmarshalJSON(data)
}

func (this *OftenTime) UnmarshalText(data []byte) (err error) {
	t := time.Time(*this)
	return t.UnmarshalText(data)
}
