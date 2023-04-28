/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {spi} from "@/macro/spi";
import {mpi} from "@/macro/mpi";
import type {Context} from "@/prsim/context";
import {Captcha} from "@/kinds/captcha";
import {Status} from "@/cause";
import {index} from "@/macro";

@spi("mesh")
export abstract class Graphics {

    /**
     * Apply a graphics captcha.
     */
    @mpi("mesh.graphics.captcha.apply", Captcha)
    captcha(ctx: Context, @index(0, 'kind') kind: string, @index(1, 'features') features: Map<string, string>): Promise<Captcha> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

    /**
     * Verify a graphics captcha value.
     */
    @mpi("mesh.graphics.captcha.verify", Boolean)
    verify(ctx: Context, @index(0, 'mno') mno: string, @index(1, 'value') value: string): Promise<boolean> {
        return Promise.reject(Status.URN_NOT_PERMIT)
    }

}