package indi.yuluo.otel.custom.config;

import indi.yuluo.otel.custom.OtelCustomApplication;
import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.baggage.propagation.W3CBaggagePropagator;
import io.opentelemetry.api.common.Attributes;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.api.trace.propagation.W3CTraceContextPropagator;
import io.opentelemetry.context.propagation.ContextPropagators;
import io.opentelemetry.context.propagation.TextMapPropagator;
import io.opentelemetry.exporter.otlp.trace.OtlpGrpcSpanExporter;
import io.opentelemetry.sdk.OpenTelemetrySdk;
import io.opentelemetry.sdk.resources.Resource;
import io.opentelemetry.sdk.trace.SdkTracerProvider;
import io.opentelemetry.sdk.trace.export.SimpleSpanProcessor;
import io.opentelemetry.semconv.resource.attributes.ResourceAttributes;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * @author yuluo
 * @author <a href="mailto:yuluo08290126@gmail.com">yuluo</a>
 */

@Configuration
public class OtelCustomConfiguration {

	@Bean
	public OpenTelemetry openTelemetry() {

		GlobalOpenTelemetry.resetForTest();

		String exporterEndpointFromEnv = System.getenv("OTLP_EXPORTER_ENDPOINT");
		String exporterEndpoint =
				exporterEndpointFromEnv != null ?
						exporterEndpointFromEnv :
						"http://localhost:4317";

		String serviceNameFromEnv = System.getenv("SERVICE_NAME");
		String serviceName =
				serviceNameFromEnv != null ?
						serviceNameFromEnv :
						"otel-example";

		OtlpGrpcSpanExporter exporter = OtlpGrpcSpanExporter.builder()
				.setEndpoint(exporterEndpoint)
				.build();



		Resource resource = Resource.getDefault()
				.merge(Resource.create(Attributes.of(
						ResourceAttributes.SERVICE_NAME, serviceName,
						ResourceAttributes.TELEMETRY_SDK_LANGUAGE, "java")));

		SdkTracerProvider tracerProvider = SdkTracerProvider.builder()
				.addSpanProcessor(SimpleSpanProcessor.create(exporter))
				.setResource(resource)
				.build();

		ContextPropagators propagators = ContextPropagators.create(
				TextMapPropagator.composite(
						W3CTraceContextPropagator.getInstance(),
						W3CBaggagePropagator.getInstance()));

		return OpenTelemetrySdk.builder()
				.setPropagators(propagators)
				.setTracerProvider(tracerProvider)
				.buildAndRegisterGlobal();
	}

	@Bean
	public Tracer tracer() {

		return openTelemetry().getTracer(OtelCustomApplication.class.getName());
	}

}
