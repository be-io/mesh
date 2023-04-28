/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import type {Context} from "@/prsim/context";

export interface Runtime {

}

export interface RuntimeHook {

    /**
     * Start Trigger when mesh runtime is start.
     */
    start(ctx: Context, runtime: Runtime): Promise<void>;

    /**
     * Stop Trigger when mesh runtime is stop.
     */
    stop(ctx: Context, runtime: Runtime): Promise<void>;

    /**
     * Refresh Trigger then mesh runtime context is refresh or metadata is refresh.
     */
    refresh(ctx: Context, runtime: Runtime): Promise<void>;

}