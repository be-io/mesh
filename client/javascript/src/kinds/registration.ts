/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Registration {

    @index(0)
    public instance_id: string = "";
    @index(5)
    public name: string = "";
    @index(10)
    public kind: string = "";
    @index(15)
    public address: string = "";
    @index(20)
    public content: any = "";
    @index(25)
    public timestamp: number = 0;
    @index(30)
    public attachments: Map<string, string> = new Map<string, string>();

}