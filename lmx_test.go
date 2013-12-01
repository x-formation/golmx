package lmx_test

import (
	"math/rand"
	"testing"
	"time"

	"x-formation/lmx"
)

type Random struct {
	rand.Source
}

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
	defer c.Close()
	if err != nil {
		t.Errorf(`expected err to be nil, got "%v" (%v) instead`, err, c.GetErrorMessage())
	}
}

func OptionTableTest(t *testing.T, c lmx.Client, options []lmx.OptionType, values []interface{}) {
	for _, opt := range options {
		for _, val := range values {
			if err := c.SetOption(opt, val); err != nil {
				t.Errorf(`expected err to be nil, got "%v" (%v:%v) instead`, err, opt, val)
			}
		}
	}
}

func TestSetOptionParameterHandling(t *testing.T) {
	c, err := lmx.NewClient()
	defer c.Close()
	if err != nil {
		t.Fatalf(`expected err to be nil, got "%v" (%v) instead`, err, c.GetErrorMessage())
	}
	boolOpts := []lmx.OptionType{
		lmx.OptExactVersion,
		lmx.OptAllowBorrow,
		lmx.OptAllowGrace,
		lmx.OptTrialVirtualMachine,
		lmx.OptTrialTerminalServer,
		lmx.OptBlacklist,
		lmx.OptAllowMultipleServers,
		lmx.OptClientHostIDToServer,
	}
	boolVals := []interface{}{0, 1, true, false, nil}
	intOpts := []lmx.OptionType{
		lmx.OptTrialDays,
		// LMX-2285
		// lmx.OptTrialUses,
		// lmx.OptAutomaticHeartbeatInterval,
		lmx.OptAutomaticHeartbeatAttempts,
		lmx.OptLicenseIdle,
		lmx.OptHostIDCacheCleanupInterval,
	}
	intVals := []interface{}{30, 60, nil}
	stringOpts := []lmx.OptionType{
		lmx.OptLicensePath,
		lmx.OptLicenseString,
		lmx.OptCustomShareString,
		lmx.OptServersideRequestString,
		// LMX-2285
		// lmx.OptCustomUsername,
		// lmx.OptCustomHostname,
		lmx.OptReservationToken,
		lmx.OptBindAddress,
	}
	stringVals := []interface{}{"feature", "", nil}
	hostIDOpts := []lmx.OptionType{
		lmx.OptHostIDEnabled,
		lmx.OptHostIDDisabled,
	}
	hostIDVals := []interface{}{
		lmx.HostIDEthernet,
		lmx.HostIDUsername,
		lmx.HostIDHostname,
		lmx.HostIDIPAddress,
		lmx.HostIDCustom,
		lmx.HostIDDongleHaspHL,
		lmx.HostIDHardDisk,
		lmx.HostIDLong,
		lmx.HostIDBios,
		lmx.HostIDWinProduct,
		lmx.HostIDAWSInstance,
		lmx.HostIDUnknown,
		lmx.HostIDAll,
	}

	OptionTableTest(t, c, boolOpts, boolVals)
	OptionTableTest(t, c, intOpts, intVals)
	OptionTableTest(t, c, stringOpts, stringVals)
	OptionTableTest(t, c, hostIDOpts, hostIDVals)
}

func TestGetHostID(t *testing.T) {
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
	defer c.Close()
	if err != nil {
		t.Fatalf(`expected err to be nil, got "%v" (%v) instead`, err, c.GetErrorMessage())
	}
	hostID, err := c.GetHostID(lmx.HostIDAll)
	if err != nil {
		t.Fatal(`expected err to be nil, got "%v" instead`, err)
	}
	if len(hostID) == 0 {
		t.Fatal("expected hostID to be non-empty")
	}
	for _, ID := range hostID {
		hostID, err := c.GetHostID(ID.Type)
		if err != nil {
			t.Fatal(`expected err to be nil, got "%v" instead`, err)
		}
		checkHostID(t, &ID)
		found := false
		for _, ID2 := range hostID {
			checkHostID(t, &ID2)
			if ID2.Type != ID.Type {
				t.Error(`expected "%v" to be equal to "%v"`, ID2.Type, ID.Type)
			}
			if ID2 == ID {
				found = true
			}
		}
		if !found {
			t.Error(`expected to find the ID in a result of c.GetHostID(ID.Type)`)
		}
	}
}

func TestGetHostIDSimple(t *testing.T) {
	c, err := lmx.NewClient()
	defer c.Close()
	if err != nil {
		t.Fatalf(`expected err to be nil, got "%v" (%v) instead`, c.GetErrorMessage())
	}
	hostID, err := c.GetHostIDSimple(lmx.HostIDAll)
	if err != nil {
		t.Fatal(`expected err to be nil, got "%v" instead`, err)
	}
	if len(hostID) == 0 {
		t.Fatal("expected hostID to be non-empty")
	}
}

func BenchmarkClientStore(b *testing.B) {
	c, err := lmx.NewClient()
	defer c.Close()
	if err != nil {
		b.Fatalf(`expected err to be nil, got "%v" (%v) instead`, err, c.GetErrorMessage())
	}
	for i := 0; i < b.N; i++ {
		filename, content := "benchmark_"+DefaultRandom.String(32), DefaultRandom.String(512)
		b.SetBytes(2 * int64(len(content)))
		if err := c.ClientStoreSave(filename, content); err != nil {
			b.Fatalf(`ClientStoreSave failed writing "%v" to "%v"`, []byte(filename), []byte(content))
		}
		if readContent, err := c.ClientStoreLoad(filename); err != nil {
			b.Fatalf(`ClientStoreSave failed reading "%v"`, []byte(filename))
		} else if readContent != content {
			b.Fatalf(`expected read data to be "%v", got "%v" instead`, []byte(readContent), []byte(content))
		}
		if err := c.ClientStoreSave(filename, ""); err != nil {
			b.Fatal(`expected err to be nil, got "%v" instead`, err)
		}
		if _, err := c.ClientStoreLoad(filename); err != lmx.ErrFileReadFailure {
			b.Fatal(`expected err to be ErrFileReadFailure, got "%v" instead`, err)
		}
	}
}
