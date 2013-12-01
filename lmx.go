package lmx

// #include <string.h>
// #include <lmx.h>
// #cgo linux freebsd LDFLAGS: -llmxclient -ldl
// LMX_HOSTID* AllocMaxHostID() {
//   return (LMX_HOSTID*)(malloc(LMX_MAX_HOSTIDS * sizeof(LMX_HOSTID)));
// }
// char* AllocLongString() {
//   return (char*)(malloc(LMX_MAX_LONG_STRING_LENGTH * sizeof(char)));
// }
// LMX_HOSTID* IterHostID(LMX_HOSTID *h, int i) {
//   return &h[i];
// }
// void* IntToPtr(int i) {
//   return (void*)(uintptr_t)(i);
// }
import "C"

import (
	"time"
	"unsafe"
)

// Client
type Client interface {
	Close() error
	Checkout(feature string, major, minor, patch int) error
	Checkin(feature string, count int) error
	SetOption(option OptionType, value interface{}) error
	GetErrorMessage() string
	GetErrorMessageSimple(status Status) string
	GetError() *ErrorInfo
	GetHostID(hostid HostIDType) ([]HostID, error)
	GetHostIDSimple(hostid HostIDType) (string, error)
	GetFeatureInfo(feature string) ([]FeatureInfo, error)
	GetLicenseInfo() ([]LicenseInfo, error)
	Heartbeat(feature string) error
	GetExpireTime(feature string) (time.Duration, error)
	ServerLog(feature string, message string) error
	ServerFunction(feature string, request string) (string, error)
	ClientStoreSave(filename string, content string) error
	ClientStoreLoad(filename string) (string, error)
}

// CgoClient
type cgoClient struct {
	handle C.LMX_HANDLE
}

// NewClient TODO(rjeczalik)
func NewClient() (Client, error) {
	c := &cgoClient{}

	if s := Status(C.LMX_Init((*C.LMX_HANDLE)(unsafe.Pointer(&c.handle)))); s != StatSuccess {
		return nil, lookupError(s)
	}

	return c, nil
}

// Close TODO(rjeczalik)
func (c *cgoClient) Close() error {
	if c.handle == nil {
		return ErrNotInitialized
	}
	C.LMX_Free(c.handle)
	return nil
}

// Checkout TODO(rjeczalik)
func (c *cgoClient) Checkout(feature string, major, minor, count int) error {
	cfeature := C.CString(feature)
	defer C.free(unsafe.Pointer(cfeature))

	s := Status(C.LMX_Checkout(c.handle, cfeature, C.int(major), C.int(minor), C.int(count)))
	return lookupError(s)
}

// Checkin TODO(rjeczalik)
func (c *cgoClient) Checkin(feature string, count int) error {
	cfeature := C.CString(feature)
	defer C.free(unsafe.Pointer(cfeature))

	s := Status(C.LMX_Checkin(c.handle, cfeature, C.int(count)))
	return lookupError(s)
}

// SetOption TODO(rjeczalik)
func (c *cgoClient) SetOption(option OptionType, value interface{}) error {
	s := StatInvalidParameter
	switch option {
	case OptExactVersion, OptAllowBorrow, OptAllowGrace, OptTrialVirtualMachine,
		OptTrialTerminalServer, OptBlacklist, OptAllowMultipleServers, OptClientHostIDToServer:
		var ok unsafe.Pointer
		if value != nil {
			switch value := value.(type) {
			case int:
				if value == 1 {
					ok = unsafe.Pointer(C.IntToPtr(1))
				}
			case bool:
				if value {
					ok = unsafe.Pointer(C.IntToPtr(1))
				}
			default:
				return ErrInvalidParameter
			}
		}
		s = Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(option), ok))
	case OptLicensePath, OptCustomShareString, OptLicenseString, OptServersideRequestString,
		OptCustomUsername, OptCustomHostname, OptReservationToken, OptBindAddress:
		var p unsafe.Pointer
		if value != nil {
			if value, ok := value.(string); !ok {
				return ErrInvalidParameter
			} else {
				p = unsafe.Pointer(C.CString(value))
			}
		}
		s = Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(option), p))
		if p != nil {
			C.free(p)
		}
	case OptTrialDays, OptTrialUses, OptAutomaticHeartbeatAttempts, OptAutomaticHeartbeatInterval,
		OptLicenseIdle, OptHostIDCacheCleanupInterval:
		var val C.int
		if value != nil {
			if value, ok := value.(int); ok {
				val = C.int(value)
			} else {
				return ErrInvalidParameter
			}
		}
		s = Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(option), unsafe.Pointer(C.IntToPtr(val))))
	case OptHostIDEnabled, OptHostIDDisabled:
		var val C.int
		if value != nil {
			if value, ok := value.(HostIDType); ok {
				val = C.int(value)
			} else {
				return ErrInvalidParameter
			}
		}
		s = Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(option), unsafe.Pointer(C.IntToPtr(val))))
	case OptCustomHostIDFunction:
		return ErrNotImplemented
	case OptHostIDCompareFunction:
		return ErrNotImplemented
	case OptHeartbeatCheckoutFailureFunction:
		return ErrNotImplemented
	case OptHeartbeatCheckoutSuccessFunction:
		return ErrNotImplemented
	case OptRetryFeatureFunction:
		return ErrNotImplemented
	case OptHeartbeatConnectionLostFunction:
		return ErrNotImplemented
	case OptHeartbeatExitFunction:
		return ErrNotImplemented
	case OptHeartbeatCallbackVendordata:
		return ErrNotImplemented
	default:
		return ErrInvalidParameter
	}
	return lookupError(s)
}

