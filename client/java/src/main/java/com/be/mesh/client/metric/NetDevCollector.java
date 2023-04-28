package com.be.mesh.client.metric;

import com.be.mesh.client.tool.Tool;
import io.prometheus.client.Collector;
import io.prometheus.client.CounterMetricFamily;
import lombok.extern.slf4j.Slf4j;
import oshi.hardware.NetworkIF;

import java.util.ArrayList;
import java.util.List;

/**
 * 基于oshi
 * linux下参考：/proc/net/dev
 * <p>
 * 格式参考：
 * Inter-|   Receive                                                |  Transmit
 * face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
 * eth0: 17568967  258205    0    0    0     0          0         0 10968283  156370    0    0    0     0       0          0
 * lo: 6001019   72596    0    0    0     0          0         0  6001019   72596    0    0    0     0       0          0
 * tunl0:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
 *
 * @author jianyue.li
 * @date 2022/4/7 8:20 PM
 */
@Slf4j
public class NetDevCollector extends Collector {

    private static final List<String> constLabels = new ArrayList<>();

    static {
        constLabels.add("device");
        constLabels.add("host");
    }

    @Override
    public List<MetricFamilySamples> collect() {
        List<MetricFamilySamples> mfs = new ArrayList<>();

        List<NetworkIF> networkIFs = MetricsExporter.SYSTEM.getHardware().getNetworkIFs();
        if (Tool.required(networkIFs)) {
            CounterMetricFamily receiveCmf = new CounterMetricFamily("node_network_receive_total", "Network device statistic receive.", constLabels);
            CounterMetricFamily transmitCmf = new CounterMetricFamily("node_network_transmit_total", "Network device statistic transmit.", constLabels);

            for (NetworkIF network : networkIFs) {
                boolean success = network.updateAttributes();
                if (!success) {
                    continue;
                }
                List<String> labelVals = new ArrayList<>();
                labelVals.add(network.getName());
                labelVals.add(Tool.HOST_NAME.get());
                receiveCmf.addMetric(labelVals, network.getBytesRecv());
                transmitCmf.addMetric(labelVals, network.getBytesSent());
            }
            mfs.add(receiveCmf);
            mfs.add(transmitCmf);
        }
        return mfs;
    }

}
