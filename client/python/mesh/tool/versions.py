#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import sys

try:
    if sys.version_info < (3, 8):
        import importlib_metadata as metadata
    else:
        from importlib import metadata
except (ImportError, OSError):
    metadata = None


def _get_version():
    default_version = '0.0.1'
    try:
        version = metadata.version('mesh')
    except (AttributeError, ImportError, OSError):
        version = default_version
    return version


__version__ = _get_version()
