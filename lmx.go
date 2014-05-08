package lmx

/*
#include <string.h>
#include <lmx.h>
#cgo linux freebsd LDFLAGS: -llmxclient -ldl
LMX_HOSTID* AllocMaxHostID() {
	return (LMX_HOSTID*)(malloc(LMX_MAX_HOSTIDS * sizeof(LMX_HOSTID)));
}
char* AllocLongString() {
	return (char*)(malloc(LMX_MAX_LONG_STRING_LENGTH * sizeof(char)));
}
LMX_FEATURE_INFO* AllocFeatureInfo() {
	return (LMX_FEATURE_INFO*)(malloc(sizeof(LMX_FEATURE_INFO)));
}
*/
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
		return nil, LookupError(s)
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

	return LookupError(Status(
		C.LMX_Checkout(c.handle, cfeature, C.int(major), C.int(minor), C.int(count))))
}

// Checkin TODO(rjeczalik)
func (c *cgoClient) Checkin(feature string, count int) error {
	cfeature := C.CString(feature)
	defer C.free(unsafe.Pointer(cfeature))

	return LookupError(Status(
		C.LMX_Checkin(c.handle, cfeature, C.int(count))))
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
	if s := Status(C.LMX_Hostid(c.handle, C.LMX_HOSTID_TYPE(t), cids, &clen)); s != StatSuccess {
		return nil, LookupError(s)
	}

	return goHostId(cids, int(clen)), nil
}

// GetHostIDSimple TODO(rjeczalik)
func (c *cgoClient) GetHostIDSimple(t HostIDType) (string, error) {
	cids := C.AllocLongString()
	defer C.free(unsafe.Pointer(cids))
	if s := Status(C.LMX_HostidSimple(c.handle, C.LMX_HOSTID_TYPE(t), cids)); s != StatSuccess {
		return "", LookupError(s)
	}

	return C.GoString(cids), nil
}

// GetLicenseInfo TODO(rjeczalik)
func (c *cgoClient) GetLicenseInfo() ([]LicenseInfo, error) {
	var cLicInfo *C.LMX_LICENSE_INFO
	if s := Status(C.LMX_GetLicenseInfo(c.handle, &cLicInfo)); s != StatSuccess {
		return nil, LookupError(s)
	}

	return goLicenseInfo(cLicInfo), nil
}

// GetFeatureInfo TODO(rjeczalik)
func (c *cgoClient) GetFeatureInfo(feature string) ([]FeatureInfo, error) {
	cfi := C.AllocFeatureInfo()
	defer C.free(unsafe.Pointer(cfi))
	cfeature := C.CString(feature)
	defer C.free(unsafe.Pointer(cfeature))
	if s := Status(C.LMX_GetFeatureInfo(c.handle, cfeature, cfi)); s != StatSuccess {
		return nil, LookupError(s)
	}

	return goFeatureInfo(cfi), nil
}

// Heartbeat TODO(rjeczalik)
func (c *cgoClient) Heartbeat(feature string) error {
	cfeature := C.CString(feature)
	defer C.free(unsafe.Pointer(cfeature))

	return LookupError(Status(C.LMX_Heartbeat(c.handle, cfeature)))
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
		err = LookupError(StatTooLateDate)
	case ret == -2:
		err = ErrDoesNotExpire
	case ret < -2:
		err = LookupError(StatUnknownFailure)
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

	return LookupError(Status(
		C.LMX_ServerLog(c.handle, cfeature, cmessage)))
}

// ServerFunction TODO(rjeczalik)
func (c *cgoClient) ServerFunction(feature, message string) (string, error) {
	if len(message) >= int(C.LMX_MAX_LONG_STRING_LENGTH) {
		return "", LookupError(StatInvalidParameter)
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
		return "", LookupError(s)
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

	return LookupError(Status(
		C.LMX_ClientStoreSave(c.handle, cfilename, ccontent)))
}

// ClientStoreLoad TODO(rjeczalik)
func (c *cgoClient) ClientStoreLoad(filename string) (string, error) {
	cfilename, ccontent := C.CString(filename), C.AllocLongString()
	defer func() {
		C.free(unsafe.Pointer(cfilename))
		C.free(unsafe.Pointer(ccontent))
	}()

	if s := Status(C.LMX_ClientStoreLoad(c.handle, cfilename, ccontent)); s != StatSuccess {
		return "", LookupError(s)
	}

	return C.GoString(ccontent), nil
}
