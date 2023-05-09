/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import {useEffect} from "react";
import WebhookIcon from '@mui/icons-material/Webhook';
import ApiIcon from '@mui/icons-material/Api';
import {Route, Routes, useNavigate} from "react-router-dom";
import Sdks from "@/views/traits/sdks";
import Cases from "@/views/traits/cases";

;

export default function Traits(props: {}) {

    const navigate = useNavigate();
    const menus: ItemOption[] = [{
        key: "sdk",
        name: '多语言SDK',
        icon: <ApiIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/traits/sdk", {replace: true});
        }
    }, {
        key: "cases",
        name: '测试用例',
        icon: <WebhookIcon color="primary"/>,
        onclick: (_: ItemOption) => {
            navigate("/mesh/traits/cases", {replace: true});
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
                    <Route path="/" element={<Sdks/>}/>
                    <Route path="/sdk" element={<Sdks/>}/>
                    <Route path="/cases" element={<Cases/>}/>
                </Routes>
            </Box>
        </Box>
    )
}