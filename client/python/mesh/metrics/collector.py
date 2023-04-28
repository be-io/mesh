#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

# import psutil
# import socket
# from prometheus_client.metrics_core import CounterMetricFamily, GaugeMetricFamily
#
# from prometheus_client.registry import Collector
#
# from typing import Iterable
#
# from prometheus_client import Metric
#
#
# class CPUCollector(Collector):
#
#     def __init__(self):
#         self.hostname = socket.gethostname()
#
#     def collect(self) -> Iterable[Metric]:
#         total_cpu_times = psutil.cpu_times(True)
#         cpu_counter = CounterMetricFamily("node_cpu_seconds_total", "Seconds the CPUs spent in each mode.",
#                                           labels=["cpu", "mode", "host"])
#         for idx, per_cpu in enumerate(total_cpu_times):
#             print(per_cpu)
#             for k, v in per_cpu._asdict().items():
#                 cpu_counter.add_metric([idx, k, self.hostname], value=v)
#         return [cpu_counter]
#
#
# class DiskStatsCollector(Collector):
#
#     def __init__(self):
#         self.hostname = socket.gethostname()
#
#     def collect(self) -> Iterable[Metric]:
#         read_completed_counter = CounterMetricFamily("node_memory_reads_completed_total",
#                                                      "The total number of reads completed successfully.",
#                                                      labels=["device", "host"])
#         read_bytes_total = CounterMetricFamily("node_memory_read_bytes_total",
#                                                "The total number of bytes read successfully.",
#                                                labels=["device", "host"])
#         writes_completed_total = CounterMetricFamily("node_memory_writes_completed_total",
#                                                      "The total number of writes completed successfully.",
#                                                      labels=["device", "host"])
#         written_bytes_total = CounterMetricFamily("node_memory_written_bytes_total",
#                                                   "The total number of bytes written successfully.",
#                                                   labels=["device", "host"])
#         # io_now = Gauge("node_memory_io_now",
#         #               "The number of I/Os currently in progress.",
#         #               ["device", "host"])
#         # io_time_seconds_total = Counter("node_memory_io_time_seconds_total",
#         #                                "Total seconds spent doing I/Os.",
#         #                                ["device", "host"])
#         di_counters = psutil.disk_io_counters(perdisk=True)
#         if not di_counters:
#             return []
#         for k, v in di_counters.items():
#             read_completed_counter.add_metric([k, self.hostname], value=v.read_count)
#             read_bytes_total.add_metric([k, self.hostname], value=v.read_bytes)
#             writes_completed_total.add_metric([k, self.hostname], value=v.write_count)
#             written_bytes_total.add_metric([k, self.hostname], value=v.write_bytes)
#             # io_now.labels(k, self.hostname).inc(v["read_count"])
#             # io_time_seconds_total.labels(k, self.hostname).inc(v["busy_time"])
#         return [read_completed_counter, read_bytes_total, writes_completed_total, written_bytes_total]
#
#
# class MemInfoCollector(Collector):
#     def __init__(self):
#         self.hostname = socket.gethostname()
#
#     def collect(self) -> Iterable[Metric]:
#         available = GaugeMetricFamily("node_memory_available_bytes",
#                                       "Memory information field available.",
#                                       labels=["host"])
#         total = CounterMetricFamily("node_memory_total_bytes",
#                                     "Memory information field total.",
#                                     labels=["host"])
#         active = GaugeMetricFamily("node_memory_active_bytes",
#                                    "Memory information field active.",
#                                    labels=["host"])
#         memory_view = psutil.virtual_memory()
#         if not memoryview:
#             return []
#         available.add_metric([self.hostname], value=memory_view.available)
#         total.add_metric([self.hostname], value=memory_view.total)
#         active.add_metric([self.hostname], value=memory_view.active)
#         ret = []
#         ret.extend([available, active, total])
#         return ret
#
#
# class NetworkCollector(Collector):
#     def __init__(self):
#         self.hostname = socket.gethostname()
#
#     def collect(self) -> Iterable[Metric]:
#         ni_counters = psutil.net_io_counters(pernic=True)
#         if not ni_counters:
#             return []
#         recv = CounterMetricFamily("node_network_receive",
#                                    "Network device statistic receive.",
#                                    labels=["device", "host"])
#         transmit = GaugeMetricFamily("node_network_transmit",
#                                      "Network device statistic transmit.",
#                                      labels=["device", "host"])
#         for k, v in ni_counters.items():
#             recv.add_metric([k, self.hostname], value=v.bytes_recv)
#             transmit.add_metric([k, self.hostname], value=v.bytes_sent)
#         ret = []
#         ret.extend([recv, transmit])
#         return ret
#
#
# CPU_COLLECTOR = CPUCollector()
# MEM_COLLECTOR = MemInfoCollector()
# DISK_COLLECTOR = DiskStatsCollector()
# NETWORK_COLLECTOR = NetworkCollector()
#
# if __name__ == "__main__":
#     cpu_collect = CPUCollector()
#     mem_collect = MemInfoCollector()
#     disk_collect = DiskStatsCollector()
#     network_collect = NetworkCollector()
#     for e in cpu_collect.collect():
#
#     for e in mem_collect.collect():
#
#     for e in disk_collect.collect():
#
#     for e in network_collect.collect():
#
