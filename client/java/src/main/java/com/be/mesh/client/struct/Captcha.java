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

/**
 * @author coyzeng@gmail.com
 */
@Data
public class Captcha implements Serializable {

    private static final long serialVersionUID = -6603795112981102318L;
    @Index(0)
    private String mno;
    @Index(5)
    private String kind;
    @Index(10)
    private byte[] mime;
    @Index(15)
    private String text;
}
