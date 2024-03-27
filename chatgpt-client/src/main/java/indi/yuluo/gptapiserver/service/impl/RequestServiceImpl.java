package indi.yuluo.gptapiserver.service.impl;

import java.net.InetSocketAddress;
import java.net.Proxy;

import indi.yuluo.gptapiserver.constant.ServiceConstants;
import indi.yuluo.gptapiserver.entity.ServiceConfig;
import indi.yuluo.gptapiserver.entity.ServiceEntity;
import indi.yuluo.gptapiserver.service.RequestService;
import indi.yuluo.gptapiserver.util.RequestParameterUtil;
import indi.yuluo.gptapiserver.util.ResponseParseUtil;
import jakarta.annotation.Resource;

import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

/**
 * @author: yuluo
 * @date: 2023/7/21 16:11
 * @description: some desc
 */

@Service
public class RequestServiceImpl implements RequestService {

	@Resource
	private RestTemplate restTemplate;

	@Resource
	private ServiceConfig serviceConfig;

	@Override
	public String request(ServiceEntity serviceEntity) {

		SimpleClientHttpRequestFactory requestFactory = new SimpleClientHttpRequestFactory();
		requestFactory.setProxy(
				new Proxy(
						Proxy.Type.HTTP,
						new InetSocketAddress(
								serviceConfig.getProxyAddress(),
								1080
						)
				)
		);
		restTemplate.setRequestFactory(requestFactory);

		String responseBody = restTemplate.postForObject(
				ServiceConstants.API_URL,
				RequestParameterUtil.getRequestHeaders(serviceEntity, serviceConfig.getApiKey()),
				String.class
		);

		return ResponseParseUtil.parseResponseBody(responseBody);
	}

}
