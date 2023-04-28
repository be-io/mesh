/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Link as RouterLink, useNavigate} from 'react-router-dom';
import {MouseEvent, useEffect, useState} from 'react';
import {
    Box,
    Button,
    Divider,
    FormControl,
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
import {SignBox, SocialOAuth2} from '@/views/home/signbox';
import * as Yup from 'yup';
import {Formik} from 'formik';
import {strengthColor, strengthIndicator} from '@/views/home/strength';
import {AiOutlineEye, AiOutlineEyeInvisible} from "react-icons/ai";

export function AuthRegister() {
    const navigate = useNavigate();

    const [level, setLevel] = useState<{ label: string, color: string }>();
    const [showPassword, setShowPassword] = useState(false);
    const onClickShowPassword = () => {
        setShowPassword(!showPassword);
    };

    const onMouseDownPassword = (event: MouseEvent) => {
        event.preventDefault();
    };

    const onChangePassword = (v: string) => {
        setLevel(strengthColor(strengthIndicator(v)));
    };

    useEffect(() => {
        onChangePassword('');
    }, []);

    return (
        <Formik
            initialValues={{
                firstname: '',
                lastname: '',
                email: '',
                company: '',
                password: '',
                submit: null
            }}
            validationSchema={Yup.object().shape({
                firstname: Yup.string().max(255).required('First Name is required'),
                lastname: Yup.string().max(255).required('Last Name is required'),
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
                        <Grid item xs={12} md={6}>
                            <Stack spacing={1}>
                                <InputLabel htmlFor="firstname-signup">First Name*</InputLabel>
                                <OutlinedInput
                                    id="firstname-login"
                                    type="firstname"
                                    value={values.firstname}
                                    name="firstname"
                                    onBlur={handleBlur}
                                    onChange={handleChange}
                                    placeholder="John"
                                    fullWidth
                                    error={Boolean(touched.firstname && errors.firstname)}
                                />
                                {touched.firstname && errors.firstname && (
                                    <FormHelperText error id="helper-text-firstname-signup">
                                        {errors.firstname}
                                    </FormHelperText>
                                )}
                            </Stack>
                        </Grid>
                        <Grid item xs={12} md={6}>
                            <Stack spacing={1}>
                                <InputLabel htmlFor="lastname-signup">Last Name*</InputLabel>
                                <OutlinedInput
                                    fullWidth
                                    error={Boolean(touched.lastname && errors.lastname)}
                                    id="lastname-signup"
                                    type="lastname"
                                    value={values.lastname}
                                    name="lastname"
                                    onBlur={handleBlur}
                                    onChange={handleChange}
                                    placeholder="Doe"
                                    inputProps={{}}
                                />
                                {touched.lastname && errors.lastname && (
                                    <FormHelperText error id="helper-text-lastname-signup">
                                        {errors.lastname}
                                    </FormHelperText>
                                )}
                            </Stack>
                        </Grid>
                        <Grid item xs={12}>
                            <Stack spacing={1}>
                                <InputLabel htmlFor="company-signup">Company</InputLabel>
                                <OutlinedInput
                                    fullWidth
                                    error={Boolean(touched.company && errors.company)}
                                    id="company-signup"
                                    value={values.company}
                                    name="company"
                                    onBlur={handleBlur}
                                    onChange={handleChange}
                                    placeholder="Demo Inc."
                                    inputProps={{}}
                                />
                                {touched.company && errors.company && (
                                    <FormHelperText error id="helper-text-company-signup">
                                        {errors.company}
                                    </FormHelperText>
                                )}
                            </Stack>
                        </Grid>
                        <Grid item xs={12}>
                            <Stack spacing={1}>
                                <InputLabel htmlFor="email-signup">Email Address*</InputLabel>
                                <OutlinedInput
                                    fullWidth
                                    error={Boolean(touched.email && errors.email)}
                                    id="email-login"
                                    type="email"
                                    value={values.email}
                                    name="email"
                                    onBlur={handleBlur}
                                    onChange={handleChange}
                                    placeholder="demo@company.com"
                                    inputProps={{}}
                                />
                                {touched.email && errors.email && (
                                    <FormHelperText error id="helper-text-email-signup">
                                        {errors.email}
                                    </FormHelperText>
                                )}
                            </Stack>
                        </Grid>
                        <Grid item xs={12}>
                            <Stack spacing={1}>
                                <InputLabel htmlFor="password-signup">Password</InputLabel>
                                <OutlinedInput
                                    fullWidth
                                    error={Boolean(touched.password && errors.password)}
                                    id="password-signup"
                                    type={showPassword ? 'text' : 'password'}
                                    value={values.password}
                                    name="password"
                                    onBlur={handleBlur}
                                    onChange={(e) => {
                                        handleChange(e);
                                        onChangePassword(e.target.value);
                                    }}
                                    endAdornment={
                                        <InputAdornment position="end">
                                            <IconButton
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
                                    placeholder="******"
                                    inputProps={{}}
                                />
                                {touched.password && errors.password && (
                                    <FormHelperText error id="helper-text-password-signup">
                                        {errors.password}
                                    </FormHelperText>
                                )}
                            </Stack>
                            <FormControl fullWidth sx={{mt: 2}}>
                                <Grid container spacing={2} alignItems="center">
                                    <Grid item>
                                        <Box sx={{
                                            bgcolor: level?.color,
                                            width: 85,
                                            height: 8,
                                            borderRadius: '7px'
                                        }}/>
                                    </Grid>
                                    <Grid item>
                                        <Typography fontSize="0.75rem">
                                            {level?.label}
                                        </Typography>
                                    </Grid>
                                </Grid>
                            </FormControl>
                        </Grid>
                        <Grid item xs={12}>
                            <Typography variant="body2">
                                By Signing up, you agree to our &nbsp;
                                <Link component={RouterLink} to="#">
                                    Terms of Service
                                </Link>
                                &nbsp; and &nbsp;
                                <Link component={RouterLink} to="#">
                                    Privacy Policy
                                </Link>
                            </Typography>
                        </Grid>
                        {errors.submit && (
                            <Grid item xs={12}>
                                <FormHelperText error>{errors.submit}</FormHelperText>
                            </Grid>
                        )}
                        <Grid item xs={12}>
                            <Button disableElevation disabled={isSubmitting} fullWidth size="large"
                                    type="submit" variant="contained" color="primary">
                                Create Account
                            </Button>
                        </Grid>
                        <Grid item xs={12}>
                            <Divider>
                                <Typography variant="caption">Sign up with</Typography>
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
}

export function Signup() {

    return (
        <SignBox>
            <Grid container spacing={3}>
                <Grid item xs={12}>
                    <Stack direction="row" justifyContent="space-between" alignItems="baseline"
                           sx={{mb: {xs: -0.5, sm: 0.5}}}>
                        <Typography variant="h3">Sign up</Typography>
                        <Typography variant="h3" sx={{textDecoration: 'none'}} color="primary">
                            <RouterLink to="/mesh/auth/signin">Already have an account?</RouterLink>
                        </Typography>
                    </Stack>
                </Grid>
                <Grid item xs={12}>
                    <AuthRegister/>
                </Grid>
            </Grid>
        </SignBox>
    )
}