// GetErrorMessage TODO(rjeczalik)
func (c *cgoClient) GetErrorMessage() string {
	return C.GoString(C.LMX_GetErrorMessage(c.handle))
}

// GetErrorMessageSimple TODO(rjeczalik)
func (c *cgoClient) GetErrorMessageSimple(s Status) string {
	return C.GoString(C.LMX_GetErrorMessageSimple(C.LMX_STATUS(s)))
}

// GetError TODO(rjeczalik) TODO(rjeczalik)
func (c *cgoClient) GetError() *ErrorInfo {
	cerr := C.LMX_GetError(c.handle)
	return &ErrorInfo{
		Status:      Status(cerr.LmxStat),
		Internal:    int(cerr.nInternal),
		Context:     int(cerr.nContext),
		Desc:        C.GoString(&cerr.szDescription[0]),
		FeatureName: C.GoString(&cerr.szFeatureName[0]),
	}
}

// GetHostID TODO(rjeczalik) TODO(rjeczalik)
func (c *cgoClient) GetHostID(t HostIDType) ([]HostID, error) {
	cids := C.AllocMaxHostID()
	defer C.free(unsafe.Pointer(cids))
	clen := C.int(0)
	ids := make([]HostID, 0)

	if s := Status(C.LMX_Hostid(c.handle, C.LMX_HOSTID_TYPE(t), cids, &clen)); s != StatSuccess {
		return nil, lookupError(s)
	}

	for i := 0; i < int(clen); i++ {
		ids = append(ids, HostID{})
		ids[i].Type = HostIDType(C.IterHostID(cids, C.int(i)).eHostidType)
		ids[i].Value = C.GoString(&C.IterHostID(cids, C.int(i)).szValue[0])
		ids[i].Desc = C.GoString(&C.IterHostID(cids, C.int(i)).szDescription[0])
	}

	return ids, nil
}

// GetHostIDSimple TODO(rjeczalik)
func (c *cgoClient) GetHostIDSimple(t HostIDType) (string, error) {
	cids := C.AllocLongString()
	defer C.free(unsafe.Pointer(cids))

	if s := Status(C.LMX_HostidSimple(c.handle, C.LMX_HOSTID_TYPE(t), cids)); s != StatSuccess {
		return "", lookupError(s)
	}

	return C.GoString(cids), nil
}

// GetFeatureInfo TODO(rjeczalik)
func (c *cgoClient) GetFeatureInfo(feature string) ([]FeatureInfo, error) {
	return nil, ErrNotImplemented
}

// GetLicenseInfo TODO(rjeczalik)
func (c *cgoClient) GetLicenseInfo() ([]LicenseInfo, error) {
	return nil, ErrNotImplemented
}

// Heartbeat TODO(rjeczalik)
func (c *cgoClient) Heartbeat(feature string) error {
	cfeature := C.CString(feature)
	defer C.free(unsafe.Pointer(cfeature))

	s := Status(C.LMX_Heartbeat(c.handle, cfeature))

	return lookupError(s)
}

// GetExpireTime TODO(rjeczalik)
func (c *cgoClient) GetExpireTime(feature string) (t time.Duration, err error) {
	cfeature := C.CString(feature)
	defer C.free(unsafe.Pointer(cfeature))

	ret := time.Duration(C.LMX_GetExpireTime(c.handle, cfeature))

	switch {
	case ret >= 0:
		t = time.Hour * ret
	case ret == -1:
		err = ErrTooLateDate
	case ret == -2:
		err = ErrDoesNotExpire
	case ret < -2:
		err = ErrUnknownFailure
	}

	return
}

// ServerLog TODO(rjeczalik)
func (c *cgoClient) ServerLog(feature, message string) error {
	cfeature, cmessage := C.CString(feature), C.CString(message)
	defer func() {
		C.free(unsafe.Pointer(cfeature))
		C.free(unsafe.Pointer(cmessage))
	}()

	s := Status(C.LMX_ServerLog(c.handle, cfeature, cmessage))

	return lookupError(s)
}

// ServerFunction TODO(rjeczalik)
func (c *cgoClient) ServerFunction(feature, message string) (string, error) {
	if len(message) >= int(C.LMX_MAX_LONG_STRING_LENGTH) {
		return "", ErrInvalidParameter
	}

	cfeature, cmessage := C.CString(feature), C.CString(message)
	cresponse := C.AllocLongString()
	defer func() {
		C.free(unsafe.Pointer(cfeature))
		C.free(unsafe.Pointer(cmessage))
		C.free(unsafe.Pointer(cresponse))
	}()

	C.strcpy(cresponse, cmessage)

	if s := Status(C.LMX_ServerFunction(c.handle, cfeature, cresponse)); s != StatSuccess {
		return "", lookupError(s)
	}

	return C.GoString(cresponse), nil
}

// ClientStoreSave TODO(rjeczalik)
func (c *cgoClient) ClientStoreSave(filename, content string) error {
	cfilename, ccontent := C.CString(filename), C.CString(content)
	defer func() {
		C.free(unsafe.Pointer(cfilename))
		C.free(unsafe.Pointer(ccontent))
	}()

	s := Status(C.LMX_ClientStoreSave(c.handle, cfilename, ccontent))

	return lookupError(s)
}

// ClientStoreLoad TODO(rjeczalik)
func (c *cgoClient) ClientStoreLoad(filename string) (string, error) {
	cfilename, ccontent := C.CString(filename), C.AllocLongString()
	defer func() {
		C.free(unsafe.Pointer(cfilename))
		C.free(unsafe.Pointer(ccontent))
	}()

	if s := Status(C.LMX_ClientStoreLoad(c.handle, cfilename, ccontent)); s != StatSuccess {
		return "", lookupError(s)
	}

	return C.GoString(ccontent), nil
}
