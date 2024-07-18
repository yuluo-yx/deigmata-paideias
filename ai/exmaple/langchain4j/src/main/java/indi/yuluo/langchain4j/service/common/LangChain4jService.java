package indi.yuluo.langchain4j.service.common;

import java.util.Map;

import dev.langchain4j.model.chat.ChatLanguageModel;
import jakarta.annotation.Resource;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class LangChain4jService {

	@Resource
	private ChatLanguageModel chatLanguageModel;

	public Map<String, String> chat(String message) {

		return Map.of(message, chatLanguageModel.generate(message));
	}

}
