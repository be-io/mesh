/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import React, {useEffect} from "react";
import AccountTreeIcon from '@mui/icons-material/AccountTree';
import ApprovalIcon from '@mui/icons-material/Approval';
import {Route, Routes, useNavigate} from "react-router-dom";
import Dashboard from "@/views/workflow/dashboard";
import Processes from "@/views/workflow/processes";

export default function Workflow(props: {}) {

    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "dashboard",
        name: '流程中心',
        icon: <AccountTreeIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/workflow/dashboard", {replace: true});
        }
    }, {
        key: "processes",
        name: '审批中心',
        icon: <ApprovalIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/workflow/processes", {replace: true});
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
                    <Route path="/" element={<Dashboard/>}/>
                    <Route path="/dashboard" element={<Dashboard/>}/>
                    <Route path="/processes" element={<Processes/>}/>
                </Routes>
            </Box>
        </Box>
    )
}