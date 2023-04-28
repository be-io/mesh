/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Box, Button, Container, Grid, Link, Stack, Theme, Typography, useMediaQuery} from '@mui/material';
import {useTheme} from '@mui/material/styles';
import {FacebookIcon, FavIcon, GoogleIcon, TwitterIcon} from '@/widget/icons';

export function SocialOAuth2() {

    const theme = useTheme();
    const matchDownSM = useMediaQuery(theme.breakpoints.down('sm'));

    const onGoogleAuth = async () => {
        // login || singup
    };

    const onTwitterAuth = async () => {
        // login || singup
    };

    const onFacebookAuth = async () => {
        // login || singup
    };

    return (
        <Stack
            direction="row"
            spacing={matchDownSM ? 1 : 2}
            justifyContent={matchDownSM ? 'space-around' : 'space-between'}
            sx={{'& .MuiButton-startIcon': {mr: matchDownSM ? 0 : 1, ml: matchDownSM ? 0 : -0.5}}}
        >
            <Button
                variant="outlined"
                color="secondary"
                fullWidth={!matchDownSM}
                startIcon={<img src={GoogleIcon} alt="Google"/>}
                onClick={onGoogleAuth}
            >
                {!matchDownSM && 'Google'}
            </Button>
            <Button
                variant="outlined"
                color="secondary"
                fullWidth={!matchDownSM}
                startIcon={<img src={TwitterIcon} alt="Twitter"/>}
                onClick={onTwitterAuth}
            >
                {!matchDownSM && 'Twitter'}
            </Button>
            <Button
                variant="outlined"
                color="secondary"
                fullWidth={!matchDownSM}
                startIcon={<img src={FacebookIcon} alt="Facebook"/>}
                onClick={onFacebookAuth}
            >
                {!matchDownSM && 'Facebook'}
            </Button>
        </Stack>
    );
};

export function SignButtom() {
    const matchDownSM = useMediaQuery((theme: Theme) => theme.breakpoints.down('sm'));

    return (
        <Container maxWidth="xl">
            <Stack
                direction={matchDownSM ? 'column' : 'row'}
                justifyContent={matchDownSM ? 'center' : 'space-between'}
                spacing={2}
                textAlign={matchDownSM ? 'center' : 'inherit'}
            >
                <Typography variant="subtitle2" color="secondary" component="span">
                    &copy; Mantis React Dashboard Template By&nbsp;
                    <Typography component={Link} variant="subtitle2" href="https://codedthemes.com" target="_blank"
                                underline="hover">
                        CodedThemes
                    </Typography>
                </Typography>

                <Stack direction={matchDownSM ? 'column' : 'row'} spacing={matchDownSM ? 1 : 3}
                       textAlign={matchDownSM ? 'center' : 'inherit'}>
                    <Typography
                        color="secondary"
                        component={Link}
                        href="https://material-ui.com/store/contributors/codedthemes/"
                        target="_blank"
                        underline="hover"
                    >
                        MUI Templates
                    </Typography>
                    <Typography
                        color="secondary"
                        component={Link}
                        href="https://codedthemes.com"
                        target="_blank"
                        underline="hover"
                    >
                        Privacy Policy
                    </Typography>
                    <Typography
                        color="secondary"
                        component={Link}
                        href="https://codedthemes.support-hub.io/"
                        target="_blank"
                        underline="hover"
                    >
                        Support
                    </Typography>
                </Stack>
            </Stack>
        </Container>
    );
};


export function SignCard(option: { children: any, other: any[] }) {
    return (
        // <MainCard
        //     sx={{
        //         maxWidth: {xs: 400, lg: 475},
        //         margin: {xs: 2.5, md: 3},
        //         '& > *': {
        //             flexGrow: 1,
        //             flexBasis: '50%'
        //         }
        //     }}
        //     content={false}
        //     {...option.other}
        //     border={false}
        //     boxShadow
        // >
        //     <Box sx={{p: {xs: 2, sm: 3, md: 4, xl: 5}}}>{option.children}</Box>
        // </MainCard>
        <Box></Box>
    )
}

export function SignBox(option: { children: any }) {
    return (
        <Box sx={{minHeight: '100vh'}}>
            <Box>
                <img height={40} src={FavIcon} alt="mesh"/>
                <span style={{marginLeft: "15px", cursor: "pointer"}}>MESH</span>
            </Box>
            <Grid
                container
                direction="column"
                justifyContent="flex-end"
                sx={{
                    minHeight: '100vh'
                }}
            >
                <Grid item xs={12} sx={{ml: 3, mt: 3}}>
                    <Box>
                        <img height={40} src={FavIcon} alt="mesh"/>
                        <span style={{marginLeft: "15px", cursor: "pointer"}}>MESH</span>
                    </Box>
                </Grid>
                <Grid item xs={12}>
                    <Grid
                        item
                        xs={12}
                        container
                        justifyContent="center"
                        alignItems="center"
                        sx={{minHeight: {xs: 'calc(100vh - 134px)', md: 'calc(100vh - 112px)'}}}
                    >
                        <Grid item>
                            <SignCard children={option.children} other={[]}/>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item xs={12} sx={{m: 3, mt: 1}}>
                    <SignButtom/>
                </Grid>
            </Grid>
        </Box>
    )
}
