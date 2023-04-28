/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import React from 'react';
import ReactDOM from 'react-dom/client';
import {CssBaseline, ThemeProvider} from "@mui/material";
import theme from "@/widget/theme";
import Mesh from '@/mesh';
import banner from "@/widget/banner";
import '@/main.css';
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';


console.log(banner("0.0.1", "000000"))

const root = ReactDOM.createRoot(document.getElementById('root')!);
root.render(
    <ThemeProvider theme={theme}>
        <CssBaseline/>
        <Mesh/>
    </ThemeProvider>
);