package com.be.mesh.client.metric;

import com.be.mesh.client.tool.Tool;
import io.prometheus.client.Collector;
import io.prometheus.client.CounterMetricFamily;
import lombok.extern.slf4j.Slf4j;
import oshi.SystemInfo;
import oshi.hardware.CentralProcessor;

import java.util.ArrayList;
import java.util.List;

/**
 * 基于oshi
 * linux下格式参考：/proc/stat
 * cpu  2441020 72974 809111 57135110 161589 0 33882 0 0 0
 * cpu0 356252 1151 89941 7043801 40536 0 8403 0 0 0
 * cpu1 276674 16984 115070 7193746 16031 0 2344 0 0 0
 * cpu2 333010 1193 87175 7090573 25337 0 5301 0 0 0
 * cpu3 277321 17281 115682 7195354 13131 0 2357 0 0 0
 * cpu4 346507 1186 85144 7084595 21593 0 4745 0 0 0
 * cpu5 256937 16913 113869 7218945 12192 0 2298 0 0 0
 * cpu6 345075 1168 86857 7083718 19938 0 6093 0 0 0
 * cpu7 249239 17094 115372 7224375 12827 0 2339 0 0 0
 * .....
 *
 * @author jianyue.li
 * @date 2022/4/8 1:31 PM
 */
@Slf4j
public class CPUCollector extends Collector {

    private static final List<String> constLabels = new ArrayList<>();

    static {
        constLabels.add("cpu");
        constLabels.add("mode");
        constLabels.add("host");
    }

    @Override
    public List<MetricFamilySamples> collect() {
        List<MetricFamilySamples> mfs = new ArrayList<>();

        SystemInfo sysInfo = MetricsExporter.SYSTEM;
        CentralProcessor cp = sysInfo.getHardware().getProcessor();
        if (cp != null) {
            long[][] cpuUsedTicks = cp.getProcessorCpuLoadTicks();
            if (cpuUsedTicks != null) {
                CounterMetricFamily cf = new CounterMetricFamily("node_cpu_seconds_total", "Seconds the CPUs spent in each mode.", constLabels);
                for (int i = 0; i < cpuUsedTicks.length; i++) {
                    String cpuId = String.valueOf(i);
                    long[] usedTicks = cpuUsedTicks[i];
                    if (usedTicks == null || usedTicks.length < 7) {
                        continue;
                    }

                    cf.addMetric(labelValues(cpuId, "user"), usedTicks[0]);
                    cf.addMetric(labelValues(cpuId, "nice"), usedTicks[1]);
                    cf.addMetric(labelValues(cpuId, "system"), usedTicks[2]);
                    cf.addMetric(labelValues(cpuId, "idle"), usedTicks[3]);
                    cf.addMetric(labelValues(cpuId, "iowait"), usedTicks[4]);
                    cf.addMetric(labelValues(cpuId, "irq"), usedTicks[5]);
                    cf.addMetric(labelValues(cpuId, "softirq"), usedTicks[6]);

                    if (usedTicks.length >= 8) {
                        cf.addMetric(labelValues(cpuId, "steal"), usedTicks[7]);
                    }
                    if (usedTicks.length >= 9) {
                        cf.addMetric(labelValues(cpuId, "states"), usedTicks[8]);
                    }
                }
                mfs.add(cf);
            }
        }
        return mfs;
    }

    /**
     * @param mode user nice...
     * @return
     */
    private List<String> labelValues(String cpuId, String mode) {
        List<String> labelVals = new ArrayList<>();
        labelVals.add(cpuId);
        labelVals.add(mode);
        labelVals.add(Tool.HOST_NAME.get());
        return labelVals;
    }

}
