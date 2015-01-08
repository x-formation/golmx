package lmx

// #include <lmx.h>
import "C"

import "errors"

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

var (
	// ErrDoesNotExpire TODO(rjeczalik)
	ErrDoesNotExpire = errors.New("The feature does not expire.")

	// ErrNotInitialized TODO(rjeczalik)
	ErrNotInitialized = errors.New("The LM-X client was not initialized.")
)

// Error TODO(rjeczalik)
type Error struct {
	Status Status
	Err    error
}

// Error implements builtin error interface.
func (e *Error) Error() string {
	return e.Err.Error()
}

func newError(s Status) error {
	return &Error{
		Status: s,
		Err:    errors.New(C.GoString(C.LMX_GetErrorMessageSimple(C.LMX_STATUS(s)))),
	}
}

var statusErrorMap = map[Status]error{
	StatSuccess:                  nil,
	StatUnknownFailure:           newError(StatUnknownFailure),
	StatInvalidParameter:         newError(StatInvalidParameter),
	StatNoNetwork:                newError(StatNoNetwork),
	StatBadLicfile:               newError(StatBadLicfile),
	StatNoMemory:                 newError(StatNoMemory),
	StatFileReadFailure:          newError(StatFileReadFailure),
	StatBadDate:                  newError(StatBadDate),
	StatBadKey:                   newError(StatBadKey),
	StatFeatureNotFound:          newError(StatFeatureNotFound),
	StatBadHostid:                newError(StatBadHostid),
	StatTooEarlyDate:             newError(StatTooEarlyDate),
	StatTooLateDate:              newError(StatTooLateDate),
	StatBadVersion:               newError(StatBadVersion),
	StatNetworkFailure:           newError(StatNetworkFailure),
	StatNoNetworkHost:            newError(StatNoNetworkHost),
	StatNetworkDeny:              newError(StatNetworkDeny),
	StatNotEnoughLicenses:        newError(StatNotEnoughLicenses),
	StatBadSystemclock:           newError(StatBadSystemclock),
	StatTsDeny:                   newError(StatTsDeny),
	StatVirtualDeny:              newError(StatVirtualDeny),
	StatBorrowTooLong:            newError(StatBorrowTooLong),
	StatFileSaveFailure:          newError(StatFileSaveFailure),
	StatAlreadyBorrowed:          newError(StatAlreadyBorrowed),
	StatBorrowReturnFailure:      newError(StatBorrowReturnFailure),
	StatServerBorrowFailure:      newError(StatServerBorrowFailure),
	StatBorrowNotEnabled:         newError(StatBorrowNotEnabled),
	StatNotBorrowed:              newError(StatNotBorrowed),
	StatDongleFailure:            newError(StatDongleFailure),
	StatSoftlimit:                newError(StatSoftlimit),
	StatBadPlatform:              newError(StatBadPlatform),
	StatTokenLoop:                newError(StatTokenLoop),
	StatBlacklist:                newError(StatBlacklist),
	StatVendorDeny:               newError(StatVendorDeny),
	StatNotNetworkFeature:        newError(StatNotNetworkFeature),
	StatBadTimezone:              newError(StatBadTimezone),
	StatServerNotInUse:           newError(StatServerNotInUse),
	StatNotImplemented:           newError(StatNotImplemented),
	StatBorrowLimitExceeded:      newError(StatBorrowLimitExceeded),
	StatServerFuncFailure:        newError(StatServerFuncFailure),
	StatHeartbeatLostLicense:     newError(StatHeartbeatLostLicense),
	StatSingleLock:               newError(StatSingleLock),
	StatAuthFailure:              newError(StatAuthFailure),
	StatNetworkSendFailure:       newError(StatNetworkSendFailure),
	StatNetworkReceiveFailure:    newError(StatNetworkReceiveFailure),
	StatQueue:                    newError(StatQueue),
	StatBadSecurityConfig:        newError(StatBadSecurityConfig),
	StatFeatureHalMismatch:       newError(StatFeatureHalMismatch),
	StatNotLocalFeature:          newError(StatNotLocalFeature),
	StatFeatureNotReplaceable:    newError(StatFeatureNotReplaceable),
	StatHostidNotAvailable:       newError(StatHostidNotAvailable),
	StatFeatureAlreadyReserved:   newError(StatFeatureAlreadyReserved),
	StatFeatureAlreadyCheckedOut: newError(StatFeatureAlreadyCheckedOut),
	StatReservationNotFound:      newError(StatReservationNotFound),
	StatApiNotReentrant:          newError(StatApiNotReentrant),
	StatLicenseUploadFailure:     newError(StatLicenseUploadFailure),
	StatInternalLicNotEmbedded:   newError(StatInternalLicNotEmbedded),
}

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
