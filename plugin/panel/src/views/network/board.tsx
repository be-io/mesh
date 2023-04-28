/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box, Grid, Paper} from "@mui/material";
import * as React from "react";
import {useEffect, useState} from "react";
import services from "@/services/service";
import {Codec, context, Mesh, ServiceLoader, Types} from "@mesh/mesh";
import {SetOption} from "@/widget/block";

class Counter {

    public total: number;
    public warnings: number;
    public errors: number;

    constructor(total: number, warnings: number, errors: number) {
        this.total = total;
        this.warnings = warnings;
        this.errors = errors;
    }
}

function Block(props: { app: SetOption }) {
    return (
        <Grid item xs={3} sm={3} md={3}>
            <Box bgcolor="background.default" sx={{margin: 'auto 20px auto auto'}}>
                <Paper elevation={3}>
                </Paper>
            </Box>
        </Grid>
    )
}

export default function Board() {

    const codec = ServiceLoader.load(Codec).getDefault();
    const endpoints = useState<Map<string, string>[]>([]);

    useEffect(() => {
        const ctx = context();
        ctx.setAttribute(Mesh.UNAME, 'mesh.dot.dashboard');
        services.endpoint.fuzzy(ctx, codec.encode('')).then(data => {
            const dict: Map<string, Map<string, Counter>> = codec.decode(data, new Types([Map, String, [Map, String, Counter]]));
            console.log(dict?.get("http")?.get("routers"));
        });
    }, []);

    return (
        <Box>
            <Box>
                {
                    endpoints.map(endpoint => {
                        return (
                            <Box></Box>
                        )
                    })
                }
            </Box>
        </Box>
    )
}