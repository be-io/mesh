/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import {useState} from "react";
import Menu, {ItemOption} from "@/widget/menu";
import AppsIcon from '@mui/icons-material/Apps';
import DiamondIcon from '@mui/icons-material/Diamond';
import SecurityIcon from '@mui/icons-material/Security';
import AltRouteIcon from '@mui/icons-material/AltRoute';
import HiveIcon from '@mui/icons-material/Hive';
import Grid4x4Icon from '@mui/icons-material/Grid4x4';
import {Route, Routes, useNavigate} from "react-router-dom";
import Cluster from "@/views/governance/cluster";
import SRM from "@/views/governance/srm";
import DRM from "@/views/governance/drm";
import Perm from "@/views/governance/perm";
import Shortcut from "@/views/governance/shortcut";
import Depends from "@/views/governance/depends";


export default function Governance(props: {}) {

    const [registrations, setRegistrations] = useState([]);
    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "cluster",
        name: '集群治理',
        icon: <AppsIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/cluster", {replace: true});
        }
    }, {
        key: "service",
        name: '服务治理',
        icon: <Grid4x4Icon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/service", {replace: true});
        }
    }, {
        key: "drm",
        name: '动态配置',
        icon: <DiamondIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/drm", {replace: true});
        }
    }, {
        key: "perms",
        name: '权限管理',
        icon: <SecurityIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/perms", {replace: true});
        }
    }, {
        key: "depends",
        name: '依赖分析',
        icon: <AltRouteIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/depends", {replace: true});
        }
    }, {
        key: "shortcut",
        name: '快捷操作',
        icon: <HiveIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/governance/shortcut", {replace: true});
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
                    <Route path="/" element={<Cluster/>}/>
                    <Route path="/cluster" element={<Cluster/>}/>
                    <Route path="/service" element={<SRM/>}/>
                    <Route path="/drm" element={<DRM/>}/>
                    <Route path="/perms" element={<Perm/>}/>
                    <Route path="/depends" element={<Depends/>}/>
                    <Route path="/shortcut" element={<Shortcut/>}/>
                </Routes>
            </Box>
        </Box>
    )
}