package indi.yuluo.gptapiserver.entity;

import java.util.HashMap;
import java.util.List;

/**
 * @author: yuluo
 * @date: 2023/7/21 16:04
 * @description: some desc
 */

public class ServiceEntity {

	private List<HashMap<String, String>> messages;

	private String model;

	private float temperature;

	public List<HashMap<String, String>> getMessages() {
		return messages;
	}

	public void setMessages(List<HashMap<String, String>> messages) {
		this.messages = messages;
	}

	public String getModel() {
		return model;
	}

	public void setModel(String model) {
		this.model = model;
	}

	public float getTemperature() {
		return temperature;
	}

	public void setTemperature(float temperature) {
		this.temperature = temperature;
	}
}
