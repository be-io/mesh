#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import redis

import mesh.log as log


def main():
    c = redis.Redis(host='127.0.0.1', port=6379, single_connection_client=True)
    c.delete('x', 'y', 'z')
    with c.pipeline(transaction=True) as p:
        log.info("z")
        p.set('x', 'x')
        p.delete('x')
        log.info(f"{p.execute()}")
        log.info(f"{c.hincrbyfloat('z', 'x', 1)}")
        log.info("x")


if __name__ == "__main__":
    main()
