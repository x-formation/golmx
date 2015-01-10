package lmx

// #include <lmx.h>
//
// void goConnectionLost(void*, char*, int, int);
// void goCheckoutFailure(void*, char*, int, LMX_STATUS);
// void goCheckoutSuccess(void*, char*, int);
// void goCheckoutRetry(void*, char*, int);
// void goExit(void*);
import "C"

import (
	"net"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// TODO(rjeczalik): ~QuoteASCII every string

//export goConnectionLost
func goConnectionLost(p unsafe.Pointer, host *C.char, port C.int, n C.int) {
	(*cgoClient)(p).heartbeat(&Heartbeat{
		Type:       OptHeartbeatConnectionLost,
		Addr:       net.JoinHostPort(C.GoString(host), strconv.Itoa(int(port))),
		Heartbeats: int(n),
	})
}

//export goCheckoutFailure
func goCheckoutFailure(p unsafe.Pointer, feature *C.char, n C.int, s C.LMX_STATUS) {
	(*cgoClient)(p).heartbeat(&Heartbeat{
		Type:     OptHeartbeatCheckoutFailure,
		Feature:  C.GoString(feature),
		Err:      ToError(Status(s)),
		Licenses: int(n),
	})
}

//export goCheckoutSuccess
func goCheckoutSuccess(p unsafe.Pointer, feature *C.char, n C.int) {
	(*cgoClient)(p).heartbeat(&Heartbeat{
		Type:     OptHeartbeatCheckoutSuccess,
		Feature:  C.GoString(feature),
		Licenses: int(n),
	})
}

//export goCheckoutRetry
func goCheckoutRetry(p unsafe.Pointer, feature *C.char, n C.int) {
	(*cgoClient)(p).heartbeat(&Heartbeat{
		Type:     OptHeartbeatRetryFeature,
		Feature:  C.GoString(feature),
		Licenses: int(n),
	})
}

//export goExit
func goExit(p unsafe.Pointer) {
	(*cgoClient)(p).heartbeat(&Heartbeat{
		Type: OptHeartbeatExit,
	})
}

func lookupCallback(opt OptionType) unsafe.Pointer {
	switch opt {
	case OptHeartbeatCheckoutFailure:
		return unsafe.Pointer(C.goCheckoutFailure)
	case OptHeartbeatCheckoutSuccess:
		return unsafe.Pointer(C.goCheckoutSuccess)
	case OptHeartbeatRetryFeature:
		return unsafe.Pointer(C.goCheckoutRetry)
	case OptHeartbeatConnectionLost:
		return unsafe.Pointer(C.goConnectionLost)
	case OptHeartbeatExit:
		return unsafe.Pointer(C.goExit)
	default:
		return unsafe.Pointer(uintptr(0))
	}
}

type cgoClient struct {
	m         sync.Mutex // protects callbacks
	callbacks map[OptionType]HeartbeatFunc
	vendor    interface{}
	handle    C.LMX_HANDLE
}

var _ Client = (*cgoClient)(nil)

func newClient() (Client, error) {
	c := &cgoClient{}

	if s := Status(C.LMX_Init((*C.LMX_HANDLE)(unsafe.Pointer(&c.handle)))); s != StatSuccess {
		return nil, ToError(s)
	}
	return c, nil
}

func (c *cgoClient) Close() error {
	if c.handle == nil {
		return ErrNotInitialized
	}
	C.LMX_Free(c.handle)
	c.handle = nil
	return nil
}

func (c *cgoClient) Checkout(feature string, major, minor, count int) error {
	cfeature := C.CString(feature)
	s := Status(C.LMX_Checkout(c.handle, cfeature, C.int(major), C.int(minor), C.int(count)))
	C.free(unsafe.Pointer(cfeature))
	return ToError(s)
}

func (c *cgoClient) Checkin(feature string, count int) error {
	cfeature := C.CString(feature)
	s := Status(C.LMX_Checkin(c.handle, cfeature, C.int(count)))
	C.free(unsafe.Pointer(cfeature))
	return ToError(s)
}

func (c *cgoClient) Error() string {
	return C.GoString(C.LMX_GetErrorMessage(c.handle))
}

func (c *cgoClient) ErrorInfo() *ErrorInfo {
	cerr := C.LMX_GetError(c.handle)
	return &ErrorInfo{
		Status:      Status(cerr.LmxStat),
		Internal:    int(cerr.nInternal),
		Context:     int(cerr.nContext),
		Desc:        C.GoString(&cerr.szDescription[0]),
		FeatureName: C.GoString(&cerr.szFeatureName[0]),
	}
}

func (c *cgoClient) HostID(t HostIDType) ([]HostID, error) {
	backing := [C.LMX_MAX_HOSTIDS]C.LMX_HOSTID{}
	cids := (*C.LMX_HOSTID)(unsafe.Pointer(&backing[0]))
	clen := C.int(0)
	s := Status(C.LMX_Hostid(c.handle, C.LMX_HOSTID_TYPE(t), cids, &clen))
	if s != StatSuccess {
		return nil, ToError(s)
	}
	return goHostId(cids, int(clen)), nil
}

func (c *cgoClient) HostIDString(t HostIDType) (string, error) {
	backing := [C.LMX_MAX_LONG_STRING_LENGTH]C.char{}
	cids := (*C.char)(unsafe.Pointer(&backing[0]))
	s := Status(C.LMX_HostidSimple(c.handle, C.LMX_HOSTID_TYPE(t), cids))
	if s != StatSuccess {
		return "", ToError(s)
	}
	return C.GoString(cids), nil
}

func (c *cgoClient) LicenseInfo() ([]LicenseInfo, error) {
	var cLicInfo *C.LMX_LICENSE_INFO
	if s := Status(C.LMX_GetLicenseInfo(c.handle, &cLicInfo)); s != StatSuccess {
		return nil, ToError(s)
	}
	return goLicenseInfo(cLicInfo), nil
}

func (c *cgoClient) FeatureInfo(feature string) ([]FeatureInfo, error) {
	var cfi C.LMX_FEATURE_INFO
	cfeature := C.CString(feature)
	s := Status(C.LMX_GetFeatureInfo(c.handle, cfeature, &cfi))
	C.free(unsafe.Pointer(cfeature))
	if s != StatSuccess {
		return nil, ToError(s)
	}
	return goFeatureInfo(&cfi), nil
}

func (c *cgoClient) Heartbeat(feature string) error {
	cfeature := C.CString(feature)
	s := Status(C.LMX_Heartbeat(c.handle, cfeature))
	C.free(unsafe.Pointer(cfeature))
	return ToError(s)
}

func (c *cgoClient) ExpireTime(feature string) (t time.Duration, err error) {
	cfeature := C.CString(feature)
	ret := time.Duration(C.LMX_GetExpireTime(c.handle, cfeature))
	C.free(unsafe.Pointer(cfeature))
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

func (c *cgoClient) ServerLog(feature, message string) error {
	cfeature, cmessage := C.CString(feature), C.CString(message)
	s := Status(C.LMX_ServerLog(c.handle, cfeature, cmessage))
	C.free(unsafe.Pointer(cfeature))
	C.free(unsafe.Pointer(cmessage))
	return ToError(s)
}

func (c *cgoClient) ServerFunction(feature, message string) (string, error) {
	if len(message) >= int(C.LMX_MAX_LONG_STRING_LENGTH) {
		return "", ErrInvalidParameter
	}
	cfeature := C.CString(feature)
	backing := [C.LMX_MAX_LONG_STRING_LENGTH]C.char{}
	for i := 0; i < len(message); i++ {
		backing[i] = C.char(message[i])
	}
	cresponse := (*C.char)(unsafe.Pointer(&backing[0]))
	s := Status(C.LMX_ServerFunction(c.handle, cfeature, cresponse))
	C.free(unsafe.Pointer(cfeature))
	if s != StatSuccess {
		return "", ToError(s)
	}
	return C.GoString(cresponse), nil
}

func (c *cgoClient) ClientStoreSave(filename, content string) error {
	cfilename := C.CString(filename)
	ccontent := C.CString(content)
	s := Status(C.LMX_ClientStoreSave(c.handle, cfilename, ccontent))
	C.free(unsafe.Pointer(cfilename))
	C.free(unsafe.Pointer(ccontent))
	return ToError(s)
}

func (c *cgoClient) ClientStoreLoad(filename string) (string, error) {
	cfilename := C.CString(filename)
	backing := [C.LMX_MAX_LONG_STRING_LENGTH]C.char{}
	ccontent := (*C.char)(unsafe.Pointer(&backing[0]))
	s := Status(C.LMX_ClientStoreLoad(c.handle, cfilename, ccontent))
	C.free(unsafe.Pointer(cfilename))
	if s != StatSuccess {
		return "", ToError(s)
	}
	return C.GoString(ccontent), nil
}

func (c *cgoClient) SetOption(option OptionType, value interface{}) error {
	s := StatInvalidParameter
	switch option {
	case OptExactVersion, OptAllowBorrow, OptAllowGrace, OptTrialVirtualMachine,
		OptTrialTerminalServer, OptBlacklist, OptLicenseIdle, OptAllowMultipleServers,
		OptClientHostIDToServer, OptAllowCheckoutLessLicenses:
		s = c.setBooleanOpt(option, value)
	case OptLicensePath, OptCustomShareString, OptLicenseString, OptServersideRequestString,
		OptCustomUsername, OptCustomHostname, OptReservationToken, OptBindAddress:
		s = c.setStringOpt(option, value)
	case OptTrialDays, OptTrialUses, OptAutomaticHeartbeatAttempts, OptAutomaticHeartbeatInterval,
		OptHostIDCacheCleanupInterval:
		s = c.setIntegerOpt(option, value)
	case OptHostIDEnabled, OptHostIDDisabled:
		s = c.setHostIdTypeOpt(option, value)
	case OptCustomHostIDFunction:
		s = StatNotImplemented
	case OptHostIDCompareFunction:
		s = StatNotImplemented
	case OptHeartbeatCheckoutFailure, OptHeartbeatCheckoutSuccess, OptHeartbeatRetryFeature,
		OptHeartbeatConnectionLost, OptHeartbeatExit:
		s = c.setCallback(option, value)
	case OptHeartbeatCallbackVendordata:
		s = c.setVendordata(value)
	}
	return ToError(s)
}

func (c *cgoClient) heartbeat(h *Heartbeat) {
	c.m.Lock()
	if c.callbacks == nil {
		c.m.Unlock()
		return
	}
	v := c.vendor
	fn, ok := c.callbacks[h.Type]
	c.m.Unlock()
	if !ok {
		return
	}
	fn(v, h) // shouldn't be executed in critical section
}

func (c *cgoClient) setCallback(opt OptionType, val interface{}) Status {
	// Initialize heartbeat callbacks. It's effectively a nop if it's already
	// been initialized.
	if s := c.setVendordata(c.vendor); s != StatSuccess {
		return s
	}
	if val == nil {
		s := Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), unsafe.Pointer(uintptr(0))))
		if s != StatSuccess {
			return s
		}
		c.m.Lock()
		delete(c.callbacks, opt)
		c.m.Unlock()
		return StatSuccess
	}
	if fn, ok := val.(HeartbeatFunc); ok {
		c.m.Lock()
		defer c.m.Unlock()
		if _, ok = c.callbacks[opt]; !ok {
			s := Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), lookupCallback(opt)))
			if s != StatSuccess {
				return s
			}
		}
		c.callbacks[opt] = fn
		return StatSuccess
	}
	return StatInvalidParameter
}

