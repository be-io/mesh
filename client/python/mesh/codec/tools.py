#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

"""
Helper functions for pickling and unpickling.  Most functions assist in
determining the type of an object.
"""
from __future__ import absolute_import, division, unicode_literals

import base64
import collections
import inspect
import io
import operator
import sys
import time
import types

import mesh.tool as tool

PY_MAJOR = sys.version_info[0]
PY2 = PY_MAJOR == 2
PY3 = PY_MAJOR == 3
PY3_ORDERED_DICT = PY3 and sys.version_info[1] >= 6  # Python 3.6+

class_types = (type,)
iterator_types = (type(iter('')),)

if PY3:
    from collections.abc import Iterator as abc_iterator

    string_types = (str,)
    numeric_types = (int, float)
    ustr = str
else:
    from collections import Iterator as abc_iterator  # noqa

    builtins = __import__('__builtin__')
    class_types += (types.ClassType,)
    encodebytes = base64.encodestring
    decodebytes = base64.decodestring
    string_types = (builtins.basestring,)
    numeric_types = (int, float, builtins.long)
    queue = __import__('Queue')
    ustr = builtins.unicode


def iterator(class_):
    if PY2 and hasattr(class_, '__next__'):
        class_.next = class_.__next__
    return class_


SEQUENCES = (list, set, tuple)
SEQUENCES_SET = {list, set, tuple}
PRIMITIVES = {ustr, bool, type(None)} | set(numeric_types)
FUNCTION_TYPES = {
    types.FunctionType,
    types.MethodType,
    types.LambdaType,
    types.BuiltinFunctionType,
    types.BuiltinMethodType,
}
NON_REDUCIBLE_TYPES = (
        {
            list,
            dict,
            set,
            tuple,
            object,
            bytes,
        }
        | PRIMITIVES
        | FUNCTION_TYPES
)
NON_CLASS_TYPES = {
                      list,
                      dict,
                      set,
                      tuple,
                      bytes,
                  } | PRIMITIVES


def is_type(obj):
    """Returns True is obj is a reference to a type.

    >>> is_type(1)
    False

    >>> is_type(object)
    True

    >>> class Klass: pass
    >>> is_type(Klass)
    True
    """
    # use "isinstance" and not "is" to allow for metaclasses
    return isinstance(obj, class_types)


def has_method(obj, name):
    # false if attribute doesn't exist
    if not hasattr(obj, name):
        return False
    func = getattr(obj, name)

    # builtin descriptors like __getnewargs__
    if isinstance(func, types.BuiltinMethodType):
        return True

    # note that FunctionType has a different meaning in py2/py3
    if not isinstance(func, (types.MethodType, types.FunctionType)):
        return False

    # need to go through __dict__'s since in py3
    # methods are essentially descriptors

    # __class__ for old-style classes
    base_type = obj if is_type(obj) else obj.__class__
    original = None
    # there is no .mro() for old-style classes
    for subtype in inspect.getmro(base_type):
        original = vars(subtype).get(name)
        if original is not None:
            break

    # name not found in the mro
    if original is None:
        return False

    # static methods are always fine
    if isinstance(original, staticmethod):
        return True

    # at this point, the method has to be an instancemthod or a classmethod
    if not hasattr(func, '__self__'):
        return False
    bound_to = getattr(func, '__self__')

    # class methods
    if isinstance(original, classmethod):
        return issubclass(base_type, bound_to)

    # bound methods
    return isinstance(obj, type(bound_to))


def is_object(obj):
    """Returns True is obj is a reference to an object instance.

    >>> is_object(1)
    True

    >>> is_object(object())
    True

    >>> is_object(lambda x: 1)
    False
    """
    return isinstance(obj, object) and not isinstance(
        obj, (type, types.FunctionType, types.BuiltinFunctionType)
    )


def is_not_class(obj):
    """Determines if the object is not a class or a class instance.
    Used for serializing properties.
    """
    return type(obj) in NON_CLASS_TYPES


def is_primitive(obj):
    """Helper method to see if the object is a basic data type. Unicode strings,
    integers, longs, floats, booleans, and None are considered primitive
    and will return True when passed into *is_primitive()*

    >>> is_primitive(3)
    True
    >>> is_primitive([4,4])
    False
    """
    return type(obj) in PRIMITIVES


def is_enum(obj):
    """Is the object an enum?"""
    return 'enum' in sys.modules and isinstance(obj, sys.modules['enum'].Enum)


def is_dictionary(obj):
    """Helper method for testing if the object is a dictionary.

    >>> is_dictionary({'key':'value'})
    True

    """
    return type(obj) is dict


