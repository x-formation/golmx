package lmx

// #include <lmx.h>
import "C"

// LicenseType
type LicenseType uint8

const (
	LicenseLocal   = LicenseType(C.LMX_TYPE_LOCAL)
	LicenseNetwork = LicenseType(C.LMX_TYPE_NETWORK)
	LicenseBorrow  = LicenseType(C.LMX_TYPE_BORROW)
	LicenseGrace   = LicenseType(C.LMX_TYPE_GRACE)
	LicenseTrial   = LicenseType(C.LMX_TYPE_TRIAL)
)

// ShareType
type ShareType uint16

const (
	ShareNone     = ShareType(C.LMX_SHARE_NONE)
	ShareHost     = ShareType(C.LMX_SHARE_HOST)
	ShareUser     = ShareType(C.LMX_SHARE_USER)
	ShareCustom   = ShareType(C.LMX_SHARE_CUSTOM)
	ShareTerminal = ShareType(C.LMX_SHARE_TS)
	ShareVirtual  = ShareType(C.LMX_SHARE_VIRTUAL)
	ShareSingle   = ShareType(C.LMX_SHARE_SINGLE)
)

// KeyType
type KeyType uint8

const (
	KeyExclusive = KeyType(C.LMX_KEYTYPE_EXCLUSIVE)
	KeyAdditive  = KeyType(C.LMX_KEYTYPE_ADDITIVE)
	KeyToken     = KeyType(C.LMX_KEYTYPE_TOKEN)
	KeyUnknown   = KeyType(C.LMX_KEYTYPE_UNKNOWN)
)

// HostIDType
type HostIDType uint8

const (
	HostIDEthernet     = HostIDType(C.LMX_HOSTID_ETHERNET)
	HostIDUsername     = HostIDType(C.LMX_HOSTID_USERNAME)
	HostIDHostname     = HostIDType(C.LMX_HOSTID_HOSTNAME)
	HostIDIPAddress    = HostIDType(C.LMX_HOSTID_IPADDRESS)
	HostIDCustom       = HostIDType(C.LMX_HOSTID_CUSTOM)
	HostIDDongleHaspHL = HostIDType(C.LMX_HOSTID_DONGLE_HASPHL)
	HostIDHardDisk     = HostIDType(C.LMX_HOSTID_HARDDISK)
	HostIDLong         = HostIDType(C.LMX_HOSTID_LONG)
	HostIDBios         = HostIDType(C.LMX_HOSTID_BIOS)
	HostIDWinProduct   = HostIDType(C.LMX_HOSTID_WIN_PRODUCT_ID)
	HostIDAWSInstance  = HostIDType(C.LMX_HOSTID_AWS_INSTANCE_ID)
	HostIDUnknown      = HostIDType(C.LMX_HOSTID_UNKNOWN)
	HostIDAll          = HostIDType(C.LMX_HOSTID_ALL)
)

// HostIDKeyType
type HostIDKeyType uint8

const (
	HostIDKeyClient = HostIDKeyType(C.LMX_CLIENT_HOSTID)
	HostIDKeyServer = HostIDKeyType(C.LMX_SERVER_HOSTID)
)

// OptionType
type OptionType uint8

