package indi.yuluo.observable.component;

import java.util.Random;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;

import org.springframework.stereotype.Component;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Component
public class HealthComponent {

	private static Counter counter;

	public HealthComponent(MeterRegistry registry) {

		counter = registry.counter("spring-boot.count");
	}

	public Integer check() {

		return new Random().nextInt(2);
	}

	public String count() {

		counter.increment();

		return "Call count, count ++";
	}

}