def is_sequence(obj):
    """Helper method to see if the object is a sequence (list, set, or tuple).

    >>> is_sequence([4])
    True

    """
    return type(obj) in SEQUENCES_SET


def is_list(obj):
    """Helper method to see if the object is a Python list.

    >>> is_list([4])
    True
    """
    return type(obj) is list


def is_set(obj):
    """Helper method to see if the object is a Python set.

    >>> is_set(set())
    True
    """
    return type(obj) is set


def is_bytes(obj):
    """Helper method to see if the object is a bytestring.

    >>> is_bytes(b'foo')
    True
    """
    return type(obj) is bytes


def is_unicode(obj):
    """Helper method to see if the object is a unicode string"""
    return type(obj) is ustr


def is_tuple(obj):
    """Helper method to see if the object is a Python tuple.

    >>> is_tuple((1,))
    True
    """
    return type(obj) is tuple


def is_dictionary_subclass(obj):
    """Returns True if *obj* is a subclass of the dict type. *obj* must be
    a subclass and not the actual builtin dict.

    >>> class Temp(dict): pass
    >>> is_dictionary_subclass(Temp())
    True
    """
    # TODO: add UserDict
    return (
            hasattr(obj, '__class__')
            and issubclass(obj.__class__, dict)
            and type(obj) is not dict
    )


def is_sequence_subclass(obj):
    """Returns True if *obj* is a subclass of list, set or tuple.

    *obj* must be a subclass and not the actual builtin, such
    as list, set, tuple, etc..

    >>> class Temp(list): pass
    >>> is_sequence_subclass(Temp())
    True
    """
    return (
            hasattr(obj, '__class__')
            and issubclass(obj.__class__, SEQUENCES)
            and not is_sequence(obj)
    )


def is_noncomplex(obj):
    """Returns True if *obj* is a special (weird) class, that is more complex
    than primitive data types, but is not a full object. Including:

        * :class:`~time.struct_time`
    """
    if type(obj) is time.struct_time:
        return True
    return False


def is_function(obj):
    """Returns true if passed a function

    >>> is_function(lambda x: 1)
    True

    >>> is_function(locals)
    True

    >>> def method(): pass
    >>> is_function(method)
    True

    >>> is_function(1)
    False
    """
    return type(obj) in FUNCTION_TYPES


def is_module_function(obj):
    """Return True if `obj` is a module-global function

    >>> import os
    >>> is_module_function(os.path.exists)
    True

    >>> is_module_function(lambda: None)
    False

    """

    return (
            hasattr(obj, '__class__')
            and isinstance(obj, (types.FunctionType, types.BuiltinFunctionType))
            and hasattr(obj, '__module__')
            and hasattr(obj, '__name__')
            and obj.__name__ != '<lambda>'
    )


def is_module(obj):
    """Returns True if passed a module

    >>> import os
    >>> is_module(os)
    True

    """
    return isinstance(obj, types.ModuleType)


def is_installed(module):
    """Tests to see if ``module`` is available on the sys.path

    >>> is_installed('sys')
    True
    >>> is_installed('hopefullythisisnotarealmodule')
    False

    """
    try:
        __import__(module)
        return True
    except ImportError:
        return False


def is_list_like(obj):
    return hasattr(obj, '__getitem__') and hasattr(obj, 'append')


def is_iterator(obj):
    return isinstance(obj, abc_iterator) and not isinstance(obj, io.IOBase)


def is_collections(obj):
    try:
        return type(obj).__module__ == 'collections'
    except Exception:
        return False


def is_reducible_sequence_subclass(obj):
    return hasattr(obj, '__class__') and issubclass(obj.__class__, SEQUENCES)


def is_reducible(obj):
    """
    Returns false if of a type which have special casing,
    and should not have their __reduce__ methods used
    """
    # defaultdicts may contain functions which we cannot serialise
    if is_collections(obj) and not isinstance(obj, collections.defaultdict):
        return True
    # We turn off the formatting in order to double the speed of the function.
    # Condensing it into one line seems to save the parser a lot of time.
    # fmt: off
    # pylint: disable=line-too-long
    if type(obj) in NON_REDUCIBLE_TYPES or obj is object or is_dictionary_subclass(obj) or isinstance(obj,
                                                                                                      types.ModuleType) or is_reducible_sequence_subclass(
        obj) or is_list_like(obj) or isinstance(getattr(obj, '__slots__', None), iterator_types) or (
            is_type(obj) and obj.__module__ == 'datetime'):  # noqa: E501
        return False
    # fmt: on
    return True


def in_dict(obj, key, default=False):
    """
    Returns true if key exists in obj.__dict__; false if not in.
    If obj.__dict__ is absent, return default
    """
    return (key in obj.__dict__) if getattr(obj, '__dict__', None) else default


