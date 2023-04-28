/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import {useState} from "react";
import Menu, {ItemOption} from "@/widget/menu";
import HiveIcon from '@mui/icons-material/Hive';
import Grid4x4Icon from '@mui/icons-material/Grid4x4';
import {Route, Routes, useNavigate} from "react-router-dom";
import SRM from "@/views/governance/srm";
import Shortcut from "@/views/governance/shortcut";


export default function Governance(props: {}) {

    const [registrations, setRegistrations] = useState([]);
    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "shortcut",
        name: '快捷操作',
        icon: <HiveIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/shortcut", {replace: true});
        }
    }, {
        key: "service",
        name: '服务治理',
        icon: <Grid4x4Icon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/service", {replace: true});
        }
    }];
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
                    <Route path="/" element={<Shortcut/>}/>
                    <Route path="/service" element={<SRM/>}/>
                    <Route path="/shortcut" element={<Shortcut/>}/>
                </Routes>
            </Box>
        </Box>
    )
}