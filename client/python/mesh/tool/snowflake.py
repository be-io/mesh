import time

import mesh.log as log
from mesh.cause import CompatibleException


class UUID(object):
    def __init__(self, worker_id, data_center_id):
        self.worker_id = worker_id
        self.data_center_id = data_center_id
        # stats
        self.ids_generated = 0

        # Tue, 21 Mar 2006 20:50:14.000 GMT
        self.epoch = 1142974214000
        self.sequence = 0
        self.worker_id_bits = 5
        self.data_center_id_bits = 5
        self.max_worker_id = -1 ^ (-1 << self.worker_id_bits)
        self.max_data_center_id = -1 ^ (-1 << self.data_center_id_bits)
        self.sequence_bits = 12

        self.worker_id_shift = self.sequence_bits
        self.data_center_id_shift = self.sequence_bits + self.worker_id_bits
        self.timestamp_left_shift = self.sequence_bits + self.worker_id_bits + self.data_center_id_bits
        self.sequence_mask = -1 ^ (-1 << self.sequence_bits)
        self.last_timestamp = -1

        # Sanity check for worker_id
        if self.worker_id > self.max_worker_id or self.worker_id < 0:
            raise CompatibleException(f"Worker id can't be greater than {self.max_worker_id} or less than 0")

        if self.data_center_id > self.max_data_center_id or self.data_center_id < 0:
            raise CompatibleException(f"Data center id can't be greater than {self.max_data_center_id} or less than 0")

        log.debug(
            f"Timestamp left shift {self.timestamp_left_shift}, data center id bits {self.data_center_id_bits}, "
            f"worker id bits {self.worker_id_bits}, sequence bits {self.sequence_bits}, worker id {self.worker_id}. ")

    @staticmethod
    def _time_gen():
        return int(time.time() * 1000)

    def _till_next_millis(self, last_timestamp):
        timestamp = self._time_gen()
        while timestamp <= last_timestamp:
            timestamp = self._time_gen()

        return timestamp

    def _next_id(self):
        timestamp = self._time_gen()

        if self.last_timestamp > timestamp:
            log.warn(f"Clock is moving backwards. Rejecting request until {self.last_timestamp}")

        if self.last_timestamp == timestamp:
            self.sequence = (self.sequence + 1) & self.sequence_mask
            if self.sequence == 0:
                timestamp = self._till_next_millis(self.last_timestamp)
        else:
            self.sequence = 0

        self.last_timestamp = timestamp

        new_id = ((timestamp - self.epoch) << self.timestamp_left_shift) | (
                self.data_center_id << self.data_center_id_shift) | (
                         self.worker_id << self.worker_id_shift) | self.sequence
        self.ids_generated += 1
        return new_id

    def get_worker_id(self):
        return self.worker_id

    def get_timestamp(self):
        return self._time_gen()

    def new_id(self):
        return self._next_id()

    def get_datacenter_id(self):
        return self.data_center_id
