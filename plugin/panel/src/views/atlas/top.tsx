/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import * as React from 'react';
import {useEffect, useState} from 'react';
import {Box, TextField} from '@mui/material';
import ManageSearchIcon from '@mui/icons-material/ManageSearch';
import service from "@/services/service";
import moment from "moment";
import {context, Service} from "@mesh/mesh";

class ServiceDefinition {

    constructor(instance_id: string, app: string, kind: string, addr: string, name: string, version: string, proto: string, codec: string, asyncable: boolean, timeout: number, retries: number, lang: string, timestamp: string) {
        this.instance_id = instance_id;
        this.app = app;
        this.kind = kind;
        this.addr = addr;
        this.name = name;
        this.version = version;
        this.proto = proto;
        this.codec = codec;
        this.asyncable = asyncable;
        this.timeout = timeout;
        this.retries = retries;
        this.lang = lang;
        this.timestamp = timestamp;
    }

    public instance_id: string
    public app: string
    public kind: string
    public addr: string
    public name: string
    public version: string
    public proto: string
    public codec: string
    public asyncable: boolean
    public timeout: number
    public retries: number
    public lang: string
    public timestamp: string
}

export default function Top() {

    const [allService, setAllService] = useState<ServiceDefinition[]>([]);
    const [services, setServices] = useState<ServiceDefinition[]>([]);
    const [index, setIndex] = React.useState(0);
    const [limit, setLimit] = React.useState(15);
    const [total, setTotal] = React.useState(0);
    const [matcher, setMatcher] = React.useState('');

    const exportServices = () => {
        service.registry.export(context(), "metadata").then(rs => {
            const definitions: ServiceDefinition[] = [];
            const registrations = rs || [];
            registrations.forEach(registration => {
                const ss: Service[] = (registration.content && registration.content.services) || [];
                ss.forEach(ms => {
                    definitions.push({
                        instance_id: registration.instance_id,
                        app: registration.name,
                        kind: registration.kind,
                        addr: registration.address,
                        name: ms.urn.substring(0, ms.urn.length - 66).split('.').reverse().join('.'),
                        version: ms.version,
                        proto: ms.proto,
                        codec: ms.codec,
                        asyncable: ms.asyncable,
                        timeout: ms.timeout,
                        retries: ms.retries,
                        lang: ms.lang,
                        timestamp: moment(registration.timestamp).format('yyyy-MM-DD HH:mm:ss'),
                    })
                })
            });
            definitions.sort((p, r) => p.name.localeCompare(r.name));
            setAllService(definitions);
            refreshService(definitions);
        }).catch(e => {
            console.log(e);
        });
    }
    useEffect(() => {
        exportServices();
    }, []);

    const refreshService = (definitions: ServiceDefinition[]) => {
        const matServices = definitions.filter(s => matcher.length === 0 || s.name.indexOf(matcher) > -1)
        setTotal(matServices.length);
        setServices(matServices.slice(index * limit, Math.min(matServices.length, (index + 1) * limit)));
    }

    const refreshService0 = () => {
        const matServices = allService.filter(s => matcher.length === 0 || s.name.indexOf(matcher) > -1)
        setTotal(matServices.length);
        setServices(matServices.slice(index * limit, Math.min(matServices.length, (index + 1) * limit)));
    }

    const onClickPagination = (event: any, idx: any) => {
        setIndex(idx);
        refreshService0();
    };

    const onChangePagination = (event: { target: { value: string | number; }; }) => {
        setLimit(+event.target.value);
        setIndex(0);
        refreshService0();
    };

    const onSetsChange = (event: { target: { value: string; }; }) => {
        setMatcher(event.target.value);
        refreshService0();
    }
    return (
        <Box sx={{width: '100%', paddingRight: '20px'}}>
            <Box sx={{
                backgroundColor: 'rgb(255,255,255)',
                display: 'flex',
                justifyContent: 'flex-start',
                alignItems: 'center',
            }}>
                <Box sx={{
                    display: 'flex',
                    justifyContent: 'flex-start',
                    alignItems: 'center',
                    margin: '20px 10px 20px 20px'
                }}>
                    <Box sx={{display: 'flex', alignItems: 'flex-end', marginLeft: '40px'}}>
                        <ManageSearchIcon sx={{color: 'action.active', mr: 1, my: 0.5}}/>
                        <TextField size="small" label="日志关键字" variant="standard" sx={{width: '40vw'}}
                                   onChange={onSetsChange}/>
                    </Box>
                </Box>
            </Box>
            <Box>
            </Box>
        </Box>
    );
}