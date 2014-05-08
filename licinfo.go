package lmx

/*
#include <lmx.h>
LMX_HOSTID* IterHostID(LMX_HOSTID* h, int i) {
	return &h[i];
}
int IterTimeZones(int* tz, int i) {
	return tz[i];
}
*/
import "C"

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func goLicenseInfo(CLi *C.LMX_LICENSE_INFO) []LicenseInfo {
	if CLi == nil {
		return []LicenseInfo{}
	}
	return append(goLicenseInfo((*C.LMX_LICENSE_INFO)(CLi.pNext)),
		LicenseInfo{
			Type:          LicenseType(CLi.eLicenseType),
			Status:        Status(CLi.LmxStat),
			Vendor:        C.GoString(&CLi.szVendorName[0]),
			Path:          C.GoString(&CLi.szPath[0]),
			Version:       C.GoString(&CLi.szVersion[0]),
			Port:          int(CLi.nPort),
			UptimeSeconds: time.Duration(CLi.nServerUptimeSeconds) * time.Second,
			Features:      goFeatureInfo(CLi.pFeature),
			Users:         goUserInfo(CLi.pUser),
		})
}

func goUserInfo(CUi *C.LMX_CLIENT_USER) []UserInfo {
	if CUi == nil {
		return []UserInfo{}
	}
	return append(goUserInfo((*C.LMX_CLIENT_USER)(CUi.pNext)),
		UserInfo{
			IP:       net.ParseIP(C.GoString(&CUi.szIP[0])),
			Username: C.GoString(&CUi.szUsername[0]),
			Hostname: C.GoString(&CUi.szHostname[0]),
			LoggedIn: toTime(C.GoString(&CUi.szLoginTime[0])),
			Leases:   goLeaseInfo(CUi.pLease),
			Queues:   goQueueInfo(CUi.pQueue),
			Denials:  goDenialInfo(CUi.pDenial),
		})
}

func goLeaseInfo(CCl *C.LMX_CLIENT_LEASE) []LeaseInfo {
	if CCl == nil {
		return []LeaseInfo{}
	}
	return append(goLeaseInfo((*C.LMX_CLIENT_LEASE)(CCl.pNext)),
		LeaseInfo{
			FeatureName:  C.GoString(&CCl.szFeatureName[0]),
			ProjectName:  C.GoString(&CCl.szProjectName[0]),
			UniqueID:     C.GoString(&CCl.szUniqueID[0]),
			Count:        int(CCl.nLeaseCount),
			Checkout:     toTime(C.GoString(&CCl.szCheckoutTime[0])),
			BorrowExpire: toTime(C.GoString(&CCl.szBorrowExpireTime[0])),
			Share: struct {
				Username string
				Hostname string
				Custom   string
			}{
				C.GoString(&CCl.szShareUsername[0]),
				C.GoString(&CCl.szShareHostname[0]),
				C.GoString(&CCl.szShareCustom[0]),
			},
		})
}

func goQueueInfo(CCq *C.LMX_CLIENT_QUEUE) []QueueInfo {
	if CCq == nil {
		return []QueueInfo{}
	}
	return append(goQueueInfo((*C.LMX_CLIENT_QUEUE)(CCq.pNext)),
		QueueInfo{
			FeatureName: C.GoString(&CCq.szFeatureName[0]),
			Count:       int(CCq.nQueueCount),
			Queued:      toTime(C.GoString(&CCq.szQueueTime[0])),
		})
}

func goDenialInfo(CCd *C.LMX_CLIENT_DENIAL) []DenialInfo {
	if CCd == nil {
		return []DenialInfo{}
	}
	return append(goDenialInfo((*C.LMX_CLIENT_DENIAL)(CCd.pNext)),
		DenialInfo{
			FeatureName: C.GoString(&CCd.szFeatureName[0]),
			UniqueID:    C.GoString(&CCd.szUniqueID[0]),
			Count:       int(CCd.nDenialsCount),
			Denied:      toTime(C.GoString(&CCd.szDenialTime[0])),
		})
}

