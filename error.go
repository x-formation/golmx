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

var (
	ErrDoesNotExpire  = errors.New("The feature does not expire.")
	ErrNotInitialized = errors.New("The LM-X client was not initialized.")
)

func LookupError(s Status) error {
	if err, ok := statusErrorMap[s]; ok {
		return err
	}
	return statusErrorMap[StatUnknownFailure]
}
