/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import '@/mesh.css';
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Header from "@/widget/header";
import Footer from "@/widget/footer";
import {Box} from "@mui/material";
import {AloneRouter, CitizenRouter, Citizens} from "@/widget/router"

function Mesh() {
    return (
        <BrowserRouter>
            <Box sx={{
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between',
                flexDirection: 'column'
            }}>
                <Routes>
                    <Route path="/*" element={
                        <>
                            <Box sx={{width: '100%'}}>
                                <Header items={Citizens}/>
                            </Box>
                            <Box sx={{width: '100%'}}>
                                <Box sx={{margin: "auto 6% auto 6%"}}>
                                    <CitizenRouter/>
                                </Box>
                            </Box>
                            <Box sx={{width: '100%'}}>
                                <Footer/>
                            </Box>
                        </>
                    }/>
                    <Route path="/mesh/x/*" element={<AloneRouter/>}/>
                </Routes>
            </Box>
        </BrowserRouter>
    );
}

export default Mesh;