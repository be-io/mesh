/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Key} from "@/prsim";
import {Invocation} from "@/mpc/invoker";
import {Consumer} from "@/mpc/consumer";
import {Codec} from "@/codec";

export class Mesh {

    /**
     * Mesh mpc remote address.
     */
    public static REMOTE: Key<string> = new Key<string>("mesh.mpc.address");
    /**
     * Remote app name.
     */
    public static REMOTE_NAME: Key<string> = new Key<string>("mesh.mpc.remote.name");
    /**
     * Mesh invocation attributes.
     */
    public static INVOCATION: Key<Invocation> = new Key<Invocation>("mesh.invocation");
    /**
     * Mesh invoke mpi name attributes.
     */
    public static UNAME: Key<string> = new Key<string>("mesh.mpc.uname");

}

class MeshFlag {

    public code: string;
    public name: string;

    constructor(code: string, name: string) {
        this.code = code;
        this.name = name;
    }
}

export class MeshFlags {

    public static HTTP: MeshFlag = new MeshFlag("00", Consumer.HTTP);
    public static GRPC: MeshFlag = new MeshFlag("01", Consumer.GRPC);
    public static MQTT: MeshFlag = new MeshFlag("02", Consumer.MQTT);
    public static TCP: MeshFlag = new MeshFlag("03", Consumer.TCP);

    public static JSON: MeshFlag = new MeshFlag("00", Codec.JSON);
    public static PROTOBUF: MeshFlag = new MeshFlag("01", Codec.PROTOBUF);
    public static XML: MeshFlag = new MeshFlag("02", Codec.XML);
    public static THRIFT: MeshFlag = new MeshFlag("03", Codec.THRIFT);
    public static YAML: MeshFlag = new MeshFlag("04", Codec.YAML);

    private static PROTO = [MeshFlags.HTTP, MeshFlags.GRPC, MeshFlags.MQTT, MeshFlags.TCP];
    private static CODEC = [MeshFlags.JSON, MeshFlags.PROTOBUF, MeshFlags.XML, MeshFlags.THRIFT, MeshFlags.YAML];

    public static ofProto(code: string): MeshFlag {
        for (let proto of MeshFlags.PROTO) {
            if (proto.code == code) {
                return proto;
            }
        }
        return MeshFlags.HTTP;
    }

    public static ofCodec(code: string): MeshFlag {
        for (let codec of MeshFlags.CODEC) {
            if (codec.code == code) {
                return codec;
            }
        }
        return MeshFlags.JSON;
    }

    public static ofName(name: string): MeshFlag {
        for (let proto of MeshFlags.PROTO) {
            if (proto.name == name) {
                return proto;
            }
        }
        for (let codec of MeshFlags.CODEC) {
            if (codec.name == name) {
                return codec;
            }
        }
        return MeshFlags.HTTP;
    }
}