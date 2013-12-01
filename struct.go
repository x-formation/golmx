package lmx

// #include <lmx.h>
import "C"

import (
	"net"
	"time"
)

// TokenDependency
type TokenDependency struct {
	FeatureName string
	Version     struct {
		Major int
		Minor int
	}
	Count int
	Next  *TokenDependency
}

// DynamicReservation
type DynamicReservation struct {
	Token        string
	InitialCount int
	Count        int
	Expire       time.Time
	Next         *DynamicReservation
}

// HostID
type HostID struct {
	Type  HostIDType
	Value string
	Desc  string
}

// FeatureInfo
type FeatureInfo struct {
	Blacklisted bool
	Name        string
	Vendor      string
	Comment     string
	Data        string
	Licensee    string
	Options     string
	Serial      string
	Path        string
	ID          string
	Key         string
	KeyComment  string
	Port        int
	SoftLimit   int
	TrialUses   int
	KeyType     KeyType
	LicenseType LicenseType
	Share       ShareType
	MinCheckout time.Duration
	Hold        time.Duration
	Borrow      time.Duration
	Grace       time.Duration
	Version     struct {
		Major int
		Minor int
	}
	Count struct {
		Available int
		Used      int
	}
	Expiration struct {
		Start time.Time
		End   time.Time
	}
	Maintenance struct {
		Start time.Time
		End   time.Time
	}
	ReservationCount struct {
		User int
		Host int
	}
	HostIDMatchRate struct {
		License int
		Actual  int
	}
	Expire      time.Time
	Issued      time.Time
	ClockCheck  ClockCheckType
	Platforms   []PlatformType
	Timezones   []int
	Client      []HostID
	Server      []HostID
	TokenDep    []TokenDependency
	Reservation []DynamicReservation
}

// QueueInfo
type QueueInfo struct {
	FeatureName string
	Count       int
	Queued      time.Time
}

// LeaseInfo
type LeaseInfo struct {
	FeatureName  string
	ProjectName  string
	UniqueID     string
	Count        int
	Checkout     time.Time
	BorrowExpire time.Time
	Share        struct {
		Username string
		Hostname string
		Custom   string
	}
}

// DenialInfo
type DenialInfo struct {
	FeatureName string
	UniqueID    string
	Count       int
	Denied      time.Time
}

// UserInfo
type UserInfo struct {
	IP       net.IP
	Username string
	Hostname string
	LoggedIn time.Time
	Leases   []LeaseInfo
	Queues   []QueueInfo
	Denials  []DenialInfo
}

// LicenseInfo
type LicenseInfo struct {
	Type          LicenseType
	Status        Status
	Vendor        string
	Path          string
	Port          int
	Version       string
	UptimeSeconds int
	Features      []FeatureInfo
	Users         []UserInfo
}

// ErrorInfo
type ErrorInfo struct {
	Status      Status
	Internal    int
	Context     int
	Desc        string
	FeatureName string
}
