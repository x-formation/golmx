package lmx

// #include <lmx.h>
import "C"

import "errors"

// ErrDoesNotExpire TODO(rjeczalik)
var ErrDoesNotExpire = errors.New("the feature does not expire")

// ErrNotInitialized TODO(rjeczalik)
var ErrNotInitialized = errors.New("the LM-X client was not initialized")

var (
	ErrUnknownFailure           = newError(StatUnknownFailure)
	ErrInvalidParameter         = newError(StatInvalidParameter)
	ErrNoNetwork                = newError(StatNoNetwork)
	ErrBadLicfile               = newError(StatBadLicfile)
	ErrNoMemory                 = newError(StatNoMemory)
	ErrFileReadFailure          = newError(StatFileReadFailure)
	ErrBadDate                  = newError(StatBadDate)
	ErrBadKey                   = newError(StatBadKey)
	ErrFeatureNotFound          = newError(StatFeatureNotFound)
	ErrBadHostid                = newError(StatBadHostid)
	ErrTooEarlyDate             = newError(StatTooEarlyDate)
	ErrTooLateDate              = newError(StatTooLateDate)
	ErrBadVersion               = newError(StatBadVersion)
	ErrNetworkFailure           = newError(StatNetworkFailure)
	ErrNoNetworkHost            = newError(StatNoNetworkHost)
	ErrNetworkDeny              = newError(StatNetworkDeny)
	ErrNotEnoughLicenses        = newError(StatNotEnoughLicenses)
	ErrBadSystemclock           = newError(StatBadSystemclock)
	ErrTsDeny                   = newError(StatTsDeny)
	ErrVirtualDeny              = newError(StatVirtualDeny)
	ErrBorrowTooLong            = newError(StatBorrowTooLong)
	ErrFileSaveFailure          = newError(StatFileSaveFailure)
	ErrAlreadyBorrowed          = newError(StatAlreadyBorrowed)
	ErrBorrowReturnFailure      = newError(StatBorrowReturnFailure)
	ErrServerBorrowFailure      = newError(StatServerBorrowFailure)
	ErrBorrowNotEnabled         = newError(StatBorrowNotEnabled)
	ErrNotBorrowed              = newError(StatNotBorrowed)
	ErrDongleFailure            = newError(StatDongleFailure)
	ErrSoftlimit                = newError(StatSoftlimit)
	ErrBadPlatform              = newError(StatBadPlatform)
	ErrTokenLoop                = newError(StatTokenLoop)
	ErrBlacklist                = newError(StatBlacklist)
	ErrVendorDeny               = newError(StatVendorDeny)
	ErrNotNetworkFeature        = newError(StatNotNetworkFeature)
	ErrBadTimezone              = newError(StatBadTimezone)
	ErrServerNotInUse           = newError(StatServerNotInUse)
	ErrNotImplemented           = newError(StatNotImplemented)
	ErrBorrowLimitExceeded      = newError(StatBorrowLimitExceeded)
	ErrServerFuncFailure        = newError(StatServerFuncFailure)
	ErrHeartbeatLostLicense     = newError(StatHeartbeatLostLicense)
	ErrSingleLock               = newError(StatSingleLock)
	ErrAuthFailure              = newError(StatAuthFailure)
	ErrNetworkSendFailure       = newError(StatNetworkSendFailure)
	ErrNetworkReceiveFailure    = newError(StatNetworkReceiveFailure)
	ErrQueue                    = newError(StatQueue)
	ErrBadSecurityConfig        = newError(StatBadSecurityConfig)
	ErrFeatureHalMismatch       = newError(StatFeatureHalMismatch)
	ErrNotLocalFeature          = newError(StatNotLocalFeature)
	ErrFeatureNotReplaceable    = newError(StatFeatureNotReplaceable)
	ErrHostidNotAvailable       = newError(StatHostidNotAvailable)
	ErrFeatureAlreadyReserved   = newError(StatFeatureAlreadyReserved)
	ErrFeatureAlreadyCheckedOut = newError(StatFeatureAlreadyCheckedOut)
	ErrReservationNotFound      = newError(StatReservationNotFound)
	ErrApiNotReentrant          = newError(StatApiNotReentrant)
	ErrLicenseUploadFailure     = newError(StatLicenseUploadFailure)
	ErrInternalLicNotEmbedded   = newError(StatInternalLicNotEmbedded)
)

// Status TODO(rjeczalik)
type Status uint8

