/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import {useEffect} from "react";
import DataObjectIcon from '@mui/icons-material/DataObject';
import DatasetIcon from '@mui/icons-material/Dataset';
import {Route, Routes, useNavigate} from "react-router-dom";
import Mysql from "@/views/store/mysql";
import Cache from "@/views/store/cache";

export default function Store(props: {}) {

    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "cache",
        name: '缓存管理',
        icon: <DataObjectIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/store/cache", {replace: true});
        }
    }, {
        key: "mysql",
        name: '数据库管理',
        icon: <DatasetIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/store/mysql", {replace: true});
        }
    },];

    useEffect(() => {
    }, []);
    return (
        <Box sx={{
            display: 'flex',
            justifyContent: 'flex-start',
            alignItems: 'flex-start',
            backgroundColor: 'rgb(244,246,249)',
            width: '100%'
        }}>
            <Menu menus={menus}/>
            <Box sx={{margin: '10px auto auto 20px', width: '100%'}}>
                <Routes>
                    <Route path="/" element={<Cache/>}/>
                    <Route path="/cache" element={<Cache/>}/>
                    <Route path="/mysql" element={<Mysql/>}/>
                </Routes>
            </Box>
        </Box>
    )
}