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
    TextField
} from '@mui/material';
import ManageSearchIcon from '@mui/icons-material/ManageSearch';
import services from "@/services/service";
import {context, Registration} from "@mesh/mesh";
import moment from "moment/moment";

export default function DNS() {

    const [registrations, setRegistrations] = useState<Registration[]>([]);
    const [filteredRegistrations, setFilteredRegistrations] = useState<Registration[]>([]);
    const [index, setIndex] = React.useState(0);
    const [limit, setLimit] = React.useState(10);
    const [total, setTotal] = React.useState(0);
    const [matcher, setMatcher] = React.useState('');

    useEffect(() => {
        console.log(1)
        services.registry.export(context(), "server").then(data => {
            data.sort((p, r) => p.instance_id.localeCompare(r.instance_id));
            setRegistrations(data);
            const fs = data.filter(s => matcher.length === 0 || s.name.indexOf(matcher) > -1)
            setTotal(fs.length);
            setFilteredRegistrations(fs.slice(index * limit, Math.min(fs.length, (index + 1) * limit)));
        }).catch(e => {
            console.log(e);
        });
    }, []);

    const refreshService = () => {
        const fs = registrations.filter(s => matcher.length === 0 || s.name.indexOf(matcher) > -1)
        setTotal(fs.length);
        setFilteredRegistrations(fs.slice(index * limit, Math.min(fs.length, (index + 1) * limit)));
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
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {filteredRegistrations.map((row) => (
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
                                            <Link href="#" underline="none">标记</Link>
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