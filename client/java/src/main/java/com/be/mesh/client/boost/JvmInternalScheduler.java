/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.boost;

import com.be.mesh.client.annotate.Binding;
import com.be.mesh.client.annotate.Bindings;
import com.be.mesh.client.annotate.Listener;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Factory;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.prsim.Scheduler;
import com.be.mesh.client.prsim.Subscriber;
import com.be.mesh.client.struct.Event;
import com.be.mesh.client.struct.Timeout;
import com.be.mesh.client.struct.Topic;
import com.be.mesh.client.tool.Once;
import com.be.mesh.client.tool.UUID;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import java.time.Duration;
import java.time.LocalDateTime;
import java.time.ZonedDateTime;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.*;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI("jvm")
public class JvmInternalScheduler implements Scheduler {

    // Cron
    private final Map<String, LocalDateTime> snapshots = new ConcurrentHashMap<>();
    // Timeout event
    private final Map<String, Timeout> timeouts = new ConcurrentHashMap<>(31);
    // Jvm workers
    private final Map<String, JvmSubscriberWorker> workers = new ConcurrentHashMap<>();
    // JVM subscribers
    private static final Once<Map<Topic, Subscriber>> SUBSCRIBERS = new Once<>();
    // Scheduled executor
    private final ScheduledExecutorService executors;

    /**
     * Creates a new timer with the default thread factory
     * ({@link Executors#defaultThreadFactory()}), default tick duration, and
     * default number of ticks per wheel.
     */
    public JvmInternalScheduler() {
        int coreSize = Runtime.getRuntime().availableProcessors();
        ThreadFactory factory = Executors.defaultThreadFactory();
        this.executors = new ScheduledThreadPoolExecutor(coreSize, factory);
    }

    @Override
    public Set<String> dump() {
        return timeouts.keySet();
    }

    @Override
    public boolean cancel(String taskId) {
        Timeout timeout = timeouts.get(taskId);
        if (null == timeout) {
            return false;
        }
        timeout.setStatus(timeout.getStatus() | 2);
        return true;
    }

    @Override
    public boolean stop(String taskId) {
        Timeout timeout = timeouts.get(taskId);
        if (null == timeout) {
            return false;
        }
        timeout.setStatus(timeout.getStatus() | 4);
        return true;
    }

    @Override
    public void emit(Topic topic) {
        //
    }

    @Override
    public boolean isShutdown() {
        return this.executors.isShutdown();
    }

    @Override
    public void shutdown(Duration duration) throws InterruptedException {
        this.executors.shutdown();
    }

    @Override
    public String timeout(Timeout timeout, Duration duration) {
        String taskId = UUID.getInstance().shortUUID();
        timeouts.put(taskId, timeout);
        this.executors.schedule(new Worker(timeout), duration.toNanos(), TimeUnit.NANOSECONDS);
        return taskId;
    }

    @Override
    public String cron(String cron, Topic binding) {
        CronExpression expression = new CronExpression(cron, true);
        String taskId = UUID.getInstance().shortUUID();
        log.info("Setup cron expression, {}, {}, notify to {}", taskId, cron, binding);
        JVMCronWorker worker = new JVMCronWorker();
        worker.setTaskId(taskId);
        worker.setCron(cron);
        worker.setExpression(expression);
        worker.setTopic(binding);
        worker.setExecutors(executors);
        ZonedDateTime now = ZonedDateTime.now();
        ZonedDateTime next = expression.nextTimeAfter(now);
        this.executors.schedule(worker, next.toEpochSecond() - now.toEpochSecond(), TimeUnit.SECONDS);
        return taskId;
    }

    @Override
    public String period(Duration duration, Topic binding) {
        return workers.computeIfAbsent(String.format("%s:%s", binding.getTopic(), binding.getCode()), key -> {
            long period = duration.toMillis();
            String taskId = UUID.getInstance().shortUUID();
            JvmSubscriberWorker worker = new JvmSubscriberWorker(taskId, binding);
            executors.scheduleAtFixedRate(worker, 0, period, TimeUnit.MILLISECONDS);
            return worker;
        }).getTaskId();
    }

    private static Map<Topic, Subscriber> getJVMSubscribers() {
        return SUBSCRIBERS.get(() -> {
            Map<Topic, Subscriber> subscribers = new HashMap<>();
            ServiceLoader.load(Factory.class).list().stream().flatMap(x -> x.getProvider(Subscriber.class)).forEach(x -> {
                Listener listener = x.getClass().getAnnotation(Listener.class);
                Binding binding = x.getClass().getAnnotation(Binding.class);
                Bindings bindings = x.getClass().getAnnotation(Bindings.class);
                if (null != listener) {
                    subscribers.put(new Topic(listener.topic(), listener.code()), x);
                }
                if (null != binding) {
                    subscribers.put(new Topic(binding.topic(), binding.code()), x);
                }
                if (null != bindings) {
                    for (Binding v : bindings.value()) {
                        subscribers.put(new Topic(v.topic(), v.code()), x);
                    }
                }
            });
            return subscribers;
        });
    }

    @AllArgsConstructor
    private static final class Worker implements Runnable {

        private final Timeout timeout;

        @Override
        public void run() {
            Mesh.contextSafeCaught(() -> {
                try {
                    if ((timeout.getStatus() & 2) == 2) {
                        log.warn("Task {} is cancelled", timeout.getTaskId());
                        return;
                    }
                    if ((timeout.getStatus() & 4) == 4) {
                        log.warn("Task {} is stopped", timeout.getTaskId());
                        return;
                    }
                    if ((timeout.getStatus() & 1) == 1) {
                        log.warn("Task {} is expired, execute now", timeout.getTaskId());
                    }
                    getJVMSubscribers().forEach((topic, sub) -> {
                        if (topic.matches(timeout.getBinding())) {
                            Event event = Event.newInstance(timeout.getEntity().readObject(), timeout.getBinding());
                            sub.subscribe(event);
                        }
                    });
                } catch (Exception e) {
                    log.warn("An exception was thrown by " + Timeout.class.getSimpleName() + '.', e);
                } finally {
                    timeout.setStatus(timeout.getStatus() | 2);
                }
            });
        }
    }

    @Getter
    @AllArgsConstructor
    public static final class JvmSubscriberWorker implements Runnable {

        private final String taskId;
        private final Topic topic;

        @Override
        public void run() {
            Mesh.contextSafeCaught(() -> {
                try {
                    getJVMSubscribers().forEach((binding, sub) -> {
                        if (binding.matches(this.topic)) {
                            Event event = Event.newInstance(null, this.topic);
                            event.setEid(taskId);
                            sub.subscribe(event);
                        }
                    });
                } catch (Throwable e) {
                    log.warn("An exception was thrown by " + Timeout.class.getSimpleName() + '.', e);
                }
            });
        }
    }

    @Data
    public static final class JVMCronWorker implements Runnable {

        private String taskId;
        private String cron;
        private CronExpression expression;
        private Topic topic;
        private ScheduledExecutorService executors;

        @Override
        public void run() {
            Mesh.contextSafeCaught(() -> {
                try {
                    ZonedDateTime now = ZonedDateTime.now();
                    ZonedDateTime next = expression.nextTimeAfter(now);
                    this.executors.schedule(this, next.toEpochSecond() - now.toEpochSecond(), TimeUnit.SECONDS);
                    //
                    JvmSubscriberWorker worker = new JvmSubscriberWorker(this.taskId, this.topic);
                    log.info("Setup cron expression, {}, {}, notify to {}", this.taskId, this.cron, this.topic);
                    worker.run();
                } catch (Throwable e) {
                    log.warn("An exception was thrown by " + Timeout.class.getSimpleName() + '.', e);
                }
            });
        }
    }
}
