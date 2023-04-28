#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from typing import Dict

from mesh.macro import index, serializable


@serializable
class Versions:

    @index(0)
    def version(self) -> str:
        """
        Platform version.

        main.sub.feature.bugfix like 1.5.0.0
        :return:
        """
        return ""

    @index(5)
    def infos(self) -> Dict[str, str]:
        """
        Network modules info.
        :return:
        """
        return {}

    def compare(self, v: str) -> int:
        """
        Is the version less than the appointed version.
        :param v: compared version
        :return: -1 true less than, 0 equals, 1 gather than
        """
        if not self.version or "" == self.version:
            return -1
        if not v or "" == v:
            return 1

        rvs = self.version.split(".")
        vs = v.split(".")
        for idx, value in enumerate(rvs):
            if vs.__len__() <= idx:
                return 0
            if vs[idx] == "*":
                continue
            if rvs[idx] > vs[idx]:
                return 1
            if rvs[idx] < vs[idx]:
                return -1
        return 0
