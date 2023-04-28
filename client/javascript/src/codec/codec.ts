/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi, Type} from "@/macro";

@spi(Codec.JSON)
export abstract class Codec {

    public static JSON: string = "json";
    public static PROTOBUF: string = "protobuf";
    public static XML: string = "xml";
    public static YAML: string = "yaml";
    public static THRIFT: string = "thrift";
    public static MESH: string = "mesh";

    /**
     * Encode object to binary array.
     *
     * @param value object
     * @return binary array
     */
    abstract encode(value: any): Uint8Array;

    /**
     * Decode to binary array to object.
     *
     * @param buffer binary array
     * @param type   object type
     * @param <T>    type
     * @return typed object
     */
    abstract decode<T>(buffer: Uint8Array, type: Type<T>): T;

    /**
     * Encode object to string.
     * @param value
     */
    encodeString(value: any): string {
        return this.stringify(this.encode(value));
    }

    /**
     * Decode object from string.
     * @param value
     * @param type
     */
    decodeString<T>(value: string, type: Type<T>): T {
        return this.decode(this.uint8ify(value), type);
    }

    /**
     * String to Uint8array
     */
    abstract uint8ify<T>(buffer: string): Uint8Array;

    /**
     * Uint8array to string.
     */
    abstract stringify<T>(buffer: Uint8Array): string;
}