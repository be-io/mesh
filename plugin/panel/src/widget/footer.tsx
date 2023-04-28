/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box} from '@mui/material';
import {useNavigate} from 'react-router-dom';

function Footer(props: any) {

    const navigate = useNavigate();
    const clickMenu = (menu: { name?: string; path: any; }) => {
        navigate(menu.path, {replace: true});
    };
    return (
        <Box sx={{
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            fontSize: 12,
            margin: '10px auto 10px auto',
            position: 'fixed',
            bottom: 0,
            left: 0,
            width: '100%',
            overflowY: 'scroll'
        }}>
            Copyright © 2023 1285.tech Inc. 保留所有权利。 隐私政策 使用条款 销售政策 法律信息 网站地图
        </Box>
    );
}

export default Footer;
