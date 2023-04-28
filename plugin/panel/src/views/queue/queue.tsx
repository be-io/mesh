/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import {useEffect} from "react";
import DataArrayIcon from '@mui/icons-material/DataArray';
import StorageIcon from '@mui/icons-material/Storage';
import {Route, Routes, useNavigate} from "react-router-dom";
import Topic from "@/views/queue/topic";
import Broker from "@/views/queue/broker";

export default function Queue(props: {}) {

    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "topic",
        name: '元数据',
        icon: <DataArrayIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/queue/topic", {replace: true});
        }
    }, {
        key: "broker",
        name: '消息队列',
        icon: <StorageIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/queue/broker", {replace: true});
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
                    <Route path="/" element={<Topic/>}/>
                    <Route path="/topic" element={<Topic/>}/>
                    <Route path="/broker" element={<Broker/>}/>
                </Routes>
            </Box>
        </Box>
    )
}