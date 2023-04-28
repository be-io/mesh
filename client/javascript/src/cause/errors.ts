/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Codeable} from "@/cause/codeable";
import {Status} from "@/cause/status";

export function errorf(e: any): Codeable {
    console.log(e);
    return Status.SYSTEM_ERROR;
}