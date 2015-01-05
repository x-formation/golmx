package lmx

/*
#include <lmx.h>
#include <stdio.h>
#include <unistd.h>

void* IntToPtr(int i);

// The gateway functions
void HeartbeatConnectionLost_cgo(void *pVendorData, const char *szHost, int nPort, int nFailedHeartbeats);
void HeartbeatCheckoutFailure_cgo(void *pVendorData, const char *szFeatureName, int nUsedLicCount);
void HeartbeatCheckoutSuccess_cgo(void *pVendorData, const char *szFeatureName, int nUsedLicCount);
void HeartbeatRetryFeature_cgo(void *pVendorData, const char *szFeatureName, int nUsedLicCount);
void HeartbeatExit_cgo(void *pVendorData);
*/
import "C"

import (
	"unsafe"
)

//export HeartbeatConnectionLostGo
func HeartbeatConnectionLostGo(pVendorData unsafe.Pointer, szHost *C.char, nPort int, nFailedHeartbeats int) {
	cb.heartbeatConnectionLost(cb.vendordata, C.GoString(szHost), int(nPort),int(nFailedHeartbeats))
}

//export HeartbeatCheckoutFailureGo
func HeartbeatCheckoutFailureGo(pVendorData unsafe.Pointer, szFeatureName *C.char, nUsedLicCount int, LmxStat Status) {
	cb.heartbeatCheckoutFailure(cb.vendordata, C.GoString(szFeatureName), int(nUsedLicCount),
        LmxStat)
}

//export HeartbeatCheckoutSuccessGo
func HeartbeatCheckoutSuccessGo(pVendorData unsafe.Pointer, szFeatureName *C.char, nUsedLicCount int) {
	cb.heartbeatCheckoutSuccess(cb.vendordata, C.GoString(szFeatureName), int(nUsedLicCount))
}

//export HeartbeatRetryFeatureGo
func HeartbeatRetryFeatureGo(pVendorData unsafe.Pointer, szFeatureName *C.char, nUsedLicCount int) {
	cb.heartbeatRetryFeature(cb.vendordata, C.GoString(szFeatureName), int(nUsedLicCount))
}

//export HeartbeatExitGo
func HeartbeatExitGo(pVendorData unsafe.Pointer) {
	cb.heartbeatExit(cb.vendordata)
}

// SetOption TODO(rjeczalik)
func (c *cgoClient) SetOption(option OptionType, value interface{}) error {
	switch option {
	case
		OptExactVersion,
		OptAllowBorrow,
		OptAllowGrace,
		OptTrialVirtualMachine,
		OptTrialTerminalServer,
		OptBlacklist,
		OptLicenseIdle,
		OptAllowMultipleServers,
		OptClientHostIDToServer,
		OptAllowCheckoutLessLicenses:
		return LookupError(c.setBooleanOpt(option, value))
	case
		OptLicensePath,
		OptCustomShareString,
		OptLicenseString,
		OptServersideRequestString,
		OptCustomUsername,
		OptCustomHostname,
		OptReservationToken,
		OptBindAddress:
		return LookupError(c.setStringOpt(option, value))
	case
		OptTrialDays,
		OptTrialUses,
		OptAutomaticHeartbeatAttempts,
		OptAutomaticHeartbeatInterval,
		OptHostIDCacheCleanupInterval:
		return LookupError(c.setIntegerOpt(option, value))
	case
		OptHostIDEnabled,
		OptHostIDDisabled:
		return LookupError(c.setHostIdTypeOpt(option, value))
	case OptCustomHostIDFunction:
		return LookupError(StatNotImplemented)
	case OptHostIDCompareFunction:
		return LookupError(StatNotImplemented)
	case OptHeartbeatCheckoutFailureFunction:
		return LookupError(c.setHeartbeatCheckoutFailureFunction(option, value))
	case OptHeartbeatCheckoutSuccessFunction:
		return LookupError(c.setHeartbeatCheckoutSuccessFunction(option, value))
	case OptHeartbeatRetryFeatureFunction:
		return LookupError(c.setHeartbeatRetryFeatureFunction(option, value))
	case OptHeartbeatConnectionLostFunction:
		return LookupError(c.setHeartbeatConnectionLostFunction(option, value))
	case OptHeartbeatExitFunction:
		return LookupError(c.setHeartbeatExitFunction(option, value))
	case OptHeartbeatCallbackVendordata:
		return LookupError(c.setHeartbeatCallbackVendordataFunction(option, value))
	default:
		return LookupError(StatInvalidParameter)
	}
	return LookupError(StatInvalidParameter)
}

func (c *cgoClient) setBooleanOpt(opt OptionType, val interface{}) Status {
	var okPtr unsafe.Pointer
	if val != nil {
		switch value := val.(type) {
		case int:
			if value == 1 {
				okPtr = unsafe.Pointer(C.IntToPtr(C.int(1)))
			}
		case bool:
			if value {
				okPtr = unsafe.Pointer(C.IntToPtr(C.int(1)))
			}
		default:
			return StatInvalidParameter
		}
	}
	return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), okPtr))
}

func (c *cgoClient) setStringOpt(opt OptionType, val interface{}) Status {
	var okPtr unsafe.Pointer
	if val != nil {
		if valstr, ok := val.(string); ok {
			okPtr = unsafe.Pointer(C.CString(valstr))
			defer C.free(unsafe.Pointer(okPtr))
		} else {
			return StatInvalidParameter
		}
	}
	return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), okPtr))
}

func (c *cgoClient) setIntegerOpt(opt OptionType, val interface{}) Status {
	if val != nil {
		if valInt, ok := val.(int); ok {
			return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt),
				unsafe.Pointer(C.IntToPtr(C.int(valInt)))))
		} else {
			return StatInvalidParameter
		}
	}
	return Status(C.LMX_SetOption(
		c.handle, C.LMX_SETTINGS(opt), *new(unsafe.Pointer)))
}

func (c *cgoClient) setHostIdTypeOpt(opt OptionType, val interface{}) Status {
	if valHostID, ok := val.(HostIDType); ok {
		return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt),
			unsafe.Pointer(C.IntToPtr(C.int(valHostID)))))
	}
	return StatInvalidParameter
}

func (c *cgoClient) setHeartbeatConnectionLostFunction(opt OptionType, val interface{}) Status {
	if callback, ok := val.(func(interface{}, string, int, int)); ok {
		cb.heartbeatConnectionLost = callback
		return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt),
			unsafe.Pointer(C.HeartbeatConnectionLost_cgo)))
	}
	return StatInvalidParameter
}

func (c *cgoClient) setHeartbeatCheckoutFailureFunction(opt OptionType, val interface{}) Status {
	if callback, ok := val.(func(interface{}, string, int, Status)); ok {
		cb.heartbeatCheckoutFailure = callback
		return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), unsafe.Pointer(C.HeartbeatCheckoutFailure_cgo)))
	}
	return StatInvalidParameter
}

func (c *cgoClient) setHeartbeatCheckoutSuccessFunction(opt OptionType, val interface{}) Status {
	if callback, ok := val.(func(interface{}, string, int)); ok {
		cb.heartbeatCheckoutSuccess = callback
		return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), unsafe.Pointer(C.HeartbeatCheckoutSuccess_cgo)))
	}
	return StatInvalidParameter
}

func (c *cgoClient) setHeartbeatRetryFeatureFunction(opt OptionType, val interface{}) Status {
	if callback, ok := val.(func(interface{}, string, int)); ok {
		cb.heartbeatRetryFeature = callback
		return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt),
			unsafe.Pointer(C.HeartbeatRetryFeature_cgo)))
	}
	return StatInvalidParameter
}

func (c *cgoClient) setHeartbeatExitFunction(opt OptionType, val interface{}) Status {
	if callback, ok := val.(func(interface{})); ok {
		cb.heartbeatExit = callback
		return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt),
			unsafe.Pointer(C.HeartbeatExit_cgo)))
	}
	return StatInvalidParameter
}

func (c *cgoClient) setHeartbeatCallbackVendordataFunction(opt OptionType, val interface{}) Status {
	if vdata, ok := val.(interface{}); ok {
		cb.vendordata = vdata
		return StatSuccess
	}
	return StatInvalidParameter
}
