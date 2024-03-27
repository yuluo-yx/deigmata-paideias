package indi.yuluo.gptapiserver.util;

import cn.hutool.json.JSONObject;

/**
 * @author: yuluo
 * @date: 2023/7/21 16:11
 * @description: some desc
 */

public class ResponseParseUtil {

	private ResponseParseUtil(){
	}

	public static String parseResponseBody(String body) {

		String choices = new JSONObject(body).get("choices").toString();
		choices = choices.substring(1);
		choices = choices.substring(0, choices.length() - 1);
		String message = new JSONObject(choices).get("message").toString();

		return new JSONObject(message).get("content").toString();
	}

}
