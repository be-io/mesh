/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import {useEffect, useState} from "react";
import {context, Route as Router} from "@mesh/mesh";
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import FitbitIcon from '@mui/icons-material/Fitbit';
import VerticalAlignTopIcon from '@mui/icons-material/VerticalAlignTop';
import {Route, Routes, useNavigate} from "react-router-dom";
import Trace from "@/views/atlas/trace";
import Log from "@/views/atlas/log";
import Top from "@/views/atlas/top";
import services from "@/services/service";

export default function Atlas(props: {}) {

    const navigate = useNavigate();
    const [nodes, setNodes] = useState<Router[]>([]);
    const [node, setNode] = useState<Router>(new Router());

    const menus: ItemOption[] = [{
        key: "log",
        name: '日志查询',
        icon: <FitbitIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/atlas/log", {replace: true});
        }
    }, {
        key: "top",
        name: '调用TOP',
        icon: <VerticalAlignTopIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/atlas/top", {replace: true});
        }
    }, {
        key: "trace",
        name: '调用链分析',
        icon: <AccountTreeIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/atlas/trace", {replace: true});
        }
    }];

    useEffect(() => {
        services.network.getRoutes(context()).then(routes => {

        })
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
                    <Route path="/" element={<Trace/>}/>
                    <Route path="/trace" element={<Trace/>}/>
                    <Route path="/log" element={<Log/>}/>
                    <Route path="/top" element={<Top/>}/>
                </Routes>
            </Box>
        </Box>
    )
}