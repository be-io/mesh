import argparse
import os
import threading
import time

import numpy as np


class MeshChannel():
    # from_party,to_party is roleid, from_party_id,to_party_id is partyid
    def __init__(self, mesh_addr, from_party, to_party, taskid, subtaskid, from_party_id, to_party_id, use_mock=False):
        from mesh import Mesh, ServiceLoader, Transport, Metadata, Routable

        transport = Routable.of(ServiceLoader.load(Transport).get("mesh"))
        self.session = transport.with_address(mesh_addr).local().open(taskid, {
            Metadata.MESH_VERSION.key(): '',
            # Metadata.MESH_TPC.key(): '',
            Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
            # Metadata.MESH_TOKEN.key(): '',
            Metadata.MESH_SESSION_ID.key(): taskid,
            # Metadata.MESH_SOURCE_INST_ID.key(): from_party_id,
            Metadata.MESH_TARGET_INST_ID.key(): to_party_id
        })

        self.from_party = from_party
        self.to_party = to_party
        self.subtaskid = subtaskid

        # for mock test
        self.use_mock = use_mock

    def send(self, data: bytes):
        if self.use_mock:
            self.mock_ch.send(data)
        else:
            self.session.push(data, {}, self.subtaskid)

    def recv(self, len=None):
        if self.use_mock:
            data = self.mock_ch.recv(len)
        else:
            data = self.session.pop(5000, self.subtaskid)

        if len is not None:
            assert (len == len(data))
        return data


parser = argparse.ArgumentParser()
parser.add_argument("-r", "--role", type=int, default=-1, choices=[-1, 0, 1],
                    help="role, defalut value is -1, mean run all role")
parser.add_argument("-t", "--type", type=str, default="grpc", choices=["grpc", "mem", "mesh"],
                    help="channel type, defalut value is grpc")
parser.add_argument("-m", "--max_packet_size", type=int,
                    default=1024 * 16, help="max packet size")
parser.add_argument("-e", "--enable_cv", type=int, default=1, choices=[0, 1],
                    help="use condition variables instead spinlock")
parser.add_argument("-p", "--peer_grpc", type=str, default="127.0.0.1:9900",
                    help="grpc addr")
parser.add_argument("-l", "--local_redis", type=str, default="tcp://127.0.0.1:6379",
                    help="redis addr")
parser.add_argument("-k", "--key", type=str, default=str(time.time()),
                    help="time.time()")

args = parser.parse_args()
print("args=", args)


def bytes_rng(seed):
    rng = np.random.default_rng(seed)

    def get_rand_bytes():
        len = rng.integers(0, args.max_packet_size)
        return rng.bytes(len)

    return get_rand_bytes


def bytes_rng_fast(seed):
    rng_init = np.random.default_rng(seed)
    bytes_init = rng_init.bytes(args.max_packet_size)

    rng = np.random.default_rng(seed)

    def get_rand_bytes():
        len = rng.integers(0, args.max_packet_size)
        return bytes_init[0: len]

    return get_rand_bytes


def create_channel(taskid: str, role_from: int, role_to: int):
    mesh_addrs = ["10.99.2.33:570", "10.99.3.33:570"]
    nodeids = ["JG0100000200000000", "JG0100000300000000"]
    return MeshChannel(mesh_addrs[role_from], role_from, role_to, taskid, "sub_taskid_" + args.key, nodeids[role_from],
                       nodeids[role_to], False)


def test_send(taskid: str, seed: int):
    rng = bytes_rng_fast(seed)

    chl = create_channel(taskid, 0, 1)
    while True:
        send_str = rng()
        chl.send(send_str)


def size_to_str(speed):
    if speed < 1000:
        return "{:.2f} B".format(speed)
    elif speed < 1000 * 1000:
        return "{:.2f} KB".format(speed / 1000)
    elif speed < 1000 * 1000 * 1000:
        return "{:.2f} MB".format(speed / 1000 / 1000)
    else:
        return "{:.2f} GB".format(speed / 1000 / 1000 / 1000)


def test_recv(taskid: str, seed: int):
    rng = bytes_rng_fast(seed)
    chl = create_channel(taskid, 1, 0)

    total_size = 0
    cur_size = 0
    t0 = time.time()
    t1 = t0

    while True:
        verify_str = rng()
        # recv_str = chl.recv(len(verify_str))
        recv_str = chl.recv()  # a little slow than with recv(len)
        assert (recv_str == verify_str)
        cur_size += len(recv_str)
        t2 = time.time()
        sec = t2 - t1
        if sec > 1:
            cur_speed = cur_size / sec
            total_size += cur_size
            avg_speed = total_size / (t2 - t0)
            print(
                f"avg: {size_to_str(avg_speed)}/s  cur: {size_to_str(cur_speed)}/s  total: {size_to_str(total_size)}")

            cur_size = 0
            t1 = t2


if __name__ == "__main__":
    os.environ.setdefault("MESH_MODE", "4")
    seed = 0
    taskid = "pygaiachannel_bench-111"
    if args.role == -1:
        t1 = threading.Thread(target=test_send, args=(taskid, seed))
        t2 = threading.Thread(target=test_recv, args=(taskid, seed))
        t1.start()
        t2.start()
        t1.join()
        t2.join()
    elif args.role == 0:
        test_send(taskid, seed)
    elif args.role == 1:
        test_recv(taskid, seed)
    else:
        assert (False)
