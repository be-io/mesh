#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import threading
from typing import List, Dict

import mesh.log as log
import mesh.tool as tool
from mesh.kinds import Topic, Timeout, Event
from mesh.macro import spi, ServiceLoader, Binding
from mesh.mpc import ServiceProxy
from mesh.prsim import Scheduler, Subscriber


@spi("python")
class PythonScheduler(Scheduler):

    def __init__(self):
        self.remote = ServiceProxy.default_proxy(Scheduler)
        self.tasks: Dict[str, Task] = {}

    def timeout(self, timeout: Timeout, duration: int) -> str:
        return ''

    def cron(self, cron: str, binding: Topic) -> str:
        return ''

    def period(self, duration: int, binding: Topic) -> str:
        task_id = tool.next_id()
        log.info(f"Next task {task_id} has been submit.")
        task = Task(duration / 1000, self.do_emit(binding))
        task.start()
        self.tasks[task_id] = task
        return task_id

    def dump(self) -> List[str]:
        task_ids = []
        for task_id, _ in self.tasks.items():
            task_ids.append(task_id)
        return task_ids

    def cancel(self, task_id: str) -> bool:
        return self.stop(task_id)

    def stop(self, task_id: str) -> bool:
        task = self.tasks.get(task_id, None)
        if task:
            task.cancel()
            self.tasks.__delitem__(task_id)
        return True

    def emit(self, topic: Topic) -> bool:
        subscribers = ServiceLoader.load(Subscriber).list('')
        for subscriber in subscribers:
            try:
                bindings = Binding.get_binding_if_present(subscriber)
                if not topic.matches(bindings):
                    continue
                subscriber.subscribe(Event.new_instance({}, topic))
            except BaseException as e:
                log.error(f"{e}")

        return True

    def do_emit(self, topic: Topic):
        def x():
            self.emit(topic)

        return x

    def shutdown(self, duration: int) -> bool:
        for _, task in self.tasks.items():
            task.cancel()
        self.tasks.clear()
        return True


class Task(threading.Thread):

    def __init__(self, interval, function, args=None, kwargs=None):
        threading.Thread.__init__(self)
        self.interval = interval
        self.function = function
        self.args = args if args is not None else []
        self.kwargs = kwargs if kwargs is not None else {}
        self.finished = threading.Event()

    def cancel(self):
        """Stop the timer if it hasn't finished yet."""
        self.finished.set()

    def run(self):
        while not self.finished.isSet():
            try:
                self.function(*self.args, **self.kwargs)
                self.finished.wait(self.interval)
            except BaseException as e:
                log.error(f"{e}")
