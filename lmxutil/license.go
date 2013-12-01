package lmxutil

import (
	"errors"
	"time"

	"x-formation/lmx"
)

// ErrNonLocalLicense
var ErrNonLocalLicense = errors.New(`lmxutil: not a local license`)

// ErrNonNetworkLicense
var ErrNonNetworkLicense = errors.New(`lmxutil: not a network license`)

// LocalLicense
type LocalLicense struct {
	Name       string
	Licensee   string
	Data       string
	Comment    string
	Options    string
	Replaces   string
	KeyComment string
	Version    struct {
		Major int
		Minor int
	}
	Issued     time.Time
	Expiration struct {
		Start time.Time
		End   time.Time
	}
	Maintenance struct {
		Start time.Time
		End   time.Time
	}
	HostIDMatchRate int
	ClockCheck      lmx.ClockCheckType
	Share           lmx.ShareType
	Timezones       []int
	Platforms       []lmx.PlatformType
	Client          []lmx.HostID
}

// NetworkLicense
type NetworkLicense struct {
	LocalLicense
	Count     int
	SoftLimit int
	HAL       bool
	Serial    string
	KeyType   lmx.KeyType
	Hold      time.Duration
	Borrow    time.Duration
	Grace     time.Duration
	NamedPool struct {
		Users int
		Hosts int
	}
	TokenDep []lmx.TokenDependency
	Server   []lmx.HostID
}

func NewLocalLicense(feature *lmx.FeatureInfo) (lic *LocalLicense, err error) {
	if feature.LicenseType != lmx.LicenseLocal {
		return nil, ErrNonLocalLicense
	}
	lic = &LocalLicense{}
	copyLocalLic(feature, lic)
	return
}

func NewNetworkLicense(feature *lmx.FeatureInfo) (lic *NetworkLicense, err error) {
	if feature.LicenseType != lmx.LicenseNetwork {
		return nil, ErrNonNetworkLicense
	}
	lic = &NetworkLicense{}
	copyLocalLic(feature, &lic.LocalLicense)
	copyNetworkLic(feature, lic)
	return
}

func copyLocalLic(f *lmx.FeatureInfo, lic *LocalLicense) {
	lic.Name = f.Name
	lic.Licensee = f.Licensee
	lic.Data = f.Data
	lic.Comment = f.Comment
	lic.Options = f.Options
	lic.Version = f.Version
	lic.Issued = f.Issued
	lic.Expiration = f.Expiration
	lic.Maintenance = f.Maintenance
	lic.HostIDMatchRate = f.HostIDMatchRate.License
	lic.ClockCheck = f.ClockCheck
	lic.Share = f.Share
	lic.KeyComment = f.KeyComment
	lic.Timezones = f.Timezones
	lic.Platforms = f.Platforms
	lic.Client = f.Client
}

func copyNetworkLic(f *lmx.FeatureInfo, lic *NetworkLicense) {
	lic.Count = f.Count.Available
	lic.SoftLimit = f.SoftLimit
	lic.Serial = f.Serial
	lic.KeyType = f.KeyType
	lic.Hold = f.Hold
	lic.Borrow = f.Borrow
	lic.Grace = f.Grace
	lic.TokenDep = f.TokenDep
	lic.Server = f.Server
}
