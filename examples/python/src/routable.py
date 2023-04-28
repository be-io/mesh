#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import os

from mesh import Network, log, mpi, Routable


class TestRoutable:

    @mpi
    def network(self) -> Routable[Network]:
        """"""
        pass

    def main(self):
        os.environ.setdefault("mesh.address", "10.99.7.33:570")
        x = self.network().any_inst("JG0100000700000000").get_environ()
        log.info(x.inst_id)


if __name__ == '__main__':
    routable = TestRoutable()
    routable.main()
