/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.tool.LDCable;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
@EqualsAndHashCode
public class Reference implements Serializable, LDCable {

    private static final long serialVersionUID = -1613551358789754322L;
    /**
     * Service urn. Like https://accessible.edge.network.omega.000110000000000000000000000.lx00000000000001.trustbe.cn
     */
    @Index(0)
    private String urn;
    /**
     * Service topic.
     */
    @Index(5)
    private String namespace;
    /**
     * Service topic.
     */
    @Index(10)
    private String name;
    /**
     * Service version.
     */
    @Index(15)
    private String version;
    /**
     * Net protocol.
     */
    @Index(20)
    private String proto;
    /**
     * Serialize protocol.
     */
    @Index(25)
    private String codec;
    /**
     * Service invoke asyncable.
     */
    @Index(30)
    private long flags;
    /**
     * Invoke timeout.
     */
    @Index(35)
    private long timeout;
    /**
     * Invoke retries.
     */
    @Index(40)
    private int retries;
    /**
     * Service node identity.
     */
    @Index(45)
    private String node;
    /**
     * Service inst identity.
     */
    @Index(50)
    private String inst;
    /**
     * Service zone.
     */
    @Index(55)
    private String zone;
    /**
     * Service cluster.
     */
    @Index(60)
    private String cluster;
    /**
     * Service cell.
     */
    @Index(65)
    private String cell;
    /**
     * Service group.
     */
    @Index(70)
    private String group;
    /**
     * Service address.
     */
    @Index(75)
    private String address;
}
