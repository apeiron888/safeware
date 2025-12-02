import React, { useState } from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { useNavigate, Link } from 'react-router-dom';
import api from '../../services/api';
import { toast } from 'react-toastify';

const Register: React.FC = () => {
    const navigate = useNavigate();
    const [isLoading, setIsLoading] = useState(false);

    const formik = useFormik({
        initialValues: {
            company_name: '',
            full_name: '',
            email: '',
            password: '',
        },
        validationSchema: Yup.object({
            company_name: Yup.string().required('Required'),
            full_name: Yup.string().required('Required'),
            email: Yup.string().email('Invalid email address').required('Required'),
            password: Yup.string().min(8, 'Must be at least 8 characters').required('Required'),
        }),
        onSubmit: async (values) => {
            setIsLoading(true);
            try {
                await api.post('/auth/register', values);
                toast.success('Registration successful! Please login.');
                navigate('/login');
            } catch (error: any) {
                toast.error(error.response?.data?.error || 'Registration failed');
            } finally {
                setIsLoading(false);
            }
        },
    });

    return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-dark-bg py-12 px-4 sm:px-6 lg:px-8 transition-colors duration-200">
            <div className="max-w-md w-full space-y-8 bg-white dark:bg-dark-surface p-8 rounded-xl shadow-lg border border-gray-100 dark:border-gray-700">
                <div>
                    <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900 dark:text-white">
                        Register Company
                    </h2>
                    <p className="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
                        Create a new SafeWare account for your organization
                    </p>
                </div>
                <form className="mt-8 space-y-6" onSubmit={formik.handleSubmit}>
                    <div className="rounded-md shadow-sm space-y-4">
                        <div>
                            <label htmlFor="company_name" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Company Name</label>
                            <input
                                id="company_name"
                                type="text"
                                className={`mt-1 appearance-none rounded-lg relative block w-full px-3 py-2 border ${formik.touched.company_name && formik.errors.company_name ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-800 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm`}
                                {...formik.getFieldProps('company_name')}
                            />
                            {formik.touched.company_name && formik.errors.company_name && (
                                <div className="text-red-500 text-xs mt-1">{formik.errors.company_name}</div>
                            )}
                        </div>

                        <div>
                            <label htmlFor="full_name" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Manager Name</label>
                            <input
                                id="full_name"
                                type="text"
                                className={`mt-1 appearance-none rounded-lg relative block w-full px-3 py-2 border ${formik.touched.full_name && formik.errors.full_name ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-800 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm`}
                                {...formik.getFieldProps('full_name')}
                            />
                            {formik.touched.full_name && formik.errors.full_name && (
                                <div className="text-red-500 text-xs mt-1">{formik.errors.full_name}</div>
                            )}
                        </div>

                        <div>
                            <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Email Address</label>
                            <input
                                id="email"
                                type="email"
                                className={`mt-1 appearance-none rounded-lg relative block w-full px-3 py-2 border ${formik.touched.email && formik.errors.email ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-800 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm`}
                                {...formik.getFieldProps('email')}
                            />
                            {formik.touched.email && formik.errors.email && (
                                <div className="text-red-500 text-xs mt-1">{formik.errors.email}</div>
                            )}
                        </div>

                        <div>
                            <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
                            <input
                                id="password"
                                type="password"
                                className={`mt-1 appearance-none rounded-lg relative block w-full px-3 py-2 border ${formik.touched.password && formik.errors.password ? 'border-red-500' : 'border-gray-300 dark:border-gray-600'} placeholder-gray-500 text-gray-900 dark:text-white dark:bg-gray-800 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm`}
                                {...formik.getFieldProps('password')}
                            />
                            {formik.touched.password && formik.errors.password && (
                                <div className="text-red-500 text-xs mt-1">{formik.errors.password}</div>
                            )}
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            disabled={isLoading}
                            className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 transition-colors duration-200"
                        >
                            {isLoading ? 'Registering...' : 'Register'}
                        </button>
                    </div>

                    <div className="text-center mt-4">
                        <p className="text-sm text-gray-600 dark:text-gray-400">
                            Already have an account? <Link to="/login" className="font-medium text-primary-600 hover:text-primary-500">Sign in</Link>
                        </p>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Register;
