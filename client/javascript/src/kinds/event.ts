/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";
import {Principal} from "@/kinds/principal";
import {Topic} from "@/kinds/topic";
import {Entity} from "@/kinds/entity";

export const messageVersion: string = "1.0.0";

export class Event {

    @index(0)
    public version: string = "";
    @index(5)
    public tid: string = "";
    @index(10)
    public sid: string = "";
    @index(15)
    public eid: string = "";
    @index(20)
    public mid: string = "";
    @index(25)
    public timestamp: string = "";
    @index(30)
    public source: Principal = new Principal();
    @index(35)
    public target: Principal = new Principal();
    @index(40)
    public binding: Topic = new Topic();
    @index(45)
    public entity: Entity = new Entity();
}