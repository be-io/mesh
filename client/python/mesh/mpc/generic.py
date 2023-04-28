#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any, Type, Dict

from mesh.cause import MeshCode
from mesh.codec import Codec, Json
from mesh.macro import mpi, ServiceLoader, Inspector, Parameters, Returns, Cause, T
from mesh.mpc.compiler import ATTACHMENTS, CODE, MESSAGE, CAUSE, CONTENT
from mesh.mpc.eden import Eden
from mesh.mpc.invoker import Invocation, Execution


class GenericParameters(Parameters, Dict):

    def map(self) -> Dict[str, Any]:
        return self

    def kind(self) -> type:
        return GenericParameters

    def get_arguments(self) -> [Any]:
        args = []
        for key, value in enumerate(self):
            if ATTACHMENTS == key:
                continue
            args.append(value)
        return args

    def set_arguments(self, arguments: [Any]):
        """ No implements """
        pass

    def get_attachments(self) -> Dict[str, str]:
        if not hasattr(self, ATTACHMENTS):
            return dict()
        attachments = self.get(ATTACHMENTS)
        if type(attachments) == dict:
            return attachments
        codec = ServiceLoader.load(Codec).get(Json)
        return codec.decode(attachments, Dict[str, str])

    def set_attachments(self, attachments: Dict[str, str]):
        self.__setitem__(ATTACHMENTS, attachments)


class GenericReturns(Returns, Dict):

    def get_code(self) -> str:
        code = self.get(CODE, None)
        if code:
            return code
        return MeshCode.UNKNOWN.get_code()

    def set_code(self, code: str):
        self.__setitem__(CODE, code)

    def get_message(self) -> str:
        message = self.get(MESSAGE, None)
        if message:
            return message
        return ""

    def set_message(self, message: str):
        self.__setitem__(MESSAGE, message)

    def get_cause(self) -> Cause:
        cause = self.get(CAUSE, None)
        if not cause:
            return cause
        if cause and type(cause == Cause):
            return cause
        codec = ServiceLoader.load(Codec).get(Json)
        return codec.decode(self[CAUSE], Cause)

    def set_cause(self, cause: Cause):
        self.__setitem__(CAUSE, cause)

    def get_content(self) -> Any:
        return self.get(CONTENT, None)

    def set_content(self, content: Any):
        self.__setitem__(CONTENT, content)


class GenericExecution(Execution):

    def __init__(self, refer: mpi, invocation: Invocation):
        self.refer = refer
        self.invocation = invocation

    def schema(self) -> T:
        eden = ServiceLoader.load(Eden).get_default()
        return eden.refer(self.refer, self.kind(), self.method())

    def inspect(self) -> Inspector:
        return self.invocation.get_inspector()

    def intype(self) -> Type[Parameters]:
        return GenericParameters

    def retype(self) -> Type[Returns]:
        return GenericReturns

    def inflect(self) -> Parameters:
        return GenericParameters()

    def reflect(self) -> Returns:
        return GenericReturns()

    def run(self, invocation: Invocation) -> Any:
        eden = ServiceLoader.load(Eden).get_default()
        execution = eden.refer(self.refer, self.kind(), self.method())
        return execution.run(invocation)
