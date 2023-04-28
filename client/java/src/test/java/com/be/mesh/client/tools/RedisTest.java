/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tools;

import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.google.common.collect.ImmutableMap;
import lombok.extern.slf4j.Slf4j;
import org.redisson.Redisson;
import org.redisson.api.RLock;
import org.redisson.api.RedissonClient;
import org.redisson.config.Config;
import org.testng.annotations.Test;
import redis.clients.jedis.Jedis;
import redis.clients.jedis.JedisPubSub;
import redis.clients.jedis.Pipeline;
import redis.clients.jedis.Response;
import redis.clients.jedis.params.ScanParams;

import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.time.Duration;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

/**
 * get/set/del hset/hget/hdel hmset/hmget incr/decr exists publish/subscribe setex/setnx expire ping lpush/lpop brpop pubsub redlock
 *
 * @author coyzeng@gmail.com
 */
@Slf4j
public class RedisTest extends JedisPubSub {

    @Test
    public void testExec() {
        Codec codec = ServiceLoader.load(Codec.class).getDefault();
        try (Jedis jedis = new Jedis(URI.create("redis://127.0.0.1:6379"))) {
            log.info("hincrby {}", jedis.hincrBy(inBytes("h1"), inBytes("v"), 1));
            log.info("hincrby {}", jedis.hincrBy(inBytes("h1"), inBytes("v"), 1));

            log.info("hincrbyfloat {}", jedis.hincrByFloat(inBytes("f1"), inBytes("v"), 1));
            log.info("hincrbyfloat {}", jedis.hincrByFloat(inBytes("f1"), inBytes("v"), 1));

            log.info("hkeys {}", codec.encodeString(jedis.hkeys(inBytes("f1"))));

            log.info("hgetall {}", codec.encodeString(jedis.hgetAll(inBytes("f1"))));

            log.info("scan-1 {}", jedis.scan("1"));
            log.info("scan-2 {}", jedis.scan(inBytes("1"), new ScanParams().count(1).match(inBytes("m1")).match("m2")));

            log.info("hexists {}", jedis.hexists(inBytes("k1"), inBytes("v")));

            log.info("hexists {}", jedis.hset(inBytes("v1"), inBytes("k"), inBytes("v")));
            log.info("hexists {}", jedis.hexists(inBytes("v1"), inBytes("k")));

            log.info("lpush {}", jedis.lpush(inBytes("k1"), inBytes("v")));
            log.info("blpop {}", forBytes(jedis.blpop(Duration.ofSeconds(10).getSeconds(), inBytes("k1"))));

            log.info("set {}", jedis.set(inBytes("x"), inBytes("y")));
            log.info("get {}", forBytes(jedis.get(inBytes("x"))));
            log.info("del {}", jedis.del(inBytes("x")));

            log.info("hset {}", jedis.hset("hx", ImmutableMap.of("x", "y")));
            log.info("hget {}", jedis.hget("hx", "x"));
            log.info("hdel {}", jedis.hdel("hx", "x"));

            log.info("hmset {}", jedis.hmset("hmx", ImmutableMap.of("x", "y")));
            log.info("hmget {}", jedis.hmget("hmx", "x"));

            log.info("incr {}", jedis.incr("z"));
            log.info("decr {}", jedis.decr("z"));

            log.info("exists {}", jedis.exists("z"));

            log.info("setex {}", jedis.setex("setex", Duration.ofSeconds(10).getSeconds(), "setex"));
            log.info("setnx {}", jedis.setnx("setnx", "setnx"));

            log.info("expire {}", jedis.expire("expire", Duration.ofSeconds(10).getSeconds()));

            log.info("ping {}", jedis.ping());

            log.info("lpush {}", jedis.lpush("lpush", "x", "y"));
            log.info("lpop {}", jedis.lpop("lpush"));

            log.info("rpush {}", jedis.rpush("rpush", "x", "y", "z"));
            log.info("rpop {}", jedis.rpop("rpush"));
            log.info("brpop {}", jedis.brpop(1, "rpush", "lpop"));

            Pipeline pipeline = jedis.pipelined();
            Response<String> pset = pipeline.set("a", "b");
            Response<String> pget = pipeline.get("a");
            Response<Long> pdel = pipeline.del("a", "b", "c", "lpush", "lpop");
            pipeline.close();
            log.info("pip set {}", pset.get());
            log.info("pip get {}", pget.get());
            log.info("pip del {}", pdel.get());

            Config config = new Config();
            config.useSingleServer().setAddress("redis://127.0.0.1:6379");
            RedissonClient client = Redisson.create(config);
            RLock lock = client.getLock("xxx");
            if (lock.tryLock()) {
                log.info("Redisson lock");
                lock.unlock();
            }
            log.info("publish {}", jedis.publish("publish", "channel-x"));
            jedis.subscribe(this, "publish");
            log.info("subscribe {}", "channel-x");
        }
    }

    @Test
    public void testExecInSingle() {
        // "redis://:redis123@10.12.0.206:6379"
        try (Jedis jedis = new Jedis(URI.create("redis://127.0.0.1:6379"))) {
            log.info("publish {}", jedis.publish("channel-x", "x"));
            jedis.subscribe(this, "channel-x", "channel-y", "channel-z");
            log.info("subscribe {} {} {}", "channel-x", "channel-y", "channel-z");
        }
    }

    private byte[] inBytes(String x) {
        return x.getBytes(StandardCharsets.UTF_8);
    }

    private String forBytes(byte[] x) {
        return new String(x, StandardCharsets.UTF_8);
    }

    private String forBytes(List<byte[]> xs) {
        return Optional.ofNullable(xs).orElseGet(Collections::emptyList).stream().map(x -> new String(x, StandardCharsets.UTF_8)).collect(Collectors.joining(","));
    }

    @Override
    public void onMessage(String channel, String message) {
        log.info("subscribe {}:{}", channel, message);
    }

}
