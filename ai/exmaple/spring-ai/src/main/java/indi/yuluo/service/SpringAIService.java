package indi.yuluo.service;

import java.util.Map;

import jakarta.annotation.Resource;

import org.springframework.ai.chat.model.ChatResponse;
import org.springframework.ai.chat.prompt.Prompt;
import org.springframework.ai.ollama.OllamaChatModel;
import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class SpringAIService {

	@Resource
	private OllamaChatModel ollamaChatModel;

	public Map<String, ChatResponse> chat(String msg) {

		return Map.of("generation", ollamaChatModel.call(new Prompt(msg)));
	}

}