func goFeatureInfo(CFi *C.LMX_FEATURE_INFO) []FeatureInfo {
	if CFi == nil {
		return []FeatureInfo{}
	}
	return append(goFeatureInfo((*C.LMX_FEATURE_INFO)(CFi.pNext)),
		FeatureInfo{
			Blacklisted:  int(CFi.nBlacklisted) == 0,
			Name:         C.GoString(&CFi.szFeatureName[0]),
			Vendor:       C.GoString(&CFi.szVendorName[0]),
			Comment:      C.GoString(&CFi.szComment[0]),
			Data:         C.GoString(&CFi.szData[0]),
			Licensee:     C.GoString(&CFi.szLicensee[0]),
			Options:      C.GoString(&CFi.szOptions[0]),
			Serial:       C.GoString(&CFi.szSN[0]),
			Path:         C.GoString(&CFi.szPath[0]),
			ID:           C.GoString(&CFi.szUniqueID[0]),
			Key:          C.GoString(&CFi.szKey[0]),
			KeyComment:   C.GoString(&CFi.szKeyComment[0]),
			Port:         int(CFi.nServerPort),
			SoftLimit:    int(CFi.nSoftLimit),
			TrialUses:    int(CFi.nTrialUses),
			KeyType:      KeyType(CFi.eKeyType),
			LicenseType:  LicenseType(CFi.eLicenseType),
			Share:        toShareTypeSlice(int(CFi.nShareCode)),
			MinCheckout:  time.Duration(CFi.nMinCheckoutSeconds) * time.Second,
			Hold:         time.Duration(CFi.nHoldSeconds) * time.Second,
			Borrow:       time.Duration(CFi.nBorrowHours) * time.Hour,
			Grace:        time.Duration(CFi.nGraceHours) * time.Hour,
			ActualBorrow: time.Duration(CFi.nActualBorrowHours) * time.Hour,
			Issued:       toTime(C.GoString(&CFi.szIssuedDate[0])),
			Expire:       toTime(C.GoString(&CFi.szActualExpireTime[0])),
			Version: struct {
				Major int
				Minor int
			}{int(CFi.nMajorVer), int(CFi.nMinorVer)},
			Count: struct {
				Available int
				Used      int
			}{int(CFi.nAvailableLicCount), int(CFi.nUsedLicCount)},
			Expiration: struct {
				Start time.Time
				End   time.Time
			}{
				toTime(C.GoString(&CFi.szStartDate[0])),
				toTime(C.GoString(&CFi.szEndDate[0])),
			},
			Maintenance: struct {
				Start time.Time
				End   time.Time
			}{
				toTime(C.GoString(&CFi.szMaintenanceStartDate[0])),
				toTime(C.GoString(&CFi.szMaintenanceEndDate[0])),
			},
			ReservationCount: struct {
				User int
				Host int
			}{int(CFi.nUserBasedCount), int(CFi.nHostBasedCount)},
			HostIDMatchRate: struct {
				License int
				Actual  int
			}{int(CFi.nHostidLicenseMatchRate), int(CFi.nHostidActualMatchRate)},
			ClockCheck:  ClockCheckType(CFi.nSystemClockCheck),
			Platforms:   toPlatformSlice(C.GoString(&CFi.szPlatforms[0])),
			Client:      goHostId(&CFi.ClientLicenseHostid[0], int(CFi.nClientLicenseHostids)),
			Server:      goHostId(&CFi.ServerLicenseHostid[0], int(CFi.nServerLicenseHostids)),
			Timezones:   goTimeZones(&CFi.sTimeZones[0], int(CFi.nTimeZonesCount)),
			TokenDep:    goTokenDependency(CFi.pTokenDependency),
			Reservation: goDynamicReservation(CFi.pReservations),
		})
}

func goTokenDependency(CTd *C.LMX_TOKEN_DEPENDENCY) []TokenDependency {
	if CTd == nil {
		return []TokenDependency{}
	}
	return append(goTokenDependency((*C.LMX_TOKEN_DEPENDENCY)(CTd.pNext)),
		TokenDependency{
			FeatureName: C.GoString(&CTd.szFeatureName[0]),
			Version: struct {
				Major int
				Minor int
			}{int(CTd.nMajorVer), int(CTd.nMinorVer)},
			Count: int(CTd.nLicCount),
		})
}

func goDynamicReservation(CDr *C.LMX_DYNAMIC_RESERVATION) []DynamicReservation {
	if CDr == nil {
		return []DynamicReservation{}
	}
	return append(goDynamicReservation((*C.LMX_DYNAMIC_RESERVATION)(CDr.pNext)),
		DynamicReservation{
			Token:        C.GoString(&CDr.szToken[0]),
			Expire:       time.Unix(int64(CDr.tActualExpireTime), 0),
			Count:        int(CDr.nLicCount),
			InitialCount: int(CDr.nStartingLicCount),
		})
}

func goHostId(ptr *C.LMX_HOSTID, length int) (hIds []HostID) {
	hIds = make([]HostID, 0, length)
	for i := 0; i < length; i++ {
		hIds = append(hIds, HostID{
			Type:  HostIDType(C.IterHostID(ptr, C.int(i)).eHostidType),
			Value: C.GoString(&C.IterHostID(ptr, C.int(i)).szValue[0]),
			Desc:  C.GoString(&C.IterHostID(ptr, C.int(i)).szDescription[0]),
		})
	}
	return
}

func goTimeZones(ptr *C.int, length int) (tz []int) {
	for i := 0; i < length; i++ {
		tz = append(tz, int(C.IterTimeZones(ptr, C.int(i))))
	}
	return
}

func toShareTypeSlice(shareCode int) (st []ShareType) {
	for _, s := range ShareTypeAll {
		if (uint16(s) & uint16(shareCode)) != 0 {
			st = append(st, s)
		}
	}
	return st
}

func toTime(date string) (t time.Time) {
	if date == "" {
		return
	}
	toks := strings.Fields(date)
	if len(toks) < 2 {
		toks = append(toks, "00:00")
	}
	if t, err := time.Parse(time.RFC3339,
		fmt.Sprintf("%sT%s:00Z", toks[0], toks[1])); err == nil {
		return t
	}
	return
}

func toPlatformSlice(platforms string) (p []PlatformType) {
	for _, platform := range PlatformAll {
		if strings.Contains(platforms, string(*platform)) {
			p = append(p, *platform)
		}
	}
	return p
}
