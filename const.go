package lmx

// #include <lmx.h>
import "C"

const Version = int(C.LMX_VERSION)

const (
	MaxVeryLongStringLength = int(C.LMX_MAX_VERY_LONG_STRING_LENGTH)
	MaxFieldLength          = int(C.LMX_MAX_FIELD_LENGTH)
	MinNameLength           = int(C.LMX_MIN_NAME_LENGTH)
	MaxNameLength           = int(C.LMX_MAX_NAME_LENGTH)
	MinCount                = int(C.LMX_MIN_COUNT)
	MaxCount                = int(C.LMX_MAX_COUNT)
	MinVersion              = int(C.LMX_MAX_VERSION)
	MaxVersion              = int(C.LMX_MIN_VERSION)
	MaxHostIDs              = int(C.LMX_MAX_HOSTIDS)
	MaxBorrowHours          = int(C.LMX_MAX_BORROW_HOURS)
	MaxGraceHours           = int(C.LMX_MAX_GRACE_HOURS)
	MaxTokenDependencies    = int(C.LMX_MAX_TOKEN_DEPENDENCIES)
	MaxTokenLoops           = int(C.LMX_MAX_TOKEN_LOOPS)
)

const (
	NoVersionRestriction = int(C.LMX_NO_VERSION_RESTRICTION)
	UnlimitedCount       = int(C.LMX_UNLIMITED_COUNT)
	LogicalCPUCount      = int(C.LMX_LOGICAL_CPU_COUNT)
	PhysicalCPUCount     = int(C.LMX_PHYSICAL_CPU_COUNT)
	AllLicenses          = int(C.LMX_ALL_LICENSES)
	AllFeatures          = string(C.LMX_ALL_FEATURES)
)