const (
	StatSuccess                  = Status(C.LMX_SUCCESS)
	StatUnknownFailure           = Status(C.LMX_UNKNOWN_ERROR)
	StatInvalidParameter         = Status(C.LMX_INVALID_PARAMETER)
	StatNoNetwork                = Status(C.LMX_NO_NETWORK)
	StatBadLicfile               = Status(C.LMX_BAD_LICFILE)
	StatNoMemory                 = Status(C.LMX_NO_MEMORY)
	StatFileReadFailure          = Status(C.LMX_FILE_READ_ERROR)
	StatBadDate                  = Status(C.LMX_BAD_DATE)
	StatBadKey                   = Status(C.LMX_BAD_KEY)
	StatFeatureNotFound          = Status(C.LMX_FEATURE_NOT_FOUND)
	StatBadHostid                = Status(C.LMX_BAD_HOSTID)
	StatTooEarlyDate             = Status(C.LMX_TOO_EARLY_DATE)
	StatTooLateDate              = Status(C.LMX_TOO_LATE_DATE)
	StatBadVersion               = Status(C.LMX_BAD_VERSION)
	StatNetworkFailure           = Status(C.LMX_NETWORK_ERROR)
	StatNoNetworkHost            = Status(C.LMX_NO_NETWORK_HOST)
	StatNetworkDeny              = Status(C.LMX_NETWORK_DENY)
	StatNotEnoughLicenses        = Status(C.LMX_NOT_ENOUGH_LICENSES)
	StatBadSystemclock           = Status(C.LMX_BAD_SYSTEMCLOCK)
	StatTsDeny                   = Status(C.LMX_TS_DENY)
	StatVirtualDeny              = Status(C.LMX_VIRTUAL_DENY)
	StatBorrowTooLong            = Status(C.LMX_BORROW_TOO_LONG)
	StatFileSaveFailure          = Status(C.LMX_FILE_SAVE_ERROR)
	StatAlreadyBorrowed          = Status(C.LMX_ALREADY_BORROWED)
	StatBorrowReturnFailure      = Status(C.LMX_BORROW_RETURN_ERROR)
	StatServerBorrowFailure      = Status(C.LMX_SERVER_BORROW_ERROR)
	StatBorrowNotEnabled         = Status(C.LMX_BORROW_NOT_ENABLED)
	StatNotBorrowed              = Status(C.LMX_NOT_BORROWED)
	StatDongleFailure            = Status(C.LMX_DONGLE_ERROR)
	StatSoftlimit                = Status(C.LMX_SOFTLIMIT)
	StatBadPlatform              = Status(C.LMX_BAD_PLATFORM)
	StatTokenLoop                = Status(C.LMX_TOKEN_LOOP)
	StatBlacklist                = Status(C.LMX_BLACKLIST)
	StatVendorDeny               = Status(C.LMX_VENDOR_DENY)
	StatNotNetworkFeature        = Status(C.LMX_NOT_NETWORK_FEATURE)
	StatBadTimezone              = Status(C.LMX_BAD_TIMEZONE)
	StatServerNotInUse           = Status(C.LMX_SERVER_NOT_IN_USE)
	StatNotImplemented           = Status(C.LMX_NOT_IMPLEMENTED)
	StatBorrowLimitExceeded      = Status(C.LMX_BORROW_LIMIT_EXCEEDED)
	StatServerFuncFailure        = Status(C.LMX_SERVER_FUNC_ERROR)
	StatHeartbeatLostLicense     = Status(C.LMX_HEARTBEAT_LOST_LICENSE)
	StatSingleLock               = Status(C.LMX_SINGLE_LOCK)
	StatAuthFailure              = Status(C.LMX_AUTH_ERROR)
	StatNetworkSendFailure       = Status(C.LMX_NETWORK_SEND_ERROR)
	StatNetworkReceiveFailure    = Status(C.LMX_NETWORK_RECEIVE_ERROR)
	StatQueue                    = Status(C.LMX_QUEUE)
	StatBadSecurityConfig        = Status(C.LMX_BAD_SECURITY_CONFIG)
	StatFeatureHalMismatch       = Status(C.LMX_FEATURE_HAL_MISMATCH)
	StatNotLocalFeature          = Status(C.LMX_NOT_LOCAL_FEATURE)
	StatFeatureNotReplaceable    = Status(C.LMX_FEATURE_NOT_REPLACEABLE)
	StatHostidNotAvailable       = Status(C.LMX_HOSTID_NOT_AVAILABLE)
	StatFeatureAlreadyReserved   = Status(C.LMX_FEATURE_ALREADY_RESERVED)
	StatFeatureAlreadyCheckedOut = Status(C.LMX_FEATURE_ALREADY_CHECKED_OUT)
	StatReservationNotFound      = Status(C.LMX_RESERVATION_NOT_FOUND)
	StatApiNotReentrant          = Status(C.LMX_API_NOT_REENTRANT)
	StatLicenseUploadFailure     = Status(C.LMX_LICENSE_UPLOAD_ERROR)
	StatInternalLicNotEmbedded   = Status(C.LMX_INTERNAL_LICENSE_NOT_EMBEDDED)
	StatSystemInterprocess       = Status(C.LMX_SYSTEM_INTERPROCESS)
	StatQueueInsert              = Status(C.LMX_QUEUE_INSERT)
	StatQueueAlready             = Status(C.LMX_QUEUE_ALREADY)
	StatReserve                  = Status(C.LMX_RESERVE)
	StatLimit                    = Status(C.LMX_LIMIT)
)

