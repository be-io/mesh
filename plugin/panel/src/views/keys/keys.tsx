/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import {useEffect} from "react";
import AllInclusiveIcon from '@mui/icons-material/AllInclusive';
import VpnLockIcon from '@mui/icons-material/VpnLock';
import SafetyCheckIcon from '@mui/icons-material/SafetyCheck';
import {Route, Routes, useNavigate} from "react-router-dom";
import License from "@/views/keys/license";
import Issue from "@/views/keys/issue";
import Token from "@/views/keys/token";

export default function Keys(props: {}) {

    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "license",
        name: '许可证书',
        icon: <AllInclusiveIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/keys/license", {replace: true});
        }
    }, {
        key: "issue",
        name: '通信证书',
        icon: <VpnLockIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/keys/issue", {replace: true});
        }
    }, {
        key: "token",
        name: '授权证书',
        icon: <SafetyCheckIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/keys/token", {replace: true});
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
                    <Route path="/" element={<License/>}/>
                    <Route path="/license" element={<License/>}/>
                    <Route path="/issue" element={<Issue/>}/>
                    <Route path="/token" element={<Token/>}/>
                </Routes>
            </Box>
        </Box>
    )
}