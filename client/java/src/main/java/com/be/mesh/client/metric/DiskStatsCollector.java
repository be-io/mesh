package com.be.mesh.client.metric;

import com.be.mesh.client.tool.Tool;
import io.prometheus.client.Collector;
import io.prometheus.client.CounterMetricFamily;
import io.prometheus.client.GaugeMetricFamily;
import lombok.extern.slf4j.Slf4j;
import oshi.hardware.HWDiskStore;

import java.util.ArrayList;
import java.util.List;

/**
 * 基于oshi
 * linux下参考：/proc/diskstats
 * 代码实现参考：node_exporter: diskstats_common.go,diskstats_linux.go
 * 参考格式：
 * 253       0 vda 47413 1 3219405 278873 8502037 5941642 163207250 13422246 0 705432 13701119
 * 253       1 vda1 47373 1 3217293 278769 8501596 5941642 163207250 13422235 0 705331 13701004
 * 253      16 vdb 977777 124 35194882 2597108 47361419 31726213 1379914744 324357285 0 5478090 326954393
 * 253      17 vdb1 977739 124 35192786 2597006 47093414 31726213 1379914744 324312851 0 5464041 326909857
 *
 * @author jianyue.li
 * @date 2022/4/7 3:28 PM
 */
@Slf4j
public class DiskStatsCollector extends Collector {

    private static final List<String> constLabels = new ArrayList<>();


    static {
        constLabels.add("device");
        constLabels.add("host");
    }

    @Override
    public List<MetricFamilySamples> collect() {
        List<MetricFamilySamples> mfs = new ArrayList<>();
        List<HWDiskStore> stores = MetricsExporter.SYSTEM.getHardware().getDiskStores();
        if (Tool.required(stores)) {
            CounterMetricFamily readCompleted = new CounterMetricFamily("node_disk_reads_completed_total", "The total number of reads completed successfully.", constLabels);
            CounterMetricFamily readBytes = new CounterMetricFamily("node_disk_read_bytes_total", "The total number of bytes read successfully.", constLabels);
            CounterMetricFamily writeCompleted = new CounterMetricFamily("node_disk_writes_completed_total", "The total number of writes completed successfully.", constLabels);
            CounterMetricFamily writeBytes = new CounterMetricFamily("node_disk_written_bytes_total", "The total number of bytes written successfully.", constLabels);
            CounterMetricFamily uiTimeSeconds = new CounterMetricFamily("node_disk_io_time_seconds_total", "Total seconds spent doing I/Os.", constLabels);
            GaugeMetricFamily uiNow = new GaugeMetricFamily("node_disk_io_now", "The total number of I/Os currently in progress.", constLabels);
            for (HWDiskStore disk : stores) {
                String devName = disk.getName();
                // add other
                List<String> labelValues = new ArrayList<>();
                labelValues.add(devName);
                labelValues.add(Tool.HOST_NAME.get());

                readCompleted.addMetric(labelValues, disk.getReads());
                readBytes.addMetric(labelValues, disk.getReadBytes());
                writeCompleted.addMetric(labelValues, disk.getWrites());
                writeBytes.addMetric(labelValues, disk.getWriteBytes());
                uiTimeSeconds.addMetric(labelValues, disk.getTransferTime() / 1000);
                uiNow.addMetric(labelValues, disk.getCurrentQueueLength());
            }

            mfs.add(readCompleted);
            mfs.add(readBytes);
            mfs.add(writeCompleted);
            mfs.add(writeBytes);
            mfs.add(uiTimeSeconds);
            mfs.add(uiNow);
        }
        return mfs;
    }

}