func newError(s Status) error {
	return &Error{
		Status: s,
		Err:    errors.New(C.GoString(C.LMX_GetErrorMessageSimple(C.LMX_STATUS(s)))),
	}
}

var statusErrorMap = map[Status]error{
	StatSuccess:                  nil,
	StatUnknownFailure:           ErrUnknownFailure,
	StatInvalidParameter:         ErrInvalidParameter,
	StatNoNetwork:                ErrNoNetwork,
	StatBadLicfile:               ErrBadLicfile,
	StatNoMemory:                 ErrNoMemory,
	StatFileReadFailure:          ErrFileReadFailure,
	StatBadDate:                  ErrBadDate,
	StatBadKey:                   ErrBadKey,
	StatFeatureNotFound:          ErrFeatureNotFound,
	StatBadHostid:                ErrBadHostid,
	StatTooEarlyDate:             ErrTooEarlyDate,
	StatTooLateDate:              ErrTooLateDate,
	StatBadVersion:               ErrBadVersion,
	StatNetworkFailure:           ErrNetworkFailure,
	StatNoNetworkHost:            ErrNoNetworkHost,
	StatNetworkDeny:              ErrNetworkDeny,
	StatNotEnoughLicenses:        ErrNotEnoughLicenses,
	StatBadSystemclock:           ErrBadSystemclock,
	StatTsDeny:                   ErrTsDeny,
	StatVirtualDeny:              ErrVirtualDeny,
	StatBorrowTooLong:            ErrBorrowTooLong,
	StatFileSaveFailure:          ErrFileSaveFailure,
	StatAlreadyBorrowed:          ErrAlreadyBorrowed,
	StatBorrowReturnFailure:      ErrBorrowReturnFailure,
	StatServerBorrowFailure:      ErrServerBorrowFailure,
	StatBorrowNotEnabled:         ErrBorrowNotEnabled,
	StatNotBorrowed:              ErrNotBorrowed,
	StatDongleFailure:            ErrDongleFailure,
	StatSoftlimit:                ErrSoftlimit,
	StatBadPlatform:              ErrBadPlatform,
	StatTokenLoop:                ErrTokenLoop,
	StatBlacklist:                ErrBlacklist,
	StatVendorDeny:               ErrVendorDeny,
	StatNotNetworkFeature:        ErrNotNetworkFeature,
	StatBadTimezone:              ErrBadTimezone,
	StatServerNotInUse:           ErrServerNotInUse,
	StatNotImplemented:           ErrNotImplemented,
	StatBorrowLimitExceeded:      ErrBorrowLimitExceeded,
	StatServerFuncFailure:        ErrServerFuncFailure,
	StatHeartbeatLostLicense:     ErrHeartbeatLostLicense,
	StatSingleLock:               ErrSingleLock,
	StatAuthFailure:              ErrAuthFailure,
	StatNetworkSendFailure:       ErrNetworkSendFailure,
	StatNetworkReceiveFailure:    ErrNetworkReceiveFailure,
	StatQueue:                    ErrQueue,
	StatBadSecurityConfig:        ErrBadSecurityConfig,
	StatFeatureHalMismatch:       ErrFeatureHalMismatch,
	StatNotLocalFeature:          ErrNotLocalFeature,
	StatFeatureNotReplaceable:    ErrFeatureNotReplaceable,
	StatHostidNotAvailable:       ErrHostidNotAvailable,
	StatFeatureAlreadyReserved:   ErrFeatureAlreadyReserved,
	StatFeatureAlreadyCheckedOut: ErrFeatureAlreadyCheckedOut,
	StatReservationNotFound:      ErrReservationNotFound,
	StatApiNotReentrant:          ErrApiNotReentrant,
	StatLicenseUploadFailure:     ErrLicenseUploadFailure,
	StatInternalLicNotEmbedded:   ErrInternalLicNotEmbedded,
}
