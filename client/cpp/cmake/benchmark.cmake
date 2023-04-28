#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

if(gRPC_BENCHMARK_PROVIDER STREQUAL "module")
  set(BENCHMARK_ENABLE_GTEST_TESTS OFF CACHE BOOL "Turn off gTest in gBenchmark")
  if(NOT BENCHMARK_ROOT_DIR)
    set(BENCHMARK_ROOT_DIR ${CMAKE_CURRENT_SOURCE_DIR}/third_party/benchmark)
  endif()
  if(EXISTS "${BENCHMARK_ROOT_DIR}/CMakeLists.txt")
    add_subdirectory(${BENCHMARK_ROOT_DIR} third_party/benchmark)
    if(TARGET benchmark)
      set(_gRPC_BENCHMARK_LIBRARIES benchmark)
    endif()
  else()
    message(WARNING "gRPC_BENCHMARK_PROVIDER is \"module\" but BENCHMARK_ROOT_DIR is wrong")
  endif()
elseif(gRPC_BENCHMARK_PROVIDER STREQUAL "package")
  # Use "CONFIG" as there is no built-in cmake module for benchmark.
  find_package(benchmark REQUIRED CONFIG)
  if(TARGET benchmark::benchmark)
    set(_gRPC_BENCHMARK_LIBRARIES benchmark::benchmark)
  endif()
  set(_gRPC_FIND_BENCHMARK "if(NOT benchmark_FOUND)\n  find_package(benchmark CONFIG)\nendif()")
elseif(gRPC_BENCHMARK_PROVIDER STREQUAL "none")
  # Benchmark is a test-only dependency and can be avoided if we're not building tests.
endif()
