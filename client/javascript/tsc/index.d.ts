/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {ParsedCommandLine} from 'typescript';
import {Plugin} from 'vite';
import {LoadResult} from "rollup";

declare class TSC implements Plugin {
    name: string;
    tsx: boolean;
    tsconfig: ParsedCommandLine;

    constructor(tsconfigPath: string, tsx: boolean);

    load(id: string, options?: {
        ssr?: boolean;
    }): (Promise<LoadResult> | LoadResult);

    transform(code: string, id: string, options?: { ssr?: boolean; }): string;

    parseTsConfig(tsconfig: string, cwd?: string): ParsedCommandLine;

    printDiagnostics(...args: any[]): void;
}

export default function tsc(tsconfigPath?: string, tsx?: boolean): TSC;
export {};
