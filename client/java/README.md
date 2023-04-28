# Mesh Java Client

[![Build Status](https://travis-ci.org/ducesoft/babel.svg?branch=master)](https://travis-ci.org/ducesoft/babel)
[![Financial Contributors on Open Collective](https://opencollective.com/babel/all/badge.svg?label=financial+contributors)](https://opencollective.com/babel) [![codecov](https://codecov.io/gh/babel/babel/branch/master/graph/badge.svg)](https://codecov.io/gh/babel/babel)
![license](https://img.shields.io/github/license/ducesoft/babel.svg)

中文版 [README](README_CN.md)

## Introduction

Mesh Java client develop kits base on JDK8.

```xml

<dependency>
    <groupId>com.be.mesh</groupId>
    <artfactId>client</artfactId>
    <version>0.0.21</version>
</dependency>
```

## Features

Any

## Get Started

Java

## loki appender

参考如下配置，可将日志写入到日志系统（目前默认是Loki，使用者不需要关注日志数据源）

```xml
<xml>
    <appender name="LOKI" class="com.be.mesh.client.boost.RemoteAppender">
        <format>
            <message>
                <!-- 以下配置使用默认值可不做修改  -->
                <pattern>l=%level h=${HOSTNAME} c=%logger{20} t=%thread | %msg</pattern>
                <!-- 日志队列长度，默认2048，可调整  -->
                <queueSize>2048</queueSize>
                <!-- 日志数量超限后的默认处理动作，默认丢弃-discard，其他可选值：wait-会阻塞线程  -->
                <messageOnFullPolicy>discard</messageOnFullPolicy>
                <!-- 日志每批上送日志数量，默认1000(条每批)，可调整 -->
                <batchMaxItems>1000</batchMaxItems>
                <!-- 日志每批上送时间间隔，默认5000(5s)，可调整 -->
                <batchInterval>5000</batchInterval>
            </message>
        </format>
    </appender>

    <root level="INFO">
        <appender-ref ref="LOKI"/>
    </root>
</xml>
```


