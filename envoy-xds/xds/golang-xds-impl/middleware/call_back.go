package middleware

import (
	"context"
	"golang-xds-impl/pkg/log"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
)

// callback 回调日志信息

type CallBacksMiddleWares struct {
	XLog log.XLogger
}

func (cbm CallBacksMiddleWares) OnFetchRequest(ctx context.Context, request *discovery.DiscoveryRequest) error {

	if cbm.XLog.Debug {
		cbm.XLog.Infof("OnFetchRequest exec ....")
	}

	return nil
}

func (cbm CallBacksMiddleWares) OnFetchResponse(request *discovery.DiscoveryRequest, response *discovery.DiscoveryResponse) {
	if cbm.XLog.Debug {
		cbm.XLog.Infof("OnFetchResponse exec ....")
	}
}

func (cbm CallBacksMiddleWares) OnStreamOpen(ctx context.Context, id int64, types string) error {
	if cbm.XLog.Debug {
		cbm.XLog.Infof("OnStreamOpen exec .... id is %d, type is %v", id, types)
	}

	return nil
}

func (cbm CallBacksMiddleWares) OnStreamClosed(id int64, node *corev3.Node) {

	if cbm.XLog.Debug {
		cbm.XLog.Infof("OnStreamClosed exec .... stream id is %d, node is %s", id, node.Id)
	}
}

func (cbm CallBacksMiddleWares) OnStreamRequest(id int64, request *discovery.DiscoveryRequest) error {

	if cbm.XLog.Debug {
		cbm.XLog.Infof("OnStreamRequest exec ....")
	}

	return nil
}

func (cbm CallBacksMiddleWares) OnStreamResponse(ctx context.Context, id int64, request *discovery.DiscoveryRequest, response *discovery.DiscoveryResponse) {

	if cbm.XLog.Debug {

		cbm.XLog.Infof("OnStreamResponse exec .... id is %d.", id)
	}
}

func (cbm CallBacksMiddleWares) OnDeltaStreamOpen(ctx context.Context, id int64, types string) error {

	if cbm.XLog.Debug {

		cbm.XLog.Infof("OnDeltaStreamOpen exec .... id is %d, type is %v", id, types)
	}

	return nil
}

func (cbm CallBacksMiddleWares) OnDeltaStreamClosed(id int64, node *corev3.Node) {

	if cbm.XLog.Debug {

		cbm.XLog.Infof("OnDeltaStreamClosed exec ....")
	}
}

func (cbm CallBacksMiddleWares) OnStreamDeltaRequest(id int64, request *discovery.DeltaDiscoveryRequest) error {

	if cbm.XLog.Debug {

		cbm.XLog.Infof("OnStreamDeltaRequest exec ....")
	}

	return nil
}

func (cbm CallBacksMiddleWares) OnStreamDeltaResponse(id int64, request *discovery.DeltaDiscoveryRequest, response *discovery.DeltaDiscoveryResponse) {

	if cbm.XLog.Debug {

		cbm.XLog.Infof("OnStreamDeltaResponse exec ....")
	}
}
