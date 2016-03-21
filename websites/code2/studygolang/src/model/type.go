package model

import (
	"errors"
	"time"
)

type OftenTime time.Time

func (self OftenTime) String() string {
	return time.Time(self).Format("2006-01-02 15:04:05")
}

func (self OftenTime) MarshalBinary() ([]byte, error) {
	return time.Time(self).MarshalBinary()
}

func (self OftenTime) MarshalJSON() ([]byte, error) {
	t := time.Time(self)
	if y := t.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(t.Format(`"2006-01-02 15:04:05"`)), nil
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
