package lmx

// #include <lmx.h>
import "C"

import "errors"

// Error
type Error struct {
	Stat Status
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func newError(s Status) error {
	return &Error{
		Stat: s,
		Err:  errors.New(C.GoString(C.LMX_GetErrorMessageSimple(C.LMX_STATUS(s)))),
	}
}

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
)

var (
	ErrDoesNotExpire  = errors.New("The feature does not expire.")
	ErrNotInitialized = errors.New("The LM-X client was not initialized.")
)

func lookupError(s Status) error {
	switch s {
	case StatSuccess:
		return nil
	case StatUnknownFailure:
		return ErrUnknownFailure
	case StatInvalidParameter:
		return ErrInvalidParameter
	case StatNoNetwork:
		return ErrNoNetwork
	case StatBadLicfile:
		return ErrBadLicfile
	case StatNoMemory:
		return ErrNoMemory
	case StatFileReadFailure:
		return ErrFileReadFailure
	case StatBadDate:
		return ErrBadDate
	case StatBadKey:
		return ErrBadKey
	case StatFeatureNotFound:
		return ErrFeatureNotFound
	case StatBadHostid:
		return ErrBadHostid
	case StatTooEarlyDate:
		return ErrTooEarlyDate
	case StatTooLateDate:
		return ErrTooLateDate
	case StatBadVersion:
		return ErrBadVersion
	case StatNetworkFailure:
		return ErrNetworkFailure
	case StatNoNetworkHost:
		return ErrNoNetworkHost
	case StatNetworkDeny:
		return ErrNetworkDeny
	case StatNotEnoughLicenses:
		return ErrNotEnoughLicenses
	case StatBadSystemclock:
		return ErrBadSystemclock
	case StatTsDeny:
		return ErrTsDeny
	case StatVirtualDeny:
		return ErrVirtualDeny
	case StatBorrowTooLong:
		return ErrBorrowTooLong
	case StatFileSaveFailure:
		return ErrFileSaveFailure
	case StatAlreadyBorrowed:
		return ErrAlreadyBorrowed
	case StatBorrowReturnFailure:
		return ErrBorrowReturnFailure
	case StatServerBorrowFailure:
		return ErrServerBorrowFailure
	case StatBorrowNotEnabled:
		return ErrBorrowNotEnabled
	case StatNotBorrowed:
		return ErrNotBorrowed
	case StatDongleFailure:
		return ErrDongleFailure
	case StatSoftlimit:
		return ErrSoftlimit
	case StatBadPlatform:
		return ErrBadPlatform
	case StatTokenLoop:
		return ErrTokenLoop
	case StatBlacklist:
		return ErrBlacklist
	case StatVendorDeny:
		return ErrVendorDeny
	case StatNotNetworkFeature:
		return ErrNotNetworkFeature
	case StatBadTimezone:
		return ErrBadTimezone
	case StatServerNotInUse:
		return ErrServerNotInUse
	case StatNotImplemented:
		return ErrNotImplemented
	case StatBorrowLimitExceeded:
		return ErrBorrowLimitExceeded
	case StatServerFuncFailure:
		return ErrServerFuncFailure
	case StatHeartbeatLostLicense:
		return ErrHeartbeatLostLicense
	case StatSingleLock:
		return ErrSingleLock
	case StatAuthFailure:
		return ErrAuthFailure
	case StatNetworkSendFailure:
		return ErrNetworkSendFailure
	case StatNetworkReceiveFailure:
		return ErrNetworkReceiveFailure
	case StatQueue:
		return ErrQueue
	case StatBadSecurityConfig:
		return ErrBadSecurityConfig
	case StatFeatureHalMismatch:
		return ErrFeatureHalMismatch
	case StatNotLocalFeature:
		return ErrNotLocalFeature
	case StatFeatureNotReplaceable:
		return ErrFeatureNotReplaceable
	case StatHostidNotAvailable:
		return ErrHostidNotAvailable
	case StatFeatureAlreadyReserved:
		return ErrFeatureAlreadyReserved
	case StatFeatureAlreadyCheckedOut:
		return ErrFeatureAlreadyCheckedOut
	case StatReservationNotFound:
		return ErrReservationNotFound
	case StatApiNotReentrant:
		return ErrApiNotReentrant
	case StatLicenseUploadFailure:
		return ErrLicenseUploadFailure
	}
	return ErrUnknownFailure
}
