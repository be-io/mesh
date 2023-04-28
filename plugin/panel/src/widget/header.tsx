/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {SyntheticEvent, useCallback, useEffect, useState} from 'react';
import {Box, Divider, IconButton, Stack, Tab, Tabs, Tooltip} from '@mui/material';
import GitHubIcon from '@mui/icons-material/GitHub';
import {useNavigate} from "react-router-dom";
import {FavIcon} from "@/widget/icons";
import {Citizen} from '@/widget/router';

function Header(props: { items: Citizen[] }) {

    const [version, setVersion] = useState(0);
    const [index, setIndex] = useState(0);

    const navigate = useNavigate();

    useEffect(() => {
        for (let idx = 0; idx < props?.items?.length; idx++) {
            if (window.location.pathname && window.location.pathname.startsWith(props?.items[idx].path)) {
                setIndex(idx);
                return;
            }
        }
    }, []);

    const onTabChange = useCallback((_event: SyntheticEvent, idx: number) => {
        setIndex(idx);
    }, []);

    const onClickMenu = useCallback((menu: { path: string }) => {
        navigate(menu.path, {replace: true});
    }, []);

    const onClickHome = useCallback(() => {
        navigate("/", {replace: true});
    }, []);

    return (
        <Box>
            <Box sx={{display: "flex", justifyContent: "space-between", alignItems: "center"}}>
                <Box sx={{
                    display: "flex", justifyContent: "flex-start", alignItems: "center", marginTop: "10px",
                }}>
                    <Box sx={{
                        display: "flex",
                        justifyContent: "flex-start",
                        alignItems: "center",
                        margin: "0px 30px 10px 20px",
                        fontWeight: "bold",
                        fontSize: "20px"
                    }} onClick={onClickHome}>
                        <img className="spin" height={40} src={FavIcon} alt="mesh"/>
                        <span style={{marginLeft: "15px", cursor: "pointer"}}>MESH</span>
                    </Box>
                    <Box>
                        <Box>
                            <Tabs value={index} onChange={onTabChange} aria-label="menu tabs">
                                {
                                    props?.items?.map((menu, idx) => {
                                        return <Tab key={menu.name} label={menu.name} value={idx}
                                                    sx={{marginBottom: "10px"}} onClick={() => onClickMenu(menu)}/>
                                    })
                                }
                            </Tabs>
                        </Box>
                    </Box>
                </Box>
                <Box sx={{marginRight: "50px"}}>
                    <Stack direction="row" spacing={1}>
                        <Tooltip title="" enterDelay={300}>
                            <IconButton
                                component="a"
                                color="primary"
                                href="https://github.com/mesh/mesh"
                                data-ga-event-category="header"
                                data-ga-event-action="github"
                            >
                                <GitHubIcon fontSize="small"/>
                            </IconButton>
                        </Tooltip>
                    </Stack>
                </Box>
            </Box>
            <Divider/>
        </Box>
    );
}

export default Header;
