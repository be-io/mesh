/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import Button from '@mui/material/Button';
import {
    Box,
    DialogActions,
    DialogContent,
    DialogTitle,
    FormControl,
    InputLabel,
    Modal,
    OutlinedInput,
} from '@mui/material';
import {styled} from '@mui/material/styles';
import Divider from '@mui/material/Divider';
import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Paper from '@mui/material/Paper';
import {ElementType, useState} from 'react';
import {VscLaw, VscWorkspaceTrusted, VscWorkspaceUntrusted} from "react-icons/vsc";


function ProjectWizardDialog(props: { enable: boolean, onClose: any, onSubmit: any; }) {

    const data = [
        {icon: <VscLaw/>, label: 'x'},
        {icon: <VscLaw/>, label: 'v'},
    ];
    const {enable, onClose, onSubmit} = props;
    const [projectKind, setProjectKind] = useState('DP');
    const [projectName, setProjectName] = useState('');

    const FireNav = styled(List)<{ component?: ElementType }>({
        '& .MuiListItemButton-root': {
            paddingLeft: 24,
            paddingRight: 24,
        },
        '& .MuiListItemIcon-root': {
            minWidth: 0,
            marginRight: 16,
        },
        '& .MuiSvgIcon-root': {
            fontSize: 20,
        },
    });

    return (
        <div>
            <Modal
                open={enable}
                hideBackdrop={true}
                onClose={onClose}
                aria-describedby='alert-dialog-slide-description'
                sx={{
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                }}
            >
                <Box sx={{height: 520, width: 680, backgroundColor: 'var(--main-modal-color)'}}>
                    <DialogTitle>新建工程</DialogTitle>
                    <DialogContent sx={{display: 'flex', alignItems: 'flex-start'}}>
                        <Paper elevation={0} sx={{width: 240}}>
                            <FireNav component='nav' disablePadding>
                                <ListItemButton component='a' href='#customized-list'>
                                    <ListItemIcon sx={{fontSize: 20}}>🔥</ListItemIcon>
                                    <ListItemText
                                        sx={{my: 0}}
                                        primary="x"
                                        primaryTypographyProps={{
                                            fontSize: 20,
                                            fontWeight: 'medium',
                                            letterSpacing: 0,
                                        }}
                                    />
                                </ListItemButton>
                                <Divider/>
                                <Box sx={{bgcolor: 'rgba(71, 98, 130, 0.2)', pb: 30}}>
                                    {
                                        data.map((item) => (
                                            <ListItemButton
                                                selected={item.label == projectKind}
                                                key={item.label}
                                                sx={{py: 0, minHeight: 32, color: 'rgba(255,255,255,.8)'}}
                                                onClick={() => setProjectKind(item.label)}
                                            >
                                                <ListItemIcon sx={{color: 'inherit'}}>
                                                    {item.icon}
                                                </ListItemIcon>
                                                <ListItemText
                                                    primary={item.label}
                                                    primaryTypographyProps={{fontSize: 14, fontWeight: 'medium'}}
                                                />
                                            </ListItemButton>
                                        ))}
                                </Box>
                            </FireNav>
                        </Paper>
                        <Box>
                            <FormControl fullWidth sx={{m: 1}}>
                                <InputLabel
                                    htmlFor='outlined-adornment-project-name'
                                    size='small'>名称</InputLabel>
                                <OutlinedInput
                                    id='outlined-adornment-project-name'
                                    value={projectName}
                                    onChange={(e) => setProjectName(e.target.value)}
                                    label="名称"
                                    size='small'
                                    sx={{width: 360}}
                                />
                            </FormControl>
                        </Box>
                    </DialogContent>
                    <DialogActions sx={{mr: 5}}>
                        <Button onClick={() => onSubmit({kind: projectKind, name: projectName})} variant='contained'
                                size='small'
                                startIcon={<VscWorkspaceTrusted/>}>
                            确认
                        </Button>
                        <Button onClick={onClose} variant='contained' size='small' startIcon={<VscWorkspaceUntrusted/>}>
                            取消
                        </Button>
                    </DialogActions>
                </Box>
            </Modal>
        </div>
    );
}

export default ProjectWizardDialog;
