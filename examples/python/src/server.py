#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import os
import traceback

import mesh
from mesh import mpi, IOStream, Mesh, KV, Network
from mesh.kinds import Paging


class RemoteService:

    def __int__(self):
        self.io_stream = mpi(IOStream)
        self.kv = mpi(KV)
        self.network = mpi(Network)

    def init_data_asset(self):
        page = Paging()
        page.sid = 'xxx'
        page.index = 0
        page.limit = 4096
        page.factor = {'dataset_id': '', 'dataset_name': ''}
        try:
            Mesh.context().set_attribute(Mesh.REMOTE_NAME, "dsconnector")
            full_filename = "/Users/ducer/com/lanxiang/mesh/examples/python/x.csv"
            os.makedirs(os.path.dirname(full_filename), exist_ok=True)
            with open(full_filename, 'wb') as bs:
                tmp_data = self.io_stream.read(page)
                if tmp_data.data:
                    bs.write(tmp_data.data)
                while tmp_data.next:
                    page.index += 1
                    tmp_data = self.io_stream.read(page)
                    if tmp_data.data:
                        bs.write(tmp_data.data)
        except Exception as e:
            traceback.print_exc()
            raise e


def main():
    os.environ.setdefault("mesh.address", "10.99.45.33:570")
    mesh.start()
    mesh.wait()


if __name__ == "__main__":
    main()
