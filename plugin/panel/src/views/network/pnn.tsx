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
    TableRow
} from '@mui/material';
import services from "@/services/service";
import {context, Paging, Route} from "@mesh/mesh";
import moment from "moment/moment";

export default function PNN() {

    const [routes, setRoutes] = useState<Route[]>([]);
    const [index, setIndex] = React.useState(0);
    const [limit, setLimit] = React.useState(10);
    const [total, setTotal] = React.useState(0);

    useEffect(() => {
        indexRoutes();
    }, []);

    const indexRoutes = () => {
        const paging = new Paging();
        paging.index = index;
        paging.limit = limit;
        services.network.index(context(), paging).then(page => {
            setTotal(page.total);
            setRoutes(page.data);
        });
    }

    const onClickPagination = (_: any, idx: any) => {
        setIndex(idx - 1);
        indexRoutes();
    };

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
                                <TableCell align="left">节点编号</TableCell>
                                <TableCell align="left">机构编号</TableCell>
                                <TableCell align="left">节点名称</TableCell>
                                <TableCell align="left">节点地址</TableCell>
                                <TableCell align="left">过期时间</TableCell>
                                <TableCell align="left">节点状态</TableCell>
                                <TableCell align="left">操作</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {routes.map((row) => (
                                <TableRow
                                    key={`${row.inst_id}.${row.name}`}
                                    sx={{'&:last-child td, &:last-child th': {border: 0}}}
                                >
                                    <TableCell align="left">{row.node_id}</TableCell>
                                    <TableCell align="left">{row.inst_id}</TableCell>
                                    <TableCell align="left">{row.name}</TableCell>
                                    <TableCell align="left">{row.address}</TableCell>
                                    <TableCell
                                        align="left">{moment(row.expire_at).format('yyyy-MM-DD HH:mm')}</TableCell>
                                    <TableCell
                                        align="left">{moment(row.expire_at).format('yyyy-MM-DD HH:mm')}</TableCell>
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