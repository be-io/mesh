/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.grpc;

import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.tool.Tool;
import io.grpc.MethodDescriptor;
import lombok.AllArgsConstructor;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;

/**
 * @author coyzeng@gmail.com
 */
@AllArgsConstructor
public class GrpcMarshaller implements MethodDescriptor.Marshaller<InputStream> {

    @Override
    public InputStream stream(InputStream input) {
        try {
            return new ByteArrayInputStream(Tool.readBytes(input));
        } catch (IOException e) {
            throw new MeshException(e);
        }
    }

    @Override
    public InputStream parse(InputStream input) {
        try {
            return new ByteArrayInputStream(Tool.readBytes(input));
        } catch (IOException e) {
            throw new MeshException(e);
        }
    }
}
