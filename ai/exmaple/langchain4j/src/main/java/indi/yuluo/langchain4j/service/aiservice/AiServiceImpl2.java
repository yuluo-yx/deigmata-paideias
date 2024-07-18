package indi.yuluo.langchain4j.service.aiservice;

import java.util.List;

import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.service.AiServices;
import jakarta.annotation.Resource;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class AiServiceImpl2 {
	@Resource
	private ChatLanguageModel chatLanguageModel;

	public void chat() {

		System.out.println("正在使用 LangChain4j 的 AiServices 功能，系统和用户 Message。您将得到以下输出：\n ========================================");
		SystemAndUserMessage utils = AiServices.create(SystemAndUserMessage.class, chatLanguageModel);

		String translation = utils.translate("Hello, how are you?", "italian");
		System.out.println(translation);

		String text = "AI, or artificial intelligence, is a branch of computer science that aims to create "
				+ "machines that mimic human intelligence. This can range from simple tasks such as recognizing "
				+ "patterns or speech to more complex tasks like making decisions or predictions.";

		List<String> bulletPoints = utils.summarize(text, 3);
		bulletPoints.forEach(System.out::println);
	}

}
