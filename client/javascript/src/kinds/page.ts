/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {index} from "@/macro/idx";

export class Page<T> {
    @index(0)
    public sid: string = "";
    @index(5)
    public index: number = 0;
    @index(10)
    public limit: number = 0;
    @index(15)
    public total: number = 0;
    @index(20)
    public next: boolean = false;
    @index(25)
    public data: T[] = [];

    public pages(): number {
        return Math.ceil(this.total / Math.max(this.limit, 1));
    }
}

export class Paging {
    @index(0)
    public sid: string = "";
    @index(5)
    public index: number = 0;
    @index(10)
    public limit: number = 10;
    @index(15)
    public factor: Map<string, any> = new Map<string, any>();
}

export class LogicTable<T> {

    constructor(values: T[], filter?: (v: T, factor: Map<string, any>) => boolean) {
        this.values = values || [];
        this.filter = filter;
    }

    public values: T[] = [];
    public filter?: (v: T, factor: Map<string, any>) => boolean = (v: T, factor: Map<string, any>) => true;
    public snapshot: Page<T> = new Page<T>();

    public index(index: Paging): Page<T> {
        const factor = index.factor || new Map<string, any>();
        const page = new Page<T>();
        page.index = index.index;
        page.sid = index.sid;
        page.limit = index.limit;
        page.total = 0;
        page.next = false;
        page.data = [];
        for (const v of this.values) {
            if (!this.filter || this.filter(v, factor)) {
                page.total++;
            }
            if (page.total >= (index.index + 1) * index.limit) {
                continue
            }
            if (page.total >= index.index * index.limit) {
                page.data.push(v);
            }
        }
        this.snapshot = page;
        return page;
    }
}