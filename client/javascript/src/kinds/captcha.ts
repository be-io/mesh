/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Captcha {

    @index(0)
    public mno: string = ""
    @index(5)
    public kind: string = "";
    @index(10)
    public mime: Uint8Array = new Uint8Array;
    @index(15)
    public text: string = "";

}