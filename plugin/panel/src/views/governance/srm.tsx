/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import * as React from 'react';
import {useEffect, useState} from 'react';
import {
    Box,
    Link,
    Pagination,
    Paper,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    TextField,
    Tooltip,
    Typography
} from '@mui/material';
import ManageSearchIcon from '@mui/icons-material/ManageSearch';
import services from "@/services/service";
import moment from "moment";
import {context, Service, URN} from "@mesh/mesh";

class ServiceInfo {
    public instance_id: string = "";
    public app: string = "";
    public kind: string = "";
    public addr: string = "";
    public name: string = "";
    public version: string = "";
    public proto: string = "";
    public codec: string = "";
    public flags: number = 0;
    public timeout: number = 0;
    public retries: number = 0;
    public lang: string = "";
    public timestamp: string = "";
}

export default function SRM() {

    const [allService, setAllService] = useState<ServiceInfo[]>([]);
    const [resources, setResources] = useState<ServiceInfo[]>([]);
    const [index, setIndex] = React.useState(0);
    const [limit, setLimit] = React.useState(10);
    const [total, setTotal] = React.useState(0);
    const [matcher, setMatcher] = React.useState('');

    const exportServices = () => {
        services.registry.export(context(), "metadata").then(rs => {
            const infos: ServiceInfo[] = [];
            const registrations = rs || [];
            registrations.forEach(registration => {
                const ss: Service[] = (registration.content && registration.content.services) || [];
                ss.forEach(ms => {
                    infos.push({
                        instance_id: registration.instance_id,
                        app: registration.name,
                        kind: registration.kind === 'metadata' ? "MPI" : '代理',
                        addr: registration.address,
                        name: registration.kind === 'metadata' ? URN.from(ms.urn).name : ms.urn,
                        version: ms.version,
                        proto: ms.proto,
                        codec: ms.codec,
                        flags: ms.flags,
                        timeout: ms.timeout,
                        retries: ms.retries,
                        lang: ms.lang,
                        timestamp: moment(registration.timestamp).format('yyyy-MM-DD HH:mm'),
                    })
                })
            });
            infos.sort((p, r) => p.name.localeCompare(r.name));
            setAllService(infos);
            const fs = infos.filter(s => matcher.length === 0 || s.name.indexOf(matcher) > -1)
            setTotal(fs.length);
            setResources(fs.slice(index * limit, Math.min(fs.length, (index + 1) * limit)));
        }).catch(e => {
            console.log(e);
        });
    }
    useEffect(() => {
        exportServices();
    }, []);

    const refreshService = () => {
        const fs = allService.filter(s => matcher.length === 0 || s.name.indexOf(matcher) > -1)
        setTotal(fs.length);
        setResources(fs.slice(index * limit, Math.min(fs.length, (index + 1) * limit)));
    }

    const onClickPagination = (_: any, idx: any) => {
        setIndex(idx - 1);
        refreshService();
    };

    const onSetsChange = (event: { target: { value: string; }; }) => {
        setMatcher(event.target.value);
        refreshService();
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
                        <TextField size="small" label="服务名" variant="standard" sx={{width: '40vw'}}
                                   onChange={onSetsChange}/>
                    </Box>
                </Box>
            </Box>
            <Box>
                <TableContainer component={Paper}>
                    <Table aria-label="services">
                        <TableHead>
                            <TableRow>
                                <TableCell align="left">应用</TableCell>
                                <TableCell align="left">类型</TableCell>
                                <TableCell align="left">地址</TableCell>
                                <TableCell align="left">名称</TableCell>
                                <TableCell align="left">版本</TableCell>
                                <TableCell align="left">协议/编码</TableCell>
                                <TableCell align="left">异步</TableCell>
                                <TableCell align="left">超时（ms）</TableCell>
                                <TableCell align="left">重试</TableCell>
                                <TableCell align="left">语言</TableCell>
                                <TableCell align="left">时间</TableCell>
                                <TableCell align="left">治理</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {resources.map((row) => (
                                <TableRow
                                    key={`${row.instance_id}.${row.name}`}
                                    sx={{'&:last-child td, &:last-child th': {border: 0}}}
                                >
                                    <TableCell size="small" align="left">{row.app}</TableCell>
                                    <TableCell size="small" align="left">{row.kind}</TableCell>
                                    <TableCell size="small" align="left">{row.addr}</TableCell>
                                    <TableCell size="small" align="left">
                                        <Tooltip title={row.name} placement="bottom-start">
                                            <Typography sx={{width: 150}} align="left" variant="caption"
                                                        color="text.secondary"
                                                        noWrap={true} paragraph={true}>
                                                {row.name}
                                            </Typography>
                                        </Tooltip>
                                    </TableCell>
                                    <TableCell size="small"
                                               align="left">{row.version === "" ? '*' : row.version}</TableCell>
                                    <TableCell size="small" align="left">GRPC/JSON</TableCell>
                                    <TableCell size="small" align="left">{row.flags ? '异步' : '同步'}</TableCell>
                                    <TableCell size="small" align="left">{row.timeout}</TableCell>
                                    <TableCell size="small" align="left">{row.retries}</TableCell>
                                    <TableCell size="small" align="left">{row.lang}</TableCell>
                                    <TableCell size="small" align="left">{row.timestamp}</TableCell>
                                    <TableCell size="small" align="left">
                                        <Stack direction="row" spacing={1}>
                                            <Link href="#" underline="none">治理</Link>
                                            <Link href="#" underline="none">测试</Link>
                                        </Stack>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                <Box sx={{display: 'flex', justifyContent: 'flex-end', marginTop: '10px'}}>
                    <Pagination variant="outlined" shape="rounded" count={Math.ceil(total / limit)} color="primary"
                                onChange={onClickPagination}/>
                </Box>
            </Box>
        </Box>
    );
}