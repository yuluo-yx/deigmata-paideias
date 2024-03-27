package indi.yuluo.envoy.xds.core;

import io.envoyproxy.controlplane.server.DiscoveryServerCallbacks;
import io.envoyproxy.controlplane.server.exception.RequestException;
import io.envoyproxy.envoy.service.discovery.v3.DeltaDiscoveryRequest;
import io.envoyproxy.envoy.service.discovery.v3.DeltaDiscoveryResponse;
import io.envoyproxy.envoy.service.discovery.v3.DiscoveryRequest;
import io.envoyproxy.envoy.service.discovery.v3.DiscoveryResponse;

/**
 * @author yuluo
 * @author 1481556636@qq.com
 */

public class CallBack implements DiscoveryServerCallbacks {

	@Override
	public void onStreamClose(long streamId, String typeUrl) {

		System.out.println("onStreamClose exec ... id " + streamId + " typeUrl " + typeUrl);
	}

	@Override
	public void onStreamCloseWithError(long streamId, String typeUrl, Throwable error) {

		System.out.println("onStreamCloseWithError exec ... id " + streamId + " typeUrl " + typeUrl);
	}

	@Override
	public void onStreamOpen(long streamId, String typeUrl) throws RequestException {

		System.out.println("onStreamOpen exec ... id " + streamId + " typeUrl " + typeUrl);
	}

	@Override
	public void onV3StreamRequest(long l, DiscoveryRequest discoveryRequest) throws RequestException {

		System.out.println("onV3StreamRequest exec ... id " + l);
	}

	@Override
	public void onV3StreamDeltaRequest(long l, DeltaDiscoveryRequest deltaDiscoveryRequest) throws RequestException {

		System.out.println("onV3StreamDeltaRequest exec ... id " + l);
	}

	@Override
	public void onV3StreamResponse(long streamId, DiscoveryRequest request, DiscoveryResponse response) {

		System.out.println("onV3StreamResponse exec ... id " + streamId);
	}

	@Override
	public void onV3StreamDeltaResponse(long streamId, DeltaDiscoveryRequest request, DeltaDiscoveryResponse response) {

		System.out.println("onV3StreamDeltaResponse exec ... id " + streamId);
	}
}
