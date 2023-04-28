package com.be.mesh.client.metric;

import com.be.mesh.client.tool.Tool;
import io.prometheus.client.Collector;
import lombok.extern.slf4j.Slf4j;
import oshi.hardware.GlobalMemory;

import java.util.ArrayList;
import java.util.List;

/**
 * 基于oshi
 * linux下参考：/proc/meminfo
 * MemTotal:       15731896 kB
 * MemFree:          222268 kB
 * MemAvailable:    4688532 kB
 * Buffers:           28424 kB
 * Cached:          4567872 kB
 * SwapCached:            0 kB
 * Active:         12123932 kB
 * Inactive:        2629856 kB
 * Active(anon):   10161464 kB
 * Inactive(anon):     2292 kB
 * Active(file):    1962468 kB
 * Inactive(file):  2627564 kB
 * Unevictable:           0 kB
 * Mlocked:               0 kB
 * SwapTotal:             0 kB
 * SwapFree:              0 kB
 * Dirty:             13220 kB
 * Writeback:             0 kB
 * AnonPages:      10157528 kB
 *
 * @author jianyue.li
 * @date 2022/4/7 2:18 PM
 */
@Slf4j
public class MeminfoCollector extends Collector {

    private static final List<String> constLabels = new ArrayList<>();

    static {
        constLabels.add("host");
    }

    private static final String help = "Memory information field %s.";

    @Override
    public List<MetricFamilySamples> collect() {
        List<MetricFamilySamples> mfs = new ArrayList<>();

        GlobalMemory mem = MetricsExporter.SYSTEM.getHardware().getMemory();
        if (mem != null) {
            long available = mem.getAvailable();
            long total = mem.getTotal();
            long active = total - available;

            mfs.add(assemble("available", String.format(help, "available"), available, Type.GAUGE));
            mfs.add(assemble("total", String.format(help, "total"), total, Type.COUNTER));
            mfs.add(assemble("active", String.format(help, "active"), active, Type.GAUGE));
        }
        return mfs;
    }

    private MetricFamilySamples assemble(String metricName, String help, double val, Type type) {
        List<MetricFamilySamples.Sample> samples = new ArrayList<>();
        List<String> labelVals = new ArrayList<>();
        labelVals.add(Tool.HOST_NAME.get());
        String mn = "node_memory_" + metricName + "_bytes";
        samples.add(new MetricFamilySamples.Sample(mn, constLabels, labelVals, val));
        return new MetricFamilySamples(mn, type, help, samples);
    }

}
