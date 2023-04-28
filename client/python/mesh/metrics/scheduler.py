#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

# import socket
# import threading
# import time
#
# from prometheus_client.registry import REGISTRY, CollectorRegistry
#
# from mesh.kinds import Topic
# from mesh.metrics.collector import CPU_COLLECTOR, MEM_COLLECTOR, DISK_COLLECTOR, NETWORK_COLLECTOR
#
#
# # start_collector: start prom collector thread to gather system info
# def start_collector():
#     COLLECTOR_SCHEDULER.start()
#
#
# class CollectorScheduler(threading.Thread):
#
#     def __init__(self, thread_name='thread-prom-collector', interval=10):
#         threading.Thread.__init__(self, name=thread_name)
#         self.hostname = socket.gethostname()
#         if interval < 1:
#             self.interval = 10
#         else:
#             self.interval = interval
#         topic = Topic()
#         topic.topic = "mesh.prom.sample.export"
#         topic.code = "*"
#         self.bindings = topic
#         # TODO load publisher
#         self.publisher = {}
#         self.default_registry = REGISTRY
#         self.registry = CollectorRegistry(auto_describe=False)
#         self.registry.register(CPU_COLLECTOR)
#         self.registry.register(MEM_COLLECTOR)
#         self.registry.register(DISK_COLLECTOR)
#         self.registry.register(NETWORK_COLLECTOR)
#
#     def run(self):
#         while True:
#             self.collect()
#             # collect per 10s
#             time.sleep(self.interval)
#
#     def collect(self):
#         samples = []
#         self_metrics = self.registry.collect()
#         if self_metrics:
#             for metric in iter(self_metrics):
#                 print(metric)
#                 samples.append(metric)
#
#         metrics = self.default_registry.collect()
#         for metric in metrics:
#             print(metric)
#             samples.append(metric)
#
#         if len(samples) > 0:
#             print(samples)
#             event = {
#                 "metricFamilySamples": samples,
#                 "instanceId": self.hostname,
#                 "jobName": "pushgateway",
#             }
#             try:
#                 # TODO publish event
#                 # self.publish(event)
#                 print(event)
#             except:
#                 print("publisher publish Exception:")
#
#     def publish(self, payload):
#         # TODO active publish
#         # self.publisher.publish_with_topic(self.bindings, payload)
#         pass
#
#
# COLLECTOR_SCHEDULER = CollectorScheduler()
#
# # for local test
# if __name__ == "__main__":
#     scheduler = CollectorScheduler()
#     scheduler.run()
