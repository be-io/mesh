#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#


import os
from glob import glob
from itertools import chain
from os.path import join

import pytest

import mesh.log as log


def call_pytest_subprocess(*tests, **kwargs):
    global _ctr

    _ctr += 1
    file_name = 'test_result.%d.xml' % _ctr
    if os.path.isfile(file_name):
        os.unlink(file_name)
    env = {}
    args = [
        '--verbose',
        '--cov-append',
        '--cov-report=',
        '--cov', 'spyne',
        '--tb=line',
        '--junitxml=%s' % file_name
    ]
    args.extend('--{0}={1}'.format(k, v) for k, v in kwargs.items())
    return call_test(pytest.main, args, tests, env)


def call_test(self, fn, a, tests, env={}):
    from multiprocessing import Process, Queue

    tests_dir = os.path.dirname(self.__file__)
    if len(tests) > 0:
        a.extend(chain(*[glob(join(tests_dir, test)) for test in tests]))

    queue = Queue()
    p = Process(target=_wrapper(fn), args=[a, queue, env])
    p.start()
    p.join()

    ret = queue.get()
    if ret == 0:
        log.info(tests or a, "OK")
    else:
        log.info(tests or a, "FAIL")

    return ret


def _wrapper(fn):
    import traceback

    def _(args, queue, env):
        log.info("env:", env)
        for k, v in env.items():
            os.environ[k] = v
        try:
            ret = fn(args)
        except SystemExit as e:
            ret = e.code
        except BaseException as e:
            log.info(traceback.format_exc())
            ret = 1

        queue.put(ret)

    return _


class RunTests:
    def run_tests(self):
        ret = 0
        tests = [
            'interface', 'model', 'multipython', 'protocol', 'util',
            'interop/test_pyramid.py',
            'interop/test_soap_client_http_twisted.py',
            'test_null_server.py',
            'test_service.py',
            'test_soft_validation.py',
            'test_sqlalchemy.py',
            'test_sqlalchemy_deprecated.py',
        ]

        log.info("Test stage 1: Unit tests")
        ret = call_pytest_subprocess(*tests, capture=self.capture) or ret

        log.info("\nTest stage 2: End-to-end tests")
        ret = call_pytest_subprocess('interop/test_httprpc.py', capture=self.capture) or ret
        ret = call_pytest_subprocess('interop/test_soap_client_http.py', capture=self.capture) or ret
        ret = call_pytest_subprocess('interop/test_soap_client_zeromq.py', capture=self.capture) or ret
        raise SystemExit(ret)


__all__ = (
    RunTests
)
