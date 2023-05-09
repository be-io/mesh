/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import * as React from 'react';
import {useEffect, useState} from 'react';
import {Box, Button, Pagination, TextField} from '@mui/material';
import ManageSearchIcon from '@mui/icons-material/ManageSearch';
import services from "@/services/service";
import {context, Registration} from "@mesh/mesh";

export default function Dashboard() {

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
                alignItems: 'center',
            }}>
                <Box sx={{display: 'flex', alignItems: 'flex-end', marginLeft: '40px'}}>
                    <ManageSearchIcon sx={{color: 'action.active', mr: 1, my: 0.5}}/>
                    <TextField size="small" label="流程" variant="standard" sx={{width: '40vw'}}
                               onChange={onSetsChange}/>
                </Box>
                <Button variant="contained" size="small" sx={{marginLeft: '100px'}}>新建</Button>
            </Box>
            <Box>
                <Box>

                </Box>
                <Box sx={{display: 'flex', justifyContent: 'flex-end', marginTop: '10px'}}>
                    <Pagination variant="outlined" shape="rounded" count={Math.ceil(total / limit)} color="primary"
                                onChange={onClickPagination}/>
                </Box>
            </Box>
        </Box>
    );
}