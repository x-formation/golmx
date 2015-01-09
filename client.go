package lmx

import "time"

// ToStatus TODO(rjeczalik)
func ToStatus(err error) Status {
	if err, ok := err.(*Error); ok {
		return err.Status
	}
	return StatUnknownFailure
}

// ToError TODO(rjeczalik)
func ToError(s Status) error {
	if err, ok := statusErrorMap[s]; ok {
		return err
	}
	return statusErrorMap[StatUnknownFailure]
}

// Error TODO(rjeczalik)
type Error struct {
	Status Status
	Err    error
}

// Error implements builtin error interface.
func (e *Error) Error() string {
	return e.Err.Error()
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

	// Err is the reason of failed checkout. It's underlying type is *Error.
	//
	// It has non-zero value only when the Type is HeartbeatCheckoutFailure.
	Err error

	// Failed tells how many heartbeats have failed prior to invoking the callback
	// (basically it must be equal to the value of OptAutomaticHeartbeatAttempts).
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

// Client TODO(rjeczalik)
type Client interface {
	// Close TODO(rjeczalik)
	Close() error

	// Checkout TODO(rjeczalik) LMX_Checkout
	Checkout(feature string, major, minor, patch int) error

	// Checkin TODO(rjeczalik) LMX_Checkin
	Checkin(feature string, count int) error

	// SetOption TODO(rjeczalik) LMX_SetOption
	SetOption(option OptionType, value interface{}) error

	// Error TODO(rjeczalik) LMX_GetErrorMessage
	Error() string

	// ErrorInfo TODO(rjeczalik) LMX_GetError
	//
	// TODO(rjeczalik): merge *ErrorInfo with *Error?
	ErrorInfo() *ErrorInfo

	// HostID TODO(rjeczalik) LMX_GetHostId
	HostID(hostid HostIDType) ([]HostID, error)

	// HostIDString TODO(rjeczalik) LMX_GetHostIdSimple
	HostIDString(hostid HostIDType) (string, error)

	// FeatureInfo TODO(rjeczalik) LMX_GetFeatureInfo
	FeatureInfo(feature string) ([]FeatureInfo, error)

	// LicenseInfo TODO(rjeczalik) LMX_GetLicenseInfo
	LicenseInfo() ([]LicenseInfo, error)

	// Heartbeat TODO(rjeczalik) LMX_Heartbeat
	Heartbeat(feature string) error

	// ExpireTime TODO(rjeczalik) LMX_GetExpireTime
	ExpireTime(feature string) (time.Duration, error)

	// ServerLog TODO(rjeczalik) LMX_GetServerLog
	ServerLog(feature string, message string) error

	// ServerFunction TODO(rjeczalik) LMX_ServerFunction
	ServerFunction(feature string, request string) (string, error)

	// ClientStoreSave TODO(rjeczalik) LMX_ClientStoreSave
	ClientStoreSave(filename string, content string) error

	// ClientStoreLoad TODO(rjeczalik) LMX_ClientStoreLoad
	ClientStoreLoad(filename string) (string, error)
}

// NewClient TODO(rjeczalik)
func NewClient() (Client, error) {
	return newClient()
}
