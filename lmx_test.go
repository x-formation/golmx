package lmx_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.int.x-formation.com/scm/dev/golmx.git"
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
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
	hostID, err := c.GetHostID(lmx.HostIDAll)
	if err != nil {
		t.Fatal(`expected err==nil, got "%v"`, err)
	}
	if len(hostID) == 0 {
		t.Fatal("expected hostID to be non-empty")
	}
	for _, ID := range hostID {
		hostID, err := c.GetHostID(ID.Type)
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
			t.Error(`expected to find the ID in a result of c.GetHostID(ID.Type)`)
		}
	}
}

func TestGetHostIDSimple(t *testing.T) {
	c, err := lmx.NewClient()
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	defer c.Close()
	hostID, err := c.GetHostIDSimple(lmx.HostIDAll)
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
	if err = c.Checkout("f2", 2, 0, 1); err != lmx.LookupError(lmx.StatBadVersion) {
		t.Errorf(`expected err=='%v', got '%v'`,
			lmx.LookupError(lmx.StatBadVersion), err)
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
	fi, err := c.GetFeatureInfo("f2")
	if err != nil {
		t.Fatalf(`expected err==nil, got "%v"`, err)
	}
	if len(fi) != 1 {
		t.Fatalf(`expected len(feature info)==1, got "%d"`, len(fi))
	}
	if fi[0].Name != "f2" {
		t.Errorf(`expected feature .Name=="f2", got "%s"`, fi[0].Name)
	}
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
		if _, err := c.ClientStoreLoad(filename); err != lmx.LookupError(lmx.StatFileReadFailure) {
			b.Fatal(`expected err==ErrFileReadFailure, got "%v"`, err)
		}
	}
}
