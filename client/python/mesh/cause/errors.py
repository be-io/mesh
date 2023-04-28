#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.cause.status import Codeable, MeshCode


class MeshException(RuntimeError, Codeable):

    def __init__(self, code: str, message: str):
        self.code = code
        self.message = message

    def __str__(self):
        return f'{self.code}, {self.message}'

    def get_code(self):
        return self.code

    def get_message(self):
        return self.message


class CompatibleException(MeshException):

    def __init__(self, message: str):
        super().__init__(MeshCode.COMPATIBLE_ERROR.get_code(), message)


class NotFoundException(MeshException):

    def __init__(self, message: str):
        super().__init__(MeshCode.NOT_FOUND.get_code(), message)


class ValidationException(MeshException):
    """
    Validate error
    """

    def __init__(self, message: str):
        super().__init__(MeshCode.VALIDATE_ERROR.get_code(), message)


class TimeoutException(MeshException):
    """
    TimeoutException
    """

    def __init__(self, message: str):
        super().__init__(MeshCode.TIMEOUT_ERROR.get_code(), message)


class NoProviderException(MeshException):
    """
    NoProviderException
    """

    def __init__(self, message: str):
        super().__init__(MeshCode.NO_PROVIDER_ERROR.get_code(), message)
