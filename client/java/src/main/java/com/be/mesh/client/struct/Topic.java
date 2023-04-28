/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import com.be.mesh.client.tool.Tool;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;

import java.io.Serializable;

/**
 * @author coyzeng@gmail.com
 */
@Data
@AllArgsConstructor
@NoArgsConstructor
@EqualsAndHashCode
public class Topic implements Serializable {

    private static final long serialVersionUID = 4533289608612533105L;

    public Topic(String topic, String code) {
        this.topic = topic;
        this.code = code;
    }

    @Index(0)
    private String topic;
    @Index(5)
    private String code;
    @Index(10)
    private String group;
    @Index(15)
    private String sets;

    public boolean matches(Topic topic) {
        return null != topic && matches(topic.getTopic(), topic.getCode());
    }

    public boolean matches(String topic, String code) {
        return null != topic && Tool.equals(this.topic, topic) && Tool.equals(this.code, code);
    }

    @Override
    public String toString() {
        return String.format("%s:%s", topic, code);
    }
}
