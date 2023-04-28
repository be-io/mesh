#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import json
import os
import platform
from typing import List, Dict

import mesh.tool as tool
from mesh.kinds import Versions
from mesh.macro import spi, mps, ServiceLoader
from mesh.prsim import Builtin, Hodor


@mps
@spi("mesh")
class MeshBuiltin(Builtin):

    def doc(self, name: str, formatter: str) -> str:
        return ""

    def version(self) -> Versions:
        def prop() -> Dict[str, str]:
            try:
                with open(os.path.join(tool.pwd(), f"{tool.get_mesh_name()}.version")) as vfd:
                    vps = json.load(vfd)
                    return vps if vps else {}
            except FileNotFoundError:
                return {}

        properties = prop()
        name = f"{tool.get_mesh_name()}_version".upper()
        mv = tool.get_property("", [name])
        version = tool.anyone(mv, properties.get("version", "1.0.0"))
        versions = Versions()
        versions.version = version
        versions.infos = {
            f"{tool.get_mesh_name()}.commit_id": properties.get("commit_id", "3dd81bc"),
            f"{tool.get_mesh_name()}.os": platform.system(),
            f"{tool.get_mesh_name()}.arch": platform.architecture()[0],
            f"{tool.get_mesh_name()}.version": version,
        }
        return versions

    def debug(self, features: Dict[str, str]):
        doors = ServiceLoader.load(Hodor).list('')
        if not doors:
            return
        for hodor in doors:
            hodor.debug(features)

    def stats(self, features: List[str]) -> Dict[str, str]:
        indies = {}
        doors = ServiceLoader.load(Hodor).list('')
        if not doors:
            return indies
        for hodor in doors:
            stats = hodor.stats(features)
            if not stats:
                continue
            for key, value in stats.items():
                indies[key] = value
        return indies
