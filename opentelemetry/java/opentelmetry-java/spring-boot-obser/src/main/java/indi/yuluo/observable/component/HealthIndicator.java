package indi.yuluo.observable.component;

import jakarta.annotation.Resource;

import org.springframework.boot.actuate.health.AbstractHealthIndicator;
import org.springframework.boot.actuate.health.Health;
import org.springframework.stereotype.Component;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Component
public class HealthIndicator extends AbstractHealthIndicator {

	@Resource
	private HealthComponent component;

	@Override
	protected void doHealthCheck(Health.Builder builder) {

		int check = component.check();

		if(check == 1) {
			builder.up()
					.withDetail("code","1000")
					.withDetail("msg","Live")
					.withDetail("data","Spring Boot Live")
					.build();
		} else {
			builder.down()
					.withDetail("code","1001")
					.withDetail("msg","Down")
					.withDetail("data","Spring Boot Down")
					.build();
		}
	}
}
