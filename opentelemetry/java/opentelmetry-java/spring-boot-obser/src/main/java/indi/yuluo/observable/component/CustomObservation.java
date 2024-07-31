package indi.yuluo.observable.component;

import java.util.concurrent.atomic.AtomicInteger;

import io.micrometer.observation.Observation;
import io.micrometer.observation.ObservationRegistry;

import org.springframework.stereotype.Component;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Component
public class CustomObservation {

	private final ObservationRegistry observationRegistry;

	CustomObservation(ObservationRegistry observationRegistry) {

		this.observationRegistry = observationRegistry;
	}

	public String productTrace() {

		AtomicInteger count = new AtomicInteger(0);

		Observation observation = Observation.createNotStarted("spring-boot-tracing-operation", this.observationRegistry);
		observation.lowCardinalityKeyValue("spring-boot-tracing-tag", "spring-boot-tracing-value");
		observation.observe(() -> {
			count.getAndIncrement();
			System.out.println("productTrace method exec...");
		});

		return "product tracing data, count" + count;
	}

}
