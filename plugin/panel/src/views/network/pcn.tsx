/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import * as React from 'react';
import {useCallback, useEffect, useState} from 'react';
import {
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Link,
    Pagination,
    Paper,
    Stack,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow
} from '@mui/material';
import services from "@/services/service";
import {context, LogicTable, Paging, Registration} from "@mesh/mesh";
import moment from "moment/moment";
import Autocomplete, {LaText} from "@/widget/autocomplete";

export default function PCN() {

    const [table, setTable] = useState<LogicTable<Registration>>(new LogicTable<Registration>([]));
    const [registrations, setRegistrations] = useState<Registration[]>([]);
    const [marker, setMarker] = useState(false);
    const [options, setOptions] = useState<LaText[]>([]);
    const [optionValues, setOptionValues] = useState<LaText[]>([]);

    useEffect(() => {
        services.registry.export(context(), "server").then(data => {
            const rs = data || [];
            rs.sort((p, r) => p.instance_id.localeCompare(r.instance_id));
            const t = new LogicTable(rs)
            setTable(t);
            const paging = new Paging();
            paging.index = 0;
            paging.limit = 10;
            setRegistrations(t.index(paging).data);
        });
        services.network.getRoutes(context()).then(routes => {
            const lts = (routes || []).map(route => {
                const lt = new LaText();
                lt.name = route.name;
                lt.code = route.inst_id;
                return lt;
            });
            setOptions(lts);
        });
    }, []);

    const onClickPagination = useCallback((_: any, idx: any) => {
        const paging = new Paging();
        paging.index = idx;
        paging.limit = 10;
        setRegistrations(table.index(paging).data);
    }, []);

    const onTagClick = useCallback(() => {
        setMarker(true);
    }, []);

    const onSaveTagClick = useCallback(() => {

    }, []);

    return (
        <Box sx={{width: '100%', paddingRight: '20px'}}>
            <Box sx={{
                backgroundColor: 'rgb(255,255,255)',
                display: 'flex',
                justifyContent: 'flex-start',
                alignItems: 'center',
            }}>
            </Box>
            <Box>
                <TableContainer component={Paper}>
                    <Table aria-label="services">
                        <TableHead>
                            <TableRow>
                                <TableCell align="left">计算节点</TableCell>
                                <TableCell align="left">节点名称</TableCell>
                                <TableCell align="left">节点版本</TableCell>
                                <TableCell align="left">节点地址</TableCell>
                                <TableCell align="left">过期时间</TableCell>
                                <TableCell align="left">节点标签</TableCell>
                                <TableCell align="left">操作</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {registrations.map((row) => (
                                <TableRow
                                    key={`${row.instance_id}.${row.name}`}
                                    sx={{'&:last-child td, &:last-child th': {border: 0}}}
                                >
                                    <TableCell align="left">{row.instance_id}</TableCell>
                                    <TableCell align="left">{row.name}</TableCell>
                                    <TableCell align="left">{row.attachments?.get("version") || '1.5.0.0'}</TableCell>
                                    <TableCell align="left">{row.address}</TableCell>
                                    <TableCell
                                        align="left">{moment(row.timestamp).format('yyyy-MM-DD HH:mm')}</TableCell>
                                    <TableCell align="left">
                                        <Stack direction="row" spacing={2}>
                                            <span>1</span>
                                        </Stack>
                                    </TableCell>
                                    <TableCell align="left">
                                        <Stack direction="row" spacing={1}>
                                            <Link href="#" underline="none" onClick={onTagClick}>标记</Link>
                                        </Stack>
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
                <Box sx={{display: 'flex', justifyContent: 'flex-end', marginTop: '10px'}}>
                    <Pagination variant="outlined" shape="rounded" count={table.snapshot.pages()} color="primary"
                                onChange={onClickPagination}/>
                </Box>
                <Dialog open={marker}>
                    <DialogTitle>标记计算节点</DialogTitle>
                    <DialogContent>
                        <DialogContentText>
                            To subscribe to this website, please enter your email address here. We
                            will send updates occasionally.
                        </DialogContentText>
                        <Box>
                            <Autocomplete id="pcn-marker" options={options} values={optionValues}/>
                        </Box>
                    </DialogContent>
                    <DialogActions>
                        <Button onClick={() => {
                            setMarker(false)
                        }}>取消</Button>
                        <Button onClick={onSaveTagClick}>保存</Button>
                    </DialogActions>
                </Dialog>
            </Box>
        </Box>
    );
}