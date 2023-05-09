/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import * as React from 'react';
import {createRef, useCallback, useEffect, useState} from 'react';
import {Box, TextField} from '@mui/material';
import ManageSearchIcon from '@mui/icons-material/ManageSearch';
import {Markmap} from "markmap-view";
import {Transformer} from "markmap-lib";
import services from "@/services/service";
import {context, Paging} from "@mesh/mesh";

const transformer = new Transformer();

export default function Trace() {

    const refSVG = createRef<SVGSVGElement>();
    const [traceId, setTraceId] = useState('');
    const [mm, setMM] = useState<Markmap>();

    useEffect(() => {
        if (mm) {
            return;
        }
        // @ts-ignore
        setMM(Markmap.create(refSVG.current));
    }, [refSVG.current]);

    const onTraceIDChange = useCallback((e: { target: { value: React.SetStateAction<string>; }; }) => {
        setTraceId(e.target.value);
    }, []);

    const onKeyboardDown = useCallback((e: { key: string; }) => {
        if (e.key === 'Enter') {
            getTraces(traceId)
        }
    }, []);

    const getTraces = (traceId: string) => {
        const paging = new Paging();
        paging.factor.set("mode", "trace");
        services.datahouse.read(context(), paging).then(page => {
            // @ts-ignore
            const {root} = transformer.transform(page.data);
            mm && mm.setData(root);
            mm && mm.fit();
        }).catch(e => {
            console.log(e);
        });
    }

    // @ts-ignore
    return (
        <Box sx={{width: '100%', paddingRight: '20px'}}>
            <Box sx={{
                backgroundColor: 'rgb(255,255,255)',
                display: 'flex',
                justifyContent: 'flex-start',
                alignItems: 'center',
                flexDirection: 'column',
            }}>
                <Box sx={{
                    display: 'flex',
                    justifyContent: 'flex-start',
                    alignItems: 'center',
                    margin: '20px 10px 20px 20px'
                }}>
                    <Box sx={{display: 'flex', alignItems: 'flex-end', marginLeft: '40px'}}>
                        <ManageSearchIcon sx={{color: 'action.active', mr: 1, my: 0.5}}/>
                        <TextField size="small" label="TraceID | TaskID" variant="standard" sx={{width: '40vw'}}
                                   value={traceId}
                                   onChange={onTraceIDChange} onKeyDown={onKeyboardDown}/>
                    </Box>
                </Box>

            </Box>
            <Box>
                <svg ref={refSVG} style={{width: '1300px', height: '900px'}}/>
            </Box>
        </Box>
    );
}