/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import '@/mesh.css';
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Header from "@/widget/header";
import Footer from "@/widget/footer";
import {Box} from "@mui/material";
import Home from "@/views/home/home";
import Citizens from "@/widget/router"

function Mesh() {
    return (
        <BrowserRouter>
            <Box sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                flexDirection: 'column'
            }}>
                <Box sx={{width: '100%'}}>
                    <Header items={Citizens}/>
                </Box>
                <Box sx={{width: '100%'}}>
                    <Box sx={{margin: "auto 6% auto 6%"}}>
                        <Routes>
                            <Route path="/" element={<Home/>}/>
                            <Route path="/mesh" element={<Home/>}/>
                            {
                                Citizens.map(c => c.route)
                            }
                        </Routes>
                    </Box>
                </Box>
                <Box sx={{width: '100%'}}>
                    <Footer/>
                </Box>
            </Box>
        </BrowserRouter>
    );
}

export default Mesh;