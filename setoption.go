package lmx

/*
#include <lmx.h>
void* IntToPtr(int i) {
	return (void*)(uintptr_t)(i);
}
*/
import "C"

import (
	"unsafe"
)

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
		return LookupError(StatNotImplemented)
	case OptHeartbeatCheckoutSuccessFunction:
		return LookupError(StatNotImplemented)
	case OptRetryFeatureFunction:
		return LookupError(StatNotImplemented)
	case OptHeartbeatConnectionLostFunction:
		return LookupError(StatNotImplemented)
	case OptHeartbeatExitFunction:
		return LookupError(StatNotImplemented)
	case OptHeartbeatCallbackVendordata:
		return LookupError(StatNotImplemented)
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
