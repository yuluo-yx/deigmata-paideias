package indi.yuluo.gptapiserver.util;

import cn.hutool.json.JSONUtil;
import indi.yuluo.gptapiserver.entity.ServiceConfig;
import indi.yuluo.gptapiserver.entity.ServiceEntity;
import jakarta.annotation.Resource;

import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;

/**
 * @author: yuluo
 * @date: 2023/7/21 16:10
 * @description: some desc
 */

public class RequestParameterUtil {

	private RequestParameterUtil(){
	}

	public static HttpEntity<String> getRequestHeaders(ServiceEntity serviceEntity, String apiKey) {

		HttpHeaders httpHeaders = new HttpHeaders();
		httpHeaders.add("Content-Type", "application/json");
		httpHeaders.add("Authorization", "Bearer  " + apiKey);

		return new HttpEntity<>(JSONUtil.toJsonStr(serviceEntity), httpHeaders);
	}


}
