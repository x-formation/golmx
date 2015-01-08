package lmx_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/x-formation/golmx"
)

type Random struct{ rand.Source }

var DefaultRandom = &Random{rand.NewSource(time.Now().Unix())}

func (r *Random) String(length int) string {
	random := make([]byte, length)
	for i := 0; i < length; i++ {
		random[i] = byte(r.Source.Int63()%0x5E + 0x20)
	}
	return string(random)
}

func TestLmx(t *testing.T) {
	c, err := lmx.NewClient()
	if err != nil {
		t.Errorf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
}

func TestHostID(t *testing.T) {
	checkHostID := func(t *testing.T, ID *lmx.HostID) {
		switch ID.Type {
		case lmx.HostIDAll, lmx.HostIDUnknown:
			t.Error("expected .Type to be not of a lmx.HostIDAll or lmx.HostIDUnknown kind")
		}
		if len(ID.Value) == 0 {
			t.Error("expected .Value to be non-empty")
		}
		if len(ID.Desc) == 0 {
			t.Error("expected .Desc to be non-empty")
		}

	}
	c, err := lmx.NewClient()
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
	hostID, err := c.HostID(lmx.HostIDAll)
	if err != nil {
		t.Fatal(`expected err==nil, got "%v"`, err)
	}
	if len(hostID) == 0 {
		t.Fatal("expected hostID to be non-empty")
	}
	for _, ID := range hostID {
		hostID, err := c.HostID(ID.Type)
		if err != nil {
			t.Fatal(`expected err==nil, got "%v"`, err)
		}
		checkHostID(t, &ID)
		found := false
		for _, ID2 := range hostID {
			checkHostID(t, &ID2)
			if ID2.Type != ID.Type {
				t.Error(`expected "%v"=="%v"`, ID2.Type, ID.Type)
			}
			if ID2 == ID {
				found = true
			}
		}
		if !found {
			t.Error(`expected to find the ID in a result of c.HostID(ID.Type)`)
		}
	}
}

func TestHostIDString(t *testing.T) {
	c, err := lmx.NewClient()
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
	hostID, err := c.HostIDString(lmx.HostIDAll)
	if err != nil {
		t.Fatal(`expected err==nil, got "%v"`, err)
	}
	if len(hostID) == 0 {
		t.Fatal("expected hostID to be non-empty")
	}
}

func TestCheckout(t *testing.T) {
	c, err := lmx.NewClient()
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
	err = c.SetOption(lmx.OptLicensePath, filepath.Join(
		os.Getenv("LMXROOT"), "examples", "licenses", "nodelocked.lic"))
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	if err = c.Checkout("f2", 2, 0, 1); err != lmx.ToError(lmx.StatBadVersion) {
		t.Errorf(`expected err=='%v', got '%v'`,
			lmx.ToError(lmx.StatBadVersion), err)
	}
	if err = c.Checkout("f2", 1, 0, 1); err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
}

func TestNodelockedFeatureInfo(t *testing.T) {
	c, err := lmx.NewClient()
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
	err = c.SetOption(lmx.OptLicensePath, filepath.Join(
		os.Getenv("LMXROOT"), "examples", "licenses", "nodelocked.lic"))
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	if err = c.Checkout("f2", 1, 0, 1); err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	fi, err := c.FeatureInfo("f2")
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	if len(fi) != 1 {
		t.Fatalf(`expected len(feature info)==1, got "%d"`, len(fi))
	}
	if fi[0].Name != "f2" {
		t.Errorf(`expected feature .Name=="f2", got "%s"`, fi[0].Name)
	}
	// TODO(rjeczalik): fix hardcoded vendor
	if fi[0].Vendor != "XFORMATION" {
		t.Errorf(`expected feature .Vendor=="XFORMATION", got "%s"`, fi[0].Vendor)
	}
	comment := "Put additional sub-licensing rules here"
	if !strings.Contains(fi[0].Comment, comment) {
		t.Errorf(`expected feature .Comment=="%s", got "%s"`, comment, fi[0].Comment)
	}
	if fi[0].Version.Major != 1 || fi[0].Version.Minor != 5 {
		t.Errorf(`expected feature .Version=={1,5}, got "{%d,%d}"`,
			fi[0].Version.Major, fi[0].Version.Minor)
	}
	expTime := time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC)
	if fi[0].Expiration.End != expTime {
		t.Errorf(`expected feature .Expiration.End=="%v", got "%v"`,
			fi[0].Expiration.End, expTime)
	}
	expTime = time.Date(2018, time.January, 1, 23, 59, 0, 0, time.UTC)
	if fi[0].Expire != expTime {
		t.Errorf(`expected feature .Expire=="%v", got "%v"`, fi[0].Expire, expTime)
	}
	if len(fi[0].Share) != 1 || fi[0].Share[0] != lmx.ShareVirtual {
		t.Errorf(`expected feature .Share==ShareVirtual, got "%v"`, fi[0].Share)
	}
}

func BenchmarkClientStore(b *testing.B) {
	c, err := lmx.NewClient()
	if err != nil {
		b.Fatalf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
	for i := 0; i < b.N; i++ {
		filename, content := "benchmark_"+DefaultRandom.String(32), DefaultRandom.String(512)
		b.SetBytes(2 * int64(len(content)))
		if err := c.ClientStoreSave(filename, content); err != nil {
			b.Fatalf(`ClientStoreSave failed writing "%v" to "%v"`, []byte(filename), []byte(content))
		}
		if readContent, err := c.ClientStoreLoad(filename); err != nil {
			b.Fatalf(`ClientStoreSave failed reading "%v"`, []byte(filename))
		} else if readContent != content {
			b.Fatalf(`expected read data to be "%v", got "%v"`, []byte(readContent), []byte(content))
		}
		if err := c.ClientStoreSave(filename, ""); err != nil {
			b.Fatal(`expected err==nil, got "%v"`, err)
		}
		if _, err := c.ClientStoreLoad(filename); err != lmx.ToError(lmx.StatFileReadFailure) {
			b.Fatal(`expected err==ErrFileReadFailure, got "%v"`, err)
		}
	}
}

type setOptTest struct {
	Opt  lmx.OptionType
	Val  []interface{}
	Stat lmx.Status
}

func TestSetOptionValue(t *testing.T) {
	c, err := lmx.NewClient()
	defer c.Close()
	if err != nil {
		t.Fatalf(`expected err to be nil, got "%v" (%v) instead`, err, c.Error())
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
			if got, want := c.SetOption(soTest.Opt, val), lmx.ToError(soTest.Stat); got != want {
				t.Errorf(`want "%v"; got "%v" (i=%d, val=%#v)`, want, got, i, val)
			}
		}
	}
}
