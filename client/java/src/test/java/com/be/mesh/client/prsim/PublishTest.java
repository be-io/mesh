/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.metric.MetricsExporter;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.spring.MeshConfiguration;
import com.be.mesh.client.struct.Event;
import com.be.mesh.client.struct.Topic;
import io.prometheus.client.Collector;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.springframework.test.context.ContextConfiguration;
import org.testng.annotations.Test;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@ContextConfiguration(classes = MeshConfiguration.class)
public class PublishTest {

    @Test
    public void loadTest() {
        Publisher publisher = ServiceLoader.load(Publisher.class).getDefault();
        log.info("{}", publisher);
    }

    @Data
    private static final class ProcessTaskCreateEvent {
        private String businessNo;
        // 申请标的物编号
        private String objectNo;
        // 申请标的物类型
        private String objectType;
        // 申请标的物名称
        private String objectName;
        // 申请标的物详情地址
        private String objectURI;
        // 申请目标机构号
        private String targetInstId;
        // 申请机构号
        private String applyInstId;
        // 申请操作员号
        private String applyOperatorId;
        // 申请来源
        private String applyFrom;
        // 申请类型
        private String applyType;
        // 申请说明
        private String applyDesc;
        // 申请时间
        private LocalDateTime applyAt;
        // 审批上下文
        private String context;
    }

    @Data
    private static final class ProcessTaskOperateEvent {
        // 审批机构号
        private String operateInstId;
        // 审批操作员号
        private String operateOperatorId;
        // 审批说明
        private String operateDesc;
        // 审批时间
        private LocalDateTime operateAt;
        // 操作类型 APPROVED 审批通过 REJECT 审批驳回 REVOKE 审批撤销
        private String operateType;
        // 申请标的物编号
        private String objectNo;
        // 申请标的物类型
        private String objectType;
        // 申请标的物名称
        private String objectName;
        // 申请标的物详情地址
        private String objectURI;
        // 申请目标机构号
        private String targetInstId;
        // 申请机构号
        private String applyInstId;
        // 申请操作员号
        private String applyOperatorId;
        // 申请来源
        private String applyFrom;
        // 申请类型
        private String applyType;
        // 申请说明
        private String applyDesc;
        // 申请时间
        private LocalDateTime applyAt;
    }

    @Test
    public void publishTest() {
        ProcessTaskCreateEvent event = new ProcessTaskCreateEvent();
        event.setBusinessNo("");
        event.setObjectNo("");
        event.setObjectType("");
        event.setObjectName("");
        event.setObjectURI("");
        event.setTargetInstId("");
        event.setApplyInstId("");
        event.setApplyOperatorId("");
        event.setApplyFrom("");
        event.setApplyType("");
        event.setApplyDesc("");
        event.setApplyAt(LocalDateTime.now());
        Publisher publisher = ServiceLoader.load(Publisher.class).getDefault();
        String eventId = publisher.publish(new Topic("gaia.pandora.process.create", "*"), event);
        log.info(eventId);
    }

    @Test
    public void publishPromSamplesTest() throws Throwable {
        MetricsExporter.ExporterSamples es = new MetricsExporter.ExporterSamples();
        es.setJobName("pushgateway");
        List<Collector.MetricFamilySamples> samples = new ArrayList<>();
        String help = "";
        List<Collector.MetricFamilySamples.Sample> childSamples = new ArrayList<>();

        List<String> labelNames = new ArrayList<>();
        labelNames.add("gc");
        List<String> labelValues = new ArrayList<>();
        labelValues.add("G1 Young Generation");
        Collector.MetricFamilySamples.Sample sml = new Collector.MetricFamilySamples.Sample("jvm_gc_collection_seconds_count", labelNames, labelValues, 3.0);

        childSamples.add(sml);

        Collector.MetricFamilySamples metricFamilySamples = new Collector.MetricFamilySamples("jvm_gc_collection_seconds", "", Collector.Type.COUNTER, help, childSamples);

        es.setMetricFamilySamples(samples);

        Publisher publisher = ServiceLoader.load(Publisher.class).getDefault();
        Mesh.contextSafeCaught(() -> {
            String eventId = publisher.publish(new Topic("prometheus.exporter.sample", "prometheus.exporter.sample"), metricFamilySamples);
            log.info(eventId);
        });
    }

    @Test
    public void broadcastTest() {
        Publisher publisher = ServiceLoader.load(Publisher.class).getDefault();
        List<String> result = publisher.broadcast(new Topic(), "ping");
        result.forEach(log::info);
    }

    @Test
    public void brokerExchangeTest() {
        String event = "{\"version\":\"1.0.0\",\"tid\":\"a5aa76c1499354608106745858\",\"eid\":\"1499354608106745859\",\"mid\":\"1499354608106745860\",\"timestamp\":\"1646308969891\",\"source\":{\"node_id\":\"LX0000010000040\",\"inst_id\":\"JG0100000400000000\"},\"target\":{\"node_id\":\"JG0100000100000000\",\"inst_id\":\"JG0100000100000000\"},\"binding\":{\"topic\":\"alliance.ctrl.pboc\",\"code\":\"alliance.ctrl.pboc\"},\"entity\":{\"codec\":\"json\",\"schema\":\"\",\"buffer\":\"{\\\"requestId\\\":\\\"4185629816024064\\\",\\\"modelName\\\":\\\"cctest三方纵向逻辑回归_A模型+B数据+C数据\\\",\\\"partnerCode\\\":\\\"JG0100000400000000\\\",\\\"partnerName\\\":\\\"194\\\",\\\"status\\\":1}\"}}";
        Publisher publisher = ServiceLoader.load(Publisher.class).getDefault();
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        List<String> eventIds = publisher.publish(Collections.singletonList(codec.decodeString(event, Types.of(Event.class))));
        log.info("{}", eventIds);
    }
}
