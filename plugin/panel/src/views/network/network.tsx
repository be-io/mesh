/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import React, {useEffect} from "react";
import FormatUnderlinedIcon from '@mui/icons-material/FormatUnderlined';
import DeviceHubIcon from '@mui/icons-material/DeviceHub';
import HttpIcon from '@mui/icons-material/Http';
import AbcIcon from '@mui/icons-material/Abc';
import DnsIcon from '@mui/icons-material/Dns';
import DashboardIcon from '@mui/icons-material/Dashboard';
import HubIcon from '@mui/icons-material/Hub';
import {Route, Routes, useNavigate} from "react-router-dom";
import Board from "@/views/network/board";
import HTTP from "@/views/network/http";
import TCP from "@/views/network/tcp";
import PCN from "@/views/network/pcn";
import UDP from "@/views/network/udp";
import DNS from "@/views/network/dns";
import PNN from "@/views/network/pnn";

export default function Network(props: {}) {

    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "cpn",
        name: '网络节点',
        icon: <HubIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/network/pnn", {replace: true});
        }
    }, {
        key: "cpn",
        name: '计算节点',
        icon: <DeviceHubIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/network/pcn", {replace: true});
        }
    }, {
        key: "http",
        name: 'HTTP协议',
        icon: <HttpIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/network/http", {replace: true});
        }
    }, {
        key: "tcp",
        name: 'TCP协议',
        icon: <AbcIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/network/tcp", {replace: true});
        }
    }, {
        key: "udp",
        name: 'UDP协议',
        icon: <FormatUnderlinedIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/network/udp", {replace: true});
        }
    }, {
        key: "dns",
        name: 'DNS协议',
        icon: <DnsIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/network/dns", {replace: true});
        }
    }, {
        key: "dashboard",
        name: '网络大盘',
        icon: <DashboardIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/network/dashboard", {replace: true});
        }
    }];

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
                    <Route path="/" element={<PNN/>}/>
                    <Route path="/dashboard" element={<Board/>}/>
                    <Route path="/http" element={<HTTP/>}/>
                    <Route path="/tcp" element={<TCP/>}/>
                    <Route path="/udp" element={<UDP/>}/>
                    <Route path="/dns" element={<DNS/>}/>
                    <Route path="/pcn" element={<PCN/>}/>
                    <Route path="/pnn" element={<PNN/>}/>
                </Routes>
            </Box>
        </Box>
    )
}