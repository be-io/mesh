/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Doc {

    @index(0)
    public metadata: Map<string, string> = new Map<string, string>();
    @index(5)
    public content: string = "";
    @index(10)
    public timestamp: number = 0;

}