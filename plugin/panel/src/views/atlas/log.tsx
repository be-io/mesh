/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import * as React from 'react';
import {useEffect, useState} from 'react';
import {Box, Chip, TextField, Typography} from '@mui/material';
import ManageSearchIcon from '@mui/icons-material/ManageSearch';
import {context, Paging} from "@mesh/mesh";
import services from "@/services/service";

export default function Log() {

    useEffect(() => {
        getMetadata();
    }, []);

    const [labels, setLabels] = useState<string[]>([]);
    const [keywords, setKeywords] = useState('');
    const [blocks, setBlocks] = useState([]);
    const [index, setIndex] = useState(0);

    const getMetadata = () => {
        const paging = new Paging();
        paging.factor.set("mode", "label");
        services.datahouse.read(context(), paging).then(page => {
            // @ts-ignore
            setLabels(page.data)
        }).catch(e => {
            console.log(e);
        });
    }

    const getText = () => {
        const ks = keywords.split("=");
        const paging = new Paging();
        paging.factor.set("mode", "text");
        paging.factor.set(ks[0], ks[1]);
        paging.limit = 4999;
        services.datahouse.read(context(), paging).then(page => {
            // @ts-ignore
            setBlocks(page.data);
        }).catch(e => {
            console.log(e);
        });
    }

    const onKeywordsChange = (e: { target: { value: React.SetStateAction<string>; }; }) => {
        setKeywords(e.target.value);
    }

    const onKeywordsClick = (label: string) => {
        setKeywords(`${label}=`);
    }

    const onKeyboardDown = (e: { key: string; }) => {
        if (e.key === 'Enter') {
            getText();
        }
    }

    const onBlockClick = (idx: number) => {
        setIndex(idx);
    };
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
                    flexDirection: 'column',
                    margin: '20px 10px 20px 20px'
                }}>
                    <Box sx={{display: 'flex', alignItems: 'flex-end', marginLeft: '40px'}}>
                        <ManageSearchIcon sx={{color: 'action.active', mr: 1, my: 0.5}}/>
                        <TextField size="small" label="日志关键字" variant="standard" sx={{width: '40vw'}}
                                   value={keywords} onChange={onKeywordsChange} onKeyDown={onKeyboardDown}/>
                    </Box>
                    <Box>
                        {
                            labels.map(label => {
                                return <Chip sx={{margin: '10px 10px auto auto'}} size="small" key={label} label={label}
                                             onClick={() => onKeywordsClick(label)} color="primary" variant="outlined"/>
                            })
                        }
                    </Box>
                </Box>
                <Box>
                    <Typography paragraph={true} align="left" variant="caption"
                                color="text.secondary" sx={{margin: 'auto 10px auto 10px'}}>
                        {
                            blocks.map((block, idx) => {
                                // @ts-ignore
                                return block.values.map((vs, i) => {
                                    return (
                                        <Box key={idx * 10000 + i}>
                                            {vs[1]}
                                        </Box>
                                    )
                                })
                            })
                        }
                    </Typography>
                </Box>
            </Box>
            <Box>
            </Box>
        </Box>
    );
}