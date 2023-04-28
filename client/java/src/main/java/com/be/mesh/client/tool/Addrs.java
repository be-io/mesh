/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.tool;

import lombok.AllArgsConstructor;
import lombok.Data;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.concurrent.ThreadLocalRandom;

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Addrs implements Serializable {

    private static final long serialVersionUID = 4291031754738037162L;
    private final List<StatefulServer> servers = new ArrayList<>();
    private final List<String> availableAddrs = new CopyOnWriteArrayList<>();

    public Addrs(String servers) {
        for (String server : Tool.split(servers, ",")) {
            if (Tool.required(server) && !this.availableAddrs.contains(server)) {
                this.servers.add(new StatefulServer(true, server));
                this.availableAddrs.add(server);
            }
        }
    }

    public String any() {
        List<String> addrs = this.availableAddrs;
        if (Tool.optional(addrs)) {
            return "";
        }
        return addrs.get(ThreadLocalRandom.current().nextInt(addrs.size()));
    }

    public String anyHost() {
        String addr = any();
        if (Tool.contains(addr, ":")) {
            return Tool.split(addr, ":")[0];
        }
        return addr;
    }

    public List<String> many() {
        List<String> addrs = this.availableAddrs;
        if (Tool.optional(addrs)) {
            return Collections.emptyList();
        }
        return new ArrayList<>(addrs);
    }

    public void available(String addr, boolean available) {
        for (StatefulServer server : this.servers) {
            if (Tool.equals(addr, server.getAddress()) && server.isAvailable() != available) {
                server.setAvailable(available);
                if (available && !this.availableAddrs.contains(addr)) {
                    this.availableAddrs.add(addr);
                }
                if (!available && this.availableAddrs.size() > 1) {
                    this.availableAddrs.remove(addr);
                }
                return;
            }
        }
    }

    @Data
    @AllArgsConstructor
    public static class StatefulServer implements Serializable {

        private static final long serialVersionUID = -273985030064725213L;
        private volatile boolean available;
        private String address;
    }
}
