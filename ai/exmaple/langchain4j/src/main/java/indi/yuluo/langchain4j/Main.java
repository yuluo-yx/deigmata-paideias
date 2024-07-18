package indi.yuluo.langchain4j;

import dev.langchain4j.model.chat.ChatLanguageModel;
import dev.langchain4j.model.ollama.OllamaChatModel;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

public class Main {

	public static void main(String[] args) {
		ChatLanguageModel model = OllamaChatModel.builder()
				.baseUrl("http://localhost:11434")
				.modelName("qwen:4b")
				.build();

		String answer = model.generate("你好");
		System.out.println(answer);
	}

}
