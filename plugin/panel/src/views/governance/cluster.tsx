/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import * as React from 'react';
import {Box, Grid} from '@mui/material';
import Sets from "@/widget/sets";
import Block, {SetOption} from "@/widget/block";

export default function Cluster(props: {}) {

    const boxSize = 5;
    const items: SetOption[][] = [];
    for (let index = 0; index < Sets.length; index++) {
        if (index % boxSize === 0) {
            items.push([])
        }
        items[Math.floor(index / boxSize)].push(Sets[index])
    }

    return (
        <Box sx={{flexGrow: 1}}>
            {
                items.map((apps, index) => {
                    return (
                        <Grid key={index} container spacing={0} columns={15}>
                            {
                                apps.map(app => {
                                    return <Block key={app.name} option={app}/>
                                })
                            }
                        </Grid>
                    )
                })
            }
        </Box>
    );
}