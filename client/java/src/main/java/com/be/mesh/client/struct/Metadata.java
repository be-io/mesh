/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.List;

/**
 * Without application, service in network must be unique.
 *
 * @author coyzeng@gmail.com
 */
@Data
public class Metadata implements Serializable {

    private static final long serialVersionUID = -1943465454220506029L;

    @Index(0)
    private List<Reference> references;

    @Index(1)
    private List<Service> services;

}
