package indi.yuluo.langchain4j.service.chathistory;

import java.util.Map;

import dev.langchain4j.data.message.UserMessage;
import dev.langchain4j.memory.ChatMemory;
import dev.langchain4j.memory.chat.MessageWindowChatMemory;
import dev.langchain4j.model.chat.ChatLanguageModel;
import jakarta.annotation.Resource;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class ChatHistoryService {

	@Resource
	private ChatLanguageModel chatLanguageModel;

	public Map<String, String> chat(String message) {

		// ChatMemory chatMemory = TokenWindowChatMemory.withMaxTokens(1000, tokenizer);
		ChatMemory chatMemory = MessageWindowChatMemory.withMaxMessages(10);

		String resp = chatLanguageModel.generate(message);

		// 历史记录保存了用户和 AI 之间的所有消息。历史记录是用户在 UI 中看到的内容。它代表了实际说过的内容。
		// 记忆会保留一些信息，这些信息会呈现给 LLM，使其表现得好像“记住”了对话一样。
		// 记忆与历史截然不同。根据所使用的记忆算法，它可以以各种方式修改历史：删除一些消息、汇总多条消息、汇总单独的消息、
		// 从消息中删除不重要的细节、在消息中注入额外信息（例如，用于 RAG）或指令（例如，用于结构化输出）等等。

		// LangChain4j 目前只提供“记忆”，不提供“历史”，如果需要保留整个历史，请手动操作
		// https://github.com/langchain4j/langchain4j-examples/blob/a37df7be40080a93a0e38231b93e1ed24ae0085b/tutorials/src/main/java/_05_Memory.java
		chatMemory.add(new UserMessage(message));

		return Map.of(message, resp);
	}

}
