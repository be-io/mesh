/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box, Button, Tab} from "@mui/material";
import Menu, {ItemOption} from "@/widget/menu";
import {useEffect, useState} from "react";
import CloudCircleIcon from '@mui/icons-material/CloudCircle';
import service from "@/services/service";
import {context, Entity} from "@mesh/mesh";
import Iframe from "react-iframe";
import {TabContext, TabList, TabPanel} from "@mui/lab";

export default function Fops(props: {}) {

    const [tab, setTab] = useState('0');
    const [nodes, setNodes] = useState<ItemOption[]>([]);
    const [node, setNode] = useState<ItemOption>(new ItemOption());
    const [remainsSeconds, setRemainsSeconds] = useState(0);

    const tabs = [
        {name: '隐私计算平台', path: 'fops/gaia/'},
        {name: 'PaaS平台', path: 'fops/paas'},
        {name: '监控平台', path: 'fops/grafana'}
    ];

    useEffect(() => {
        let rs = remainsSeconds
        const timer = setInterval(() => {
            rs = Math.max(rs - 1, 0);
            setRemainsSeconds(rs);
        }, 1000);
        setTimeout(() => {
            clearInterval(timer);
        }, remainsSeconds * 1000);
        service.network.getRoutes(context()).then(routes => {
            const rs: ItemOption[] = routes.map(r => {
                return {
                    key: r.node_id,
                    name: r.name,
                    icon: <CloudCircleIcon color="primary"/>,
                    onclick: (o: ItemOption) => {
                        setNode(o);
                    }
                }
            });
            setNodes(rs);
            setNode(rs[0]);
        })
    }, []);

    const onTabChange = (e: any, v: any) => {
        setTab(v);
        service.kv.get(context(), `mesh.plugin.proxy.fops.auth.${node.key}`).then(data => {
            const x = data.readObject();
            console.log(x);
        })
    }

    const onFopsAuthorizeClick = () => {
        service.kv.put(context(), `mesh.plugin.proxy.fops.auth.${node.key}`, Entity.wrap(10 * 60 * 1000)).then(data => {

        });
    }

    const onFopsUnAuthorizeClick = () => {
        service.kv.remove(context(), `mesh.plugin.proxy.fops.auth.${node.key}`).then(data => {

        });
    }

    const AuthorizeView = () => {
        const days = Math.floor(remainsSeconds / 60 / 60 / 24);
        const hours = Math.floor(remainsSeconds / 60 / 60);
        const minutes = Math.floor(remainsSeconds / 60);
        const seconds = Math.floor(remainsSeconds % 60);
        if (remainsSeconds > 0) {
            return (
                <Box>
                    <Button variant="outlined" color="warning" onClick={onFopsUnAuthorizeClick}>
                        解除授权 | 授权剩余时间{days}天{hours}时{minutes}分{seconds}秒
                    </Button>

                </Box>
            )
        }
        return (
            <Box>
                <Button variant="outlined" color="primary" onClick={onFopsAuthorizeClick}>
                    授权{node.name}联邦运维
                </Button>
            </Box>
        )
    }

    return (
        <Box sx={{
            display: 'flex',
            justifyContent: 'flex-start',
            alignItems: 'flex-start',
            backgroundColor: 'rgb(244,246,249)',
            width: '100%'
        }}>
            <Menu menus={nodes}/>
            <Box sx={{margin: '10px auto auto 20px', width: '100%'}}>
                <Box sx={{margin: 'auto auto auto 18px'}}>
                    <AuthorizeView/>
                </Box>
                <Box>
                    <TabContext value={tab}>
                        <Box sx={{borderBottom: 1, borderColor: 'divider'}}>
                            <TabList onChange={onTabChange}>
                                {
                                    tabs.map((attr, idx) => {
                                        return <Tab key={`${node.key}-${attr.name}-${idx}`} label={attr.name}
                                                    value={`${idx}`}/>
                                    })
                                }
                            </TabList>
                        </Box>
                        {
                            tabs.map((attr, idx) => {
                                return (
                                    <TabPanel key={`${node.key}-${attr.name}-${idx}`} value={`${idx}`}>
                                        <Iframe url={`mesh/${node.key}/${attr.path}`}
                                                width="100%"
                                                height="100%"
                                                id=""
                                                className=""
                                                display="block"
                                                position="relative"/>
                                    </TabPanel>
                                )
                            })
                        }
                    </TabContext>
                </Box>

            </Box>
        </Box>
    )
}