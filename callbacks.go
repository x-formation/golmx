package lmx

/*
#include <lmx.h>
#include <stdio.h>
#include <unistd.h>

void* IntToPtr(int i) {
      return (void*)(uintptr_t)(i);
}

// The gateway functions
void HeartbeatConnectionLost_cgo(void *pVendorData, const char *szHost, int nPort, int nFailedHeartbeats)
{
        HeartbeatConnectionLostGo(pVendorData, szHost, nPort, nFailedHeartbeats);
}

void HeartbeatCheckoutFailure_cgo(void *pVendorData, const char *szFeatureName, int nUsedLicCount, LMX_STATUS LmxStat)
{
        HeartbeatCheckoutFailureGo(pVendorData, szFeatureName, nUsedLicCount, LmxStat);
}

void HeartbeatCheckoutSuccess_cgo(void *pVendorData, const char *szFeatureName, int nUsedLicCount)
{
        HeartbeatCheckoutSuccessGo(pVendorData, szFeatureName, nUsedLicCount);
}

void HeartbeatRetryFeature_cgo(void *pVendorData, const char *szFeatureName, int nUsedLicCount)
{
        HeartbeatRetryFeatureGo(pVendorData, szFeatureName, nUsedLicCount);
}

void HeartbeatExit_cgo(void *pVendorData)
{
        HeartbeatExitGo(pVendorData);
}
*/
import "C"

type callbacks struct {
	customHostID             func(*HostID, *int) Status
	hostIDCompare            func(HostIDKeyType, HostID, int, HostID, int) Status
	heartbeatConnectionLost  func(interface{}, string, int, int)
	heartbeatCheckoutFailure func(interface{}, string, int, Status)
	heartbeatCheckoutSuccess func(interface{}, string, int)
	heartbeatRetryFeature    func(interface{}, string, int)
	heartbeatExit            func(interface{})
	vendordata               interface{}
}

// TODO: Put the struct elsewhere like the client (but the client is not accessible from Go
// gateways)
var cb callbacks
