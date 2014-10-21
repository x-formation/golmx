package lmx_test

import (
	"testing"

	"git.int.x-formation.com/scm/dev/golmx.git"
)

type setOptTest struct {
	Opt  lmx.OptionType
	Val  []interface{}
	Stat lmx.Status
}

func TestSetOptionValue(t *testing.T) {
	c, err := lmx.NewClient()
	defer c.Close()
	if err != nil {
		t.Fatalf(`expected err to be nil, got "%v" (%v) instead`, err, c.GetErrorMessage())
	}

	booleanValid := []interface{}{0, 1, true, false, nil}
	integerValid := []interface{}{30, 60, nil}
	stringValid := []interface{}{"feature", "", nil}
	hostIDValid := []interface{}{
		lmx.HostIDEthernet, lmx.HostIDUsername, lmx.HostIDHostname,
		lmx.HostIDIPAddress, lmx.HostIDCustom, lmx.HostIDDongleHaspHL,
		lmx.HostIDHardDisk, lmx.HostIDLong, lmx.HostIDBios,
		lmx.HostIDWinProduct, lmx.HostIDAWSInstance, lmx.HostIDUnknown,
		lmx.HostIDAll,
	}

	tests := []setOptTest{
		// boolean options.
		{lmx.OptExactVersion, booleanValid, lmx.StatSuccess},
		{lmx.OptExactVersion, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptAllowBorrow, booleanValid, lmx.StatSuccess},
		{lmx.OptAllowBorrow, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptAllowGrace, booleanValid, lmx.StatSuccess},
		{lmx.OptAllowGrace, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptTrialVirtualMachine, booleanValid, lmx.StatSuccess},
		{lmx.OptTrialVirtualMachine, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptTrialTerminalServer, booleanValid, lmx.StatSuccess},
		{lmx.OptTrialTerminalServer, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptBlacklist, booleanValid, lmx.StatSuccess},
		{lmx.OptBlacklist, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptLicenseIdle, booleanValid, lmx.StatSuccess},
		{lmx.OptLicenseIdle, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptAllowMultipleServers, booleanValid, lmx.StatSuccess},
		{lmx.OptAllowMultipleServers, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptClientHostIDToServer, booleanValid, lmx.StatSuccess},
		{lmx.OptClientHostIDToServer, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptAllowCheckoutLessLicenses, booleanValid, lmx.StatSuccess},
		{lmx.OptAllowCheckoutLessLicenses, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		// integer options.
		{lmx.OptTrialDays, integerValid, lmx.StatSuccess},
		{lmx.OptTrialDays, []interface{}{0}, lmx.StatSuccess},
		{lmx.OptTrialDays,
			[]interface{}{"invalid", false, 100, -2}, lmx.StatInvalidParameter},
		{lmx.OptTrialUses, integerValid, lmx.StatSuccess},
		{lmx.OptTrialUses,
			[]interface{}{"invalid", false, -2}, lmx.StatInvalidParameter},
		{lmx.OptAutomaticHeartbeatInterval, integerValid, lmx.StatSuccess},
		{lmx.OptAutomaticHeartbeatInterval,
			[]interface{}{"invalid", false, 23, 15*60 + 1}, lmx.StatInvalidParameter},
		{lmx.OptAutomaticHeartbeatAttempts, integerValid, lmx.StatSuccess},
		{lmx.OptAutomaticHeartbeatAttempts, []interface{}{"invalid"}, lmx.StatInvalidParameter},
		{lmx.OptHostIDCacheCleanupInterval, integerValid, lmx.StatSuccess},
		{lmx.OptTrialDays, []interface{}{0}, lmx.StatSuccess},
		{lmx.OptHostIDCacheCleanupInterval,
			[]interface{}{"invalid", false, -2}, lmx.StatInvalidParameter},
		// string options.
		{lmx.OptLicensePath, stringValid, lmx.StatSuccess},
		{lmx.OptLicensePath, []interface{}{0, false}, lmx.StatInvalidParameter},
		{lmx.OptLicenseString, stringValid, lmx.StatSuccess},
		{lmx.OptLicenseString, []interface{}{0, false}, lmx.StatInvalidParameter},
		{lmx.OptCustomShareString, stringValid, lmx.StatSuccess},
		{lmx.OptCustomShareString, []interface{}{0, false}, lmx.StatInvalidParameter},
		{lmx.OptServersideRequestString, stringValid, lmx.StatSuccess},
		{lmx.OptServersideRequestString, []interface{}{0, false}, lmx.StatInvalidParameter},
		{lmx.OptCustomUsername, // can be set only once.
			[]interface{}{"username"}, lmx.StatSuccess},
		{lmx.OptCustomUsername,
			[]interface{}{"hostname", "", nil, 0, false}, lmx.StatInvalidParameter},
		{lmx.OptCustomHostname, // can be set only once.
			[]interface{}{"hostname"}, lmx.StatSuccess},
		{lmx.OptCustomHostname,
			[]interface{}{"hostname", "", nil, 0, false}, lmx.StatInvalidParameter},
		{lmx.OptReservationToken, stringValid, lmx.StatSuccess},
		{lmx.OptReservationToken, []interface{}{0, false}, lmx.StatInvalidParameter},
		{lmx.OptBindAddress, stringValid, lmx.StatSuccess},
		{lmx.OptBindAddress, []interface{}{0, false}, lmx.StatInvalidParameter},
		// hostId enable/disable options.
		{lmx.OptHostIDEnabled, hostIDValid, lmx.StatSuccess},
		{lmx.OptHostIDEnabled, []interface{}{nil, 0, false}, lmx.StatInvalidParameter},
		{lmx.OptHostIDDisabled, hostIDValid, lmx.StatSuccess},
		{lmx.OptHostIDDisabled, []interface{}{nil, 0, false}, lmx.StatInvalidParameter},
	}

	optionTableTest(t, c, tests)
}

func optionTableTest(t *testing.T, c lmx.Client, tests []setOptTest) {
	for i, soTest := range tests {
		for _, val := range soTest.Val {
			if got, want := c.SetOption(soTest.Opt, val), lmx.LookupError(soTest.Stat); got != want {
				t.Errorf(`want "%v"; got "%v" (i=%d, val=%#v)`, want, got, i, val)
			}
		}
	}
}
