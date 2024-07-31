package indi.yuluo.otel.custom.service;

import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.StatusCode;
import io.opentelemetry.api.trace.Tracer;
import jakarta.annotation.Resource;

import org.springframework.stereotype.Service;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Service
public class OtelService {

	@Resource
	private Tracer tracer;

	public String hello() {

		String info;

		Span span = tracer.spanBuilder("GET /otel").startSpan();

		try {
			info = "Hello otel!";

			span.setAttribute("/otel -> info", "Hello otel!");
		} catch (Exception e) {
			span.recordException(e)
					.setStatus(StatusCode.ERROR, e.getMessage()
					);
			return "Error";
		} finally {
			span.end();
		}

		return info;
	}

}
