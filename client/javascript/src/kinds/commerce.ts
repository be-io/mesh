/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro";
import {License} from "@/kinds/license";
import {Environ} from "@/kinds/environ";

export class CommerceLicense {
    @index(0)
    public cipher: string = "";
    @index(5)
    public explain: License = new License();
}

export class CommerceEnviron {
    @index(0)
    public cipher: string = "";
    @index(5)
    public explain: Environ = new Environ();
    @index(10)
    public node_key: string = "";
}