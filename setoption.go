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
	"unsafe"
)

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
		Status:   Status(s),
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

func gocallback(opt OptionType) unsafe.Pointer {
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

// Heartbeat represents heartbeat context for callbacks.
type Heartbeat struct {
	// Type denotes the type of the heartbeat function called. It holds one of the
	// following values:
	//
	//   - OptHeartbeatRetryFeature
	//   - OptHeartbeatCheckoutFailure
	//   - OptHeartbeatCheckoutSuccess
	//   - OptHeartbeatConnectionLost
	//   - OptHeartbeatExit
	//
	// It is always non-zero.
	Type OptionType

	// Addr is license server's network addres of the form "host:port".
	//
	// It has non-zero value only when the Type is HeartbeatConnectionLost.
	Addr string

	// Features is the name of a feature that failed.
	//
	// It has non-zero value for all the callback types, except HeartbeatExit.
	Feature string

	// Status is the reason of failed checkout.
	//
	// It has non-zero value only when the Type is HeartbeatCheckoutFailure.
	Status Status

	// Failed tells how many heartbeats have failed prior to invoking the callback
	// (basically it must be equal to value of OptAutomaticHeartbeatAttempts).
	//
	// It has non-zero value only when the Type is HeartbeatConnectionLost.
	Heartbeats int

	// Licenses is the number of licenses.
	//
	// It has non-zero value only when the Type is one of HeartbeatCheckout* types.
	Licenses int
}

// HeartbeatFunc TODO(rjeczalik)
type HeartbeatFunc func(interface{}, *Heartbeat)

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

// SetOption TODO(rjeczalik)
func (c *cgoClient) SetOption(option OptionType, value interface{}) error {
	switch option {
	case OptExactVersion, OptAllowBorrow, OptAllowGrace, OptTrialVirtualMachine,
		OptTrialTerminalServer, OptBlacklist, OptLicenseIdle, OptAllowMultipleServers,
		OptClientHostIDToServer, OptAllowCheckoutLessLicenses:
		return LookupError(c.setBooleanOpt(option, value))
	case OptLicensePath, OptCustomShareString, OptLicenseString, OptServersideRequestString,
		OptCustomUsername, OptCustomHostname, OptReservationToken, OptBindAddress:
		return LookupError(c.setStringOpt(option, value))
	case OptTrialDays, OptTrialUses, OptAutomaticHeartbeatAttempts, OptAutomaticHeartbeatInterval,
		OptHostIDCacheCleanupInterval:
		return LookupError(c.setIntegerOpt(option, value))
	case OptHostIDEnabled, OptHostIDDisabled:
		return LookupError(c.setHostIdTypeOpt(option, value))
	case OptCustomHostIDFunction:
		return LookupError(StatNotImplemented)
	case OptHostIDCompareFunction:
		return LookupError(StatNotImplemented)
	case OptHeartbeatCheckoutFailure, OptHeartbeatCheckoutSuccess, OptHeartbeatRetryFeature,
		OptHeartbeatConnectionLost, OptHeartbeatExit:
		return LookupError(c.setCallback(option, value))
	case OptHeartbeatCallbackVendordata:
		return LookupError(c.setVendordata(option, value))
	default:
		return LookupError(StatInvalidParameter)
	}
	return LookupError(StatInvalidParameter)
}

func (c *cgoClient) setCallback(opt OptionType, val interface{}) Status {
	// Initialize heartbeat callbacks.It's effectively a nop if it's already
	// been initialized.
	if s := c.setCallback(OptHeartbeatCallbackVendordata, c.vendor); s != StatSuccess {
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
	}
	if fn, ok := val.(HeartbeatFunc); ok {
		c.m.Lock()
		defer c.m.Unlock()
		if _, ok = c.callbacks[opt]; !ok {
			s := Status(C.LMX_SetOption(c.handle, C.LMX_SETTINGS(opt), gocallback(opt)))
			if s != StatSuccess {
				return s
			}
		}
		c.callbacks[opt] = fn
		return StatSuccess
	}
	return StatInvalidParameter
}

func (c *cgoClient) setVendordata(_ OptionType, val interface{}) Status {
	c.m.Lock()
	defer c.m.Unlock()
	if c.callbacks == nil {
		s := Status(C.LMX_SetOption(c.handle, C.LMX_OPT_HEARTBEAT_CALLBACK_VENDORDATA,
			unsafe.Pointer(c)))
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