const (
	OptExactVersion                     = OptionType(C.LMX_OPT_EXACT_VERSION)
	OptLicensePath                      = OptionType(C.LMX_OPT_LICENSE_PATH)
	OptCustomHostIDFunction             = OptionType(C.LMX_OPT_CUSTOM_HOSTID_FUNCTION)
	OptHostIDCompareFunction            = OptionType(C.LMX_OPT_HOSTID_COMPARE_FUNCTION)
	OptHeartbeatCheckoutFailureFunction = OptionType(C.LMX_OPT_HEARTBEAT_CHECKOUT_FAILURE_FUNCTION)
	OptHeartbeatCheckoutSuccessFunction = OptionType(C.LMX_OPT_HEARTBEAT_CHECKOUT_SUCCESS_FUNCTION)
	OptRetryFeatureFunction             = OptionType(C.LMX_OPT_HEARTBEAT_RETRY_FEATURE_FUNCTION)
	OptHeartbeatConnectionLostFunction  = OptionType(C.LMX_OPT_HEARTBEAT_CONNECTION_LOST_FUNCTION)
	OptHeartbeatExitFunction            = OptionType(C.LMX_OPT_HEARTBEAT_EXIT_FUNCTION)
	OptHeartbeatCallbackVendordata      = OptionType(C.LMX_OPT_HEARTBEAT_CALLBACK_VENDORDATA)
	OptAllowBorrow                      = OptionType(C.LMX_OPT_ALLOW_BORROW)
	OptAllowGrace                       = OptionType(C.LMX_OPT_ALLOW_GRACE)
	OptTrialDays                        = OptionType(C.LMX_OPT_TRIAL_DAYS)
	OptTrialUses                        = OptionType(C.LMX_OPT_TRIAL_USES)
	OptTrialVirtualMachine              = OptionType(C.LMX_OPT_TRIAL_VIRTUAL_MACHINE)
	OptTrialTerminalServer              = OptionType(C.LMX_OPT_TRIAL_TERMINAL_SERVER)
	OptAutomaticHeartbeatAttempts       = OptionType(C.LMX_OPT_AUTOMATIC_HEARTBEAT_ATTEMPTS)
	OptAutomaticHeartbeatInterval       = OptionType(C.LMX_OPT_AUTOMATIC_HEARTBEAT_INTERVAL)
	OptCustomShareString                = OptionType(C.LMX_OPT_CUSTOM_SHARE_STRING)
	OptLicenseString                    = OptionType(C.LMX_OPT_LICENSE_STRING)
	OptServersideRequestString          = OptionType(C.LMX_OPT_SERVERSIDE_REQUEST_STRING)
	OptLicenseIdle                      = OptionType(C.LMX_OPT_LICENSE_IDLE)
	OptCustomUsername                   = OptionType(C.LMX_OPT_CUSTOM_USERNAME)
	OptCustomHostname                   = OptionType(C.LMX_OPT_CUSTOM_HOSTNAME)
	OptBlacklist                        = OptionType(C.LMX_OPT_BLACKLIST)
	OptAllowMultipleServers             = OptionType(C.LMX_OPT_ALLOW_MULTIPLE_SERVERS)
	OptHostIDCacheCleanupInterval       = OptionType(C.LMX_OPT_HOSTID_CACHE_CLEANUP_INTERVAL)
	OptReservationToken                 = OptionType(C.LMX_OPT_RESERVATION_TOKEN)
	OptBindAddress                      = OptionType(C.LMX_OPT_BIND_ADDRESS)
	OptClientHostIDToServer             = OptionType(C.LMX_OPT_CLIENT_HOSTIDS_TO_SERVER)
	OptHostIDEnabled                    = OptionType(C.LMX_OPT_HOSTID_ENABLED)
	OptHostIDDisabled                   = OptionType(C.LMX_OPT_HOSTID_DISABLED)
)

// PlatformType
type PlatformType string

var (
	PlatformWin32_x86        = PlatformType("Win32_x86")
	PlatformWin64_x64        = PlatformType("Win64_x64")
	PlatformMacOSX_Universal = PlatformType("Macosx_Universal")
	PlatformLinux_x86        = PlatformType("Linux_x86")
	PlatformLinux_x86_64     = PlatformType("Linux_x64")
	PlatformLinux_arm        = PlatformType("Linux_arm")
	PlatformFreeBSD_x86_64   = PlatformType("FreeBSD_x64")
	PlatformSolaris_x86_64   = PlatformType("Solaris_x64")
	PlatformSolaris_sparc    = PlatformType("Solaris_sparc")
	PlatformSolaris_sparc64  = PlatformType("Solaris_sparc64")
	PlatformAIX_ppc          = PlatformType("Aix_ppc")
	PlatformAIX_ppc64        = PlatformType("Aix_ppc64")
	PlatformHPUX_ia64        = PlatformType("Hpux_ia64")
)

var PlatformAll = []*PlatformType{
	&PlatformWin32_x86,
	&PlatformWin64_x64,
	&PlatformMacOSX_Universal,
	&PlatformLinux_x86,
	&PlatformLinux_x86_64,
	&PlatformLinux_arm,
	&PlatformFreeBSD_x86_64,
	&PlatformSolaris_x86_64,
	&PlatformSolaris_sparc,
	&PlatformSolaris_sparc64,
	&PlatformAIX_ppc,
	&PlatformAIX_ppc64,
	&PlatformHPUX_ia64,
}

// ClockCheckType
type ClockCheckType []byte

var (
	ClockInternet = ClockCheckType("INTERNET")
	ClockLocal    = ClockCheckType("LOCAL")
	ClockDisabled = ClockCheckType("FALSE")
)
