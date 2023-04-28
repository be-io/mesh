package com.be.mesh.client.metric;

import com.be.mesh.client.annotate.Listener;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.Publisher;
import com.be.mesh.client.prsim.RuntimeHook;
import com.be.mesh.client.prsim.Scheduler;
import com.be.mesh.client.prsim.Subscriber;
import com.be.mesh.client.struct.Event;
import com.be.mesh.client.struct.Topic;
import com.be.mesh.client.tool.Mode;
import com.be.mesh.client.tool.Tool;
import io.prometheus.client.Collector;
import io.prometheus.client.hotspot.DefaultExports;
import lombok.extern.slf4j.Slf4j;
import oshi.SystemInfo;

import java.time.Duration;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.List;

/**
 * exporter data via registry scheduler, single thread
 * <p>
 * 同时，需要添加以下3个依赖：
 * com.github.oshi:oshi-core，建议版本：6.1.5
 * net.java.dev.jna:jna-platform，建议版本：5.11.0
 * net.java.dev.jna:jna，建议版本：5.11.0
 * <p>
 * * @author jianyue.li
 */
@Slf4j
@Listener(topic = "mesh.prom.sample.period")
public class MetricsExporter implements RuntimeHook, Subscriber {

    private final Topic periodBinding = new Topic("mesh.prom.sample.period", "*");
    private final Topic sampleBinding = new Topic("mesh.prom.sample.export", "*");
    private final Publisher publisher = ServiceLoader.load(Publisher.class).getDefault();
    private final Scheduler scheduler = ServiceLoader.load(Scheduler.class).getDefault();
    public static final SystemInfo SYSTEM = new SystemInfo();

    @Override
    public void start() throws Throwable {
        refresh();
    }

    @Override
    public void refresh() throws Throwable {
        if (Tool.isClassPresent("io.prometheus.client.Collector") && Tool.MESH_MODE.get().match(Mode.Metrics)) {
            scheduler.period(Duration.ofSeconds(10), this.periodBinding);

            DefaultExports.register(Registry.REGISTRY.get());
            new MeminfoCollector().register(Registry.REGISTRY.get());
            new DiskStatsCollector().register(Registry.REGISTRY.get());
            new NetDevCollector().register(Registry.REGISTRY.get());
            new CPUCollector().register(Registry.REGISTRY.get());
        }
    }

    @Override
    public void stop() throws Throwable {
        //
    }

    @Override
    public void subscribe(Event event) {
        Enumeration<Collector.MetricFamilySamples> enumeration = Registry.REGISTRY.get().metricFamilySamples();
        List<Collector.MetricFamilySamples> metricFamilySamples = new ArrayList<>();
        if (enumeration != null) {
            while (enumeration.hasMoreElements()) {
                Collector.MetricFamilySamples samples = enumeration.nextElement();
                metricFamilySamples.add(samples);
            }
        }
        try {
            ExporterSamples es = new ExporterSamples();
            es.setInstanceId(Tool.HOST_NAME.get());
            es.setMetricFamilySamples(metricFamilySamples);
            publisher.publish(this.sampleBinding, es);
        } catch (Exception e) {
            log.warn("Publish metrics event with error. ", e);
        }
    }

    public static class ExporterSamples {

        private List<Collector.MetricFamilySamples> metricFamilySamples;

        /**
         * such as:theta-xx-yy
         */
        private String instanceId;

        private String jobName = "pushgateway";

        public List<Collector.MetricFamilySamples> getMetricFamilySamples() {
            return metricFamilySamples;
        }

        public void setMetricFamilySamples(List<Collector.MetricFamilySamples> metricFamilySamples) {
            this.metricFamilySamples = metricFamilySamples;
        }

        public String getInstanceId() {
            return instanceId;
        }

        public void setInstanceId(String instanceId) {
            this.instanceId = instanceId;
        }

        public String getJobName() {
            return jobName;
        }

        public void setJobName(String jobName) {
            this.jobName = jobName;
        }
    }

}
