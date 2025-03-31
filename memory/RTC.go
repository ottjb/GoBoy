package memory

import (
	"fmt"
	"time"
)

type RTC struct {
	seconds  byte
	minutes  byte
	hours    byte
	daysLow  byte
	daysHigh byte
	latched  bool
	lastSync time.Time
}

func NewRTC() *RTC {
	rtc := &RTC{
		seconds:  0,
		minutes:  0,
		hours:    0,
		daysLow:  0,
		daysHigh: 0,
		latched:  false,
		lastSync: time.Now(),
	}
	fmt.Println("RTC Initialized:", rtc)
	return rtc
}

func (r *RTC) IsHalted() bool {
	return (r.daysHigh & 0x40) != 0
}

func (r *RTC) SetHalt(halt bool) {
	if halt {
		r.daysHigh |= 0x40 // Set bit 6
	} else {
		r.daysHigh &= 0xBF // Clear bit 6
	}
}

func (r *RTC) GetTotalDays() int {
	highBit := int(r.daysHigh&0x01) << 8
	return highBit | int(r.daysLow)
}

func (r *RTC) HasOverflowed() bool {
	return (r.daysHigh & 0x80) != 0
}
