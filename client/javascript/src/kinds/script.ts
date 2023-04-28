/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export const EXPRESSION = "EXPRESSION";
export const VALUE = "VALUE";
export const SCRIPT = "SCRIPT"

export class Script {

    @index(0)
    public code: string = ""
    @index(5)
    public name: string = "";
    @index(10)
    public desc: string = "";
    @index(15)
    public kind: string = "";
    @index(20)
    public expr: string = "";
    @index(25)
    public attachment: Map<string, string> = new Map();

}