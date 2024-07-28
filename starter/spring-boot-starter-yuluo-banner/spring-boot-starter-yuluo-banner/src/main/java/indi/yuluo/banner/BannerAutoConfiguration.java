package indi.yuluo.banner;

import org.springframework.boot.Banner;
import org.springframework.boot.ResourceBanner;
import org.springframework.boot.autoconfigure.AutoConfiguration;
import org.springframework.boot.autoconfigure.condition.ConditionalOnMissingBean;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Scope;
import org.springframework.core.io.ClassPathResource;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@AutoConfiguration
public class BannerAutoConfiguration {

	@Bean
	@Scope("prototype")
	@ConditionalOnMissingBean
	public Banner banner() {

		return new ResourceBanner(new ClassPathResource("banner.txt"));
	}

}
