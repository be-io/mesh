/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import React from 'react';
import {
    Button,
    Checkbox,
    Divider,
    FormControlLabel,
    FormHelperText,
    Grid,
    IconButton,
    InputAdornment,
    InputLabel,
    Link,
    OutlinedInput,
    Stack,
    Typography
} from '@mui/material';
import {Link as RouterLink} from 'react-router-dom';
import {AiOutlineEye, AiOutlineEyeInvisible} from "react-icons/ai";
import * as Yup from 'yup';
import {Formik} from 'formik';
import {SignBox, SocialOAuth2} from '@/views/home/signbox';

export function AuthLogin() {

    const [checked, setChecked] = React.useState(false);
    const [showPassword, setShowPassword] = React.useState(false);
    const onClickShowPassword = (event: any) => {
        setShowPassword(!showPassword);
    };

    const onMouseDownPassword = (event: any) => {
        event.preventDefault();
    };

    return (
        <Formik
            initialValues={{
                email: 'info@codedthemes.com',
                password: '123456',
                submit: null
            }}
            validationSchema={Yup.object().shape({
                email: Yup.string().email('Must be a valid email').max(255).required('Email is required'),
                password: Yup.string().max(255).required('Password is required')
            })}
            onSubmit={async (values, {setErrors, setStatus, setSubmitting}) => {
                try {
                    setStatus({success: false});
                    setSubmitting(false);
                } catch (e) {
                    setStatus({success: false});
                    setErrors({submit: `${e}`});
                    setSubmitting(false);
                }
            }}
        >
            {({errors, handleBlur, handleChange, handleSubmit, isSubmitting, touched, values}) => (
                <form noValidate onSubmit={handleSubmit}>
                    <Grid container spacing={3}>
                        <Grid item xs={12}>
                            <Stack spacing={1}>
                                <InputLabel htmlFor="email-login">Email Address</InputLabel>
                                <OutlinedInput
                                    id="email-login"
                                    type="email"
                                    value={values.email}
                                    name="email"
                                    onBlur={handleBlur}
                                    onChange={handleChange}
                                    placeholder="Enter email address"
                                    fullWidth
                                    error={Boolean(touched.email && errors.email)}
                                />
                                {touched.email && errors.email && (
                                    <FormHelperText error id="standard-weight-helper-text-email-login">
                                        {errors.email}
                                    </FormHelperText>
                                )}
                            </Stack>
                        </Grid>
                        <Grid item xs={12}>
                            <Stack spacing={1}>
                                <InputLabel htmlFor="password-login">Password</InputLabel>
                                <OutlinedInput
                                    fullWidth
                                    error={Boolean(touched.password && errors.password)}
                                    id="-password-login"
                                    type={showPassword ? 'text' : 'password'}
                                    value={values.password}
                                    name="password"
                                    onBlur={handleBlur}
                                    onChange={handleChange}
                                    endAdornment={
                                        <InputAdornment position="end">
                                            <IconButton
                                                component="a"
                                                aria-label="toggle password visibility"
                                                onClick={onClickShowPassword}
                                                onMouseDown={onMouseDownPassword}
                                                edge="end"
                                                size="large"
                                            >
                                                {showPassword ? <AiOutlineEye/> : <AiOutlineEyeInvisible/>}
                                            </IconButton>
                                        </InputAdornment>
                                    }
                                    placeholder="Enter password"
                                />
                                {touched.password && errors.password && (
                                    <FormHelperText error id="standard-weight-helper-text-password-login">
                                        {errors.password}
                                    </FormHelperText>
                                )}
                            </Stack>
                        </Grid>

                        <Grid item xs={12} sx={{mt: -1}}>
                            <Stack direction="row" justifyContent="space-between" alignItems="center" spacing={2}>
                                <FormControlLabel
                                    control={
                                        <Checkbox
                                            checked={checked}
                                            onChange={(event) => setChecked(event.target.checked)}
                                            name="checked"
                                            color="primary"
                                            size="small"
                                        />
                                    }
                                    label={<Typography variant="h6">Keep me sign in</Typography>}
                                />
                                <Link variant="h6" component={RouterLink} to="" color="text.primary">
                                    Forgot Password?
                                </Link>
                            </Stack>
                        </Grid>
                        {errors.submit && (
                            <Grid item xs={12}>
                                <FormHelperText error>{errors.submit}</FormHelperText>
                            </Grid>
                        )}
                        <Grid item xs={12}>
                            <Button disableElevation disabled={isSubmitting} fullWidth size="large"
                                    type="submit" variant="contained" color="primary">
                                Login
                            </Button>
                        </Grid>
                        <Grid item xs={12}>
                            <Divider>
                                <Typography variant="caption"> Login with</Typography>
                            </Divider>
                        </Grid>
                        <Grid item xs={12}>
                            <SocialOAuth2/>
                        </Grid>
                    </Grid>
                </form>
            )}
        </Formik>
    );
};

export function Signin() {
    return (
        <SignBox>
            <Grid container spacing={3}>
                <Grid item xs={12}>
                    <Stack direction="row" justifyContent="space-between" alignItems="baseline"
                           sx={{mb: {xs: -0.5, sm: 0.5}}}>
                        <Typography variant="h3">Login</Typography>
                        <Typography sx={{textDecoration: 'none'}} color="primary">
                            <RouterLink to="/mesh/auth/siginup">Don&apos;t have an account?</RouterLink>
                        </Typography>
                    </Stack>
                </Grid>
                <Grid item xs={12}>
                    <AuthLogin/>
                </Grid>
            </Grid>
        </SignBox>
    )
}
