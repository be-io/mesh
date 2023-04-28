/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";


export class Topic {

    @index(0)
    public topic: string = "";
    @index(5)
    public code: string = "";
    @index(10)
    public group: string = "";
    @index(15)
    public sets: string = "";


}
