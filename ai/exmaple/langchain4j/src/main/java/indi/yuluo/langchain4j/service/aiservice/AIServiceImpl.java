package indi.yuluo.langchain4j.service.aiservice;

import dev.langchain4j.service.AiServices;
import dev.langchain4j.model.chat.ChatLanguageModel;
import jakarta.annotation.Resource;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class AIServiceImpl {

	@Resource
	private ChatLanguageModel chatLanguageModel;

	public void chat(String userMsg) {

		Assistant assistant = AiServices.create(Assistant.class, chatLanguageModel);
		System.out.println("正在使用 LangChain4j 的 AiServices 功能，您将得到以下输出：\n ========================================");
		System.out.println(assistant.chat(userMsg));
	}

}