def in_slots(obj, key, default=False):
    """
    Returns true if key exists in obj.__slots__; false if not in.
    If obj.__slots__ is absent, return default
    """
    return (key in obj.__slots__) if getattr(obj, '__slots__', None) else default


def has_reduce(obj):
    """
    Tests if __reduce__ or __reduce_ex__ exists in the object dict or
    in the class dicts of every class in the MRO *except object*.

    Returns a tuple of booleans (has_reduce, has_reduce_ex)
    """

    if not is_reducible(obj) or is_type(obj):
        return (False, False)

    # in this case, reduce works and is desired
    # notwithstanding depending on default object
    # reduce
    if is_noncomplex(obj):
        return (False, True)

    has_reduce = False
    has_reduce_ex = False

    REDUCE = '__reduce__'
    REDUCE_EX = '__reduce_ex__'

    # For object instance
    has_reduce = in_dict(obj, REDUCE) or in_slots(obj, REDUCE)
    has_reduce_ex = in_dict(obj, REDUCE_EX) or in_slots(obj, REDUCE_EX)

    # turn to the MRO
    for base in type(obj).__mro__:
        if is_reducible(base):
            has_reduce = has_reduce or in_dict(base, REDUCE)
            has_reduce_ex = has_reduce_ex or in_dict(base, REDUCE_EX)
        if has_reduce and has_reduce_ex:
            return (has_reduce, has_reduce_ex)

    # for things that don't have a proper dict but can be
    # getattred (rare, but includes some builtins)
    cls = type(obj)
    object_reduce = getattr(object, REDUCE)
    object_reduce_ex = getattr(object, REDUCE_EX)
    if not has_reduce:
        has_reduce_cls = getattr(cls, REDUCE, False)
        if has_reduce_cls is not object_reduce:
            has_reduce = has_reduce_cls

    if not has_reduce_ex:
        has_reduce_ex_cls = getattr(cls, REDUCE_EX, False)
        if has_reduce_ex_cls is not object_reduce_ex:
            has_reduce_ex = has_reduce_ex_cls

    return (has_reduce, has_reduce_ex)


def translate_module_name(module):
    """Rename builtin modules to a consistent module name.

    Prefer the more modern naming.

    This is used so that references to Python's `builtins` module can
    be loaded in both Python 2 and 3.  We remap to the "__builtin__"
    name and unmap it when importing.

    Map the Python2 `exceptions` module to `builtins` because
    `builtins` is a superset and contains everything that is
    available in `exceptions`, which makes the translation simpler.

    See untranslate_module_name() for the reverse operation.
    """
    lookup = dict(__builtin__='builtins', exceptions='builtins')
    return lookup.get(module, module)


def _0_9_6_compat_untranslate(module):
    """Provide compatibility for pickles created with jsonpickle 0.9.6 and
    earlier, remapping `exceptions` and `__builtin__` to `builtins`.
    """
    lookup = dict(__builtin__='builtins', exceptions='builtins')
    return lookup.get(module, module)


def untranslate_module_name(module):
    """Rename module names mention in JSON to names that we can import

    This reverses the translation applied by translate_module_name() to
    a module name available to the current version of Python.

    """
    return _0_9_6_compat_untranslate(module)


def importable_name(cls):
    """
    >>> class Example(object):
    ...     pass

    >>> ex = Example()
    >>> importable_name(ex.__class__) == 'jsonpickle.util.Example'
    True
    >>> importable_name(type(25)) == 'builtins.int'
    True
    >>> importable_name(None.__class__) == 'builtins.NoneType'
    True
    >>> importable_name(False.__class__) == 'builtins.bool'
    True
    >>> importable_name(AttributeError) == 'builtins.AttributeError'
    True

    """
    # Use the fully-qualified name if available (Python >= 3.3)
    origin = tool.get_origin(cls)
    name = origin.__name__ if origin else getattr(cls, '__qualname__', cls.__name__)
    module = translate_module_name(cls.__module__)
    return '{}.{}'.format(module, name)


def get_raw_type(cls):
    origin = tool.get_origin(cls)
    return origin if origin else cls


def b64encode(data):
    """
    Encode binary data to ascii text in base64. Data must be bytes.
    """
    return base64.b64encode(data).decode('ascii')


def b64decode(payload):
    """
    Decode payload - must be ascii text.
    """
    return base64.b64decode(payload)


def b85encode(data):
    """
    Encode binary data to ascii text in base85. Data must be bytes.
    """
    return base64.b85encode(data).decode('ascii')


def b85decode(payload):
    """
    Decode payload - must be ascii text.
    """
    return base64.b85decode(payload)


def itemgetter(obj, getter=operator.itemgetter(0)):
    return ustr(getter(obj))