func (c *cgoClient) setVendordata(val interface{}) Status {
	const opt = C.LMX_SETTINGS(OptHeartbeatCallbackVendordata)
	c.m.Lock()
	defer c.m.Unlock()
	if c.callbacks == nil {
		s := Status(C.LMX_SetOption(c.handle, opt, unsafe.Pointer(c)))
		if s != StatSuccess {
			return s
		}
		c.callbacks = make(map[OptionType]HeartbeatFunc)
	}
	c.vendor = val
	return StatSuccess
}

func (c *cgoClient) setBooleanOpt(opt OptionType, val interface{}) Status {
	var p unsafe.Pointer
	if val != nil {
		switch value := val.(type) {
		case int:
			if value == 1 {
				p = unsafe.Pointer(uintptr(1))
			}
		case bool:
			if value {
				p = unsafe.Pointer(uintptr(1))
			}
		default:
			return StatInvalidParameter
		}
	}
	return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), p))
}

func (c *cgoClient) setStringOpt(opt OptionType, val interface{}) Status {
	var p unsafe.Pointer
	if val != nil {
		s, ok := val.(string)
		if !ok {
			return StatInvalidParameter
		}
		p = unsafe.Pointer(C.CString(s))
		defer C.free(unsafe.Pointer(p))
	}
	return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), p))
}

func (c *cgoClient) setIntegerOpt(opt OptionType, val interface{}) Status {
	var p unsafe.Pointer
	if val != nil {
		n, ok := val.(int)
		if !ok {
			return StatInvalidParameter
		}
		p = unsafe.Pointer(uintptr(n))
	}
	return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), p))
}

func (c *cgoClient) setHostIdTypeOpt(opt OptionType, val interface{}) Status {
	if valHostID, ok := val.(HostIDType); ok {
		return Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), unsafe.Pointer(uintptr(valHostID))))
	}
	return StatInvalidParameter
}
