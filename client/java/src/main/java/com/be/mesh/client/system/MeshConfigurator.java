/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.system;

import com.be.mesh.client.annotate.Binding;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.prsim.Subscriber;
import com.be.mesh.client.prsim.TheOne;
import com.be.mesh.client.struct.Event;
import com.be.mesh.client.struct.Profile;
import com.be.mesh.client.tool.Tool;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
@SPI("mesh")
@Binding(topic = "mesh.theone.event", code = "refresh")
public class MeshConfigurator implements TheOne, Subscriber {

    private final TheOne proxy = ServiceProxy.proxy(TheOne.class);
    private final Map<String, List<Watcher>> watchers = new ConcurrentHashMap<>();


    @Override
    public void pain(Profile profile) {
        proxy.pain(profile);
    }

    @Override
    public Profile gain(String dataId) {
        return proxy.gain(dataId);
    }

    @Override
    public void follow(String dataId, Watcher watcher) {
        watchers.computeIfAbsent(dataId, key -> new ArrayList<>()).add(watcher);
    }

    @Override
    public void subscribe(Event event) {
        Profile profile = event.getEntity().tryReadObject(Types.of(Profile.class));
        if (Tool.required(watchers.get(profile.getDataId()))) {
            watchers.get(profile.getDataId()).forEach(x -> x.receive(profile));
        }
    }
}
