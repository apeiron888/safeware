import React, { Fragment, useState } from 'react';
import { Dialog, Transition } from '@headlessui/react';
import { Link, useLocation, Outlet } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import {
    HiHome,
    HiUsers,
    HiOfficeBuilding,
    HiClipboardList,
    HiMenu,
    HiX,
    HiLogout,
    HiMoon,
    HiSun
} from 'react-icons/hi';
import clsx from 'clsx';

const Layout: React.FC = () => {
    const { user, logout } = useAuth();
    const location = useLocation();
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [darkMode, setDarkMode] = useState(false);

    const toggleDarkMode = () => {
        setDarkMode(!darkMode);
        if (!darkMode) {
            document.documentElement.classList.add('dark');
        } else {
            document.documentElement.classList.remove('dark');
        }
    };

    const navigation = [
        { name: 'Dashboard', href: `/${user?.role?.toLowerCase()}/dashboard`, icon: HiHome, roles: ['Manager', 'Supervisor', 'Staff', 'Auditor'] },
        { name: 'Employees', href: '/manager/employees', icon: HiUsers, roles: ['Manager'] },
        { name: 'Employees', href: '/supervisor/employees', icon: HiUsers, roles: ['Supervisor'] },
        { name: 'Warehouses', href: '/manager/warehouses', icon: HiOfficeBuilding, roles: ['Manager'] },
        { name: 'Audit Logs', href: '/manager/logs', icon: HiClipboardList, roles: ['Manager'] },
        { name: 'Warehouses', href: '/auditor/warehouses', icon: HiOfficeBuilding, roles: ['Auditor'] },
        { name: 'Audit Logs', href: '/auditor/logs', icon: HiClipboardList, roles: ['Auditor'] },
    ];

    const filteredNavigation = navigation.filter(item => user?.role && item.roles.includes(user.role));

    return (
        <div className="h-screen flex overflow-hidden bg-gray-100 dark:bg-dark-bg transition-colors duration-200">
            <Transition.Root show={sidebarOpen} as={Fragment}>
                <Dialog as="div" className="fixed inset-0 flex z-40 md:hidden" onClose={setSidebarOpen}>
                    <Transition.Child
                        as={Fragment}
                        enter="transition-opacity ease-linear duration-300"
                        enterFrom="opacity-0"
                        enterTo="opacity-100"
                        leave="transition-opacity ease-linear duration-300"
                        leaveFrom="opacity-100"
                        leaveTo="opacity-0"
                    >
                        <Dialog.Overlay className="fixed inset-0 bg-gray-600 bg-opacity-75" />
                    </Transition.Child>
                    <Transition.Child
                        as={Fragment}
                        enter="transition ease-in-out duration-300 transform"
                        enterFrom="-translate-x-full"
                        enterTo="translate-x-0"
                        leave="transition ease-in-out duration-300 transform"
                        leaveFrom="translate-x-0"
                        leaveTo="-translate-x-full"
                    >
                        <div className="relative flex-1 flex flex-col max-w-xs w-full bg-white dark:bg-dark-surface">
                            <Transition.Child
                                as={Fragment}
                                enter="ease-in-out duration-300"
                                enterFrom="opacity-0"
                                enterTo="opacity-100"
                                leave="ease-in-out duration-300"
                                leaveFrom="opacity-100"
                                leaveTo="opacity-0"
                            >
                                <div className="absolute top-0 right-0 -mr-12 pt-2">
                                    <button
                                        type="button"
                                        className="ml-1 flex items-center justify-center h-10 w-10 rounded-full focus:outline-none focus:ring-2 focus:ring-inset focus:ring-white"
                                        onClick={() => setSidebarOpen(false)}
                                    >
                                        <span className="sr-only">Close sidebar</span>
                                        <HiX className="h-6 w-6 text-white" aria-hidden="true" />
                                    </button>
                                </div>
                            </Transition.Child>
                            <div className="flex-1 h-0 pt-5 pb-4 overflow-y-auto">
                                <div className="flex-shrink-0 flex items-center px-4">
                                    <h1 className="text-2xl font-bold text-primary-600">Vaultory</h1>
                                </div>
                                <nav className="mt-5 px-2 space-y-1">
                                    {filteredNavigation.map((item) => (
                                        <Link
                                            key={item.name}
                                            to={item.href}
                                            className={clsx(
                                                location.pathname.startsWith(item.href)
                                                    ? 'bg-primary-50 dark:bg-primary-900 text-primary-600 dark:text-primary-200'
                                                    : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700',
                                                'group flex items-center px-2 py-2 text-base font-medium rounded-md'
                                            )}
                                        >
                                            <item.icon
                                                className={clsx(
                                                    location.pathname.startsWith(item.href) ? 'text-primary-600 dark:text-primary-200' : 'text-gray-400 dark:text-gray-400 group-hover:text-gray-500',
                                                    'mr-4 flex-shrink-0 h-6 w-6'
                                                )}
                                                aria-hidden="true"
                                            />
                                            {item.name}
                                        </Link>
                                    ))}
                                </nav>
                            </div>
                        </div>
                    </Transition.Child>
                    <div className="flex-shrink-0 w-14" aria-hidden="true">
                        {/* Force sidebar to shrink to fit close icon */}
                    </div>
                </Dialog>
            </Transition.Root>

            {/* Static sidebar for desktop */}
            <div className="hidden md:flex md:flex-shrink-0">
                <div className="flex flex-col w-64">
                    <div className="flex-1 flex flex-col min-h-0 bg-white dark:bg-dark-surface border-r border-gray-200 dark:border-gray-700">
                        <div className="flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
                            <div className="flex items-center flex-shrink-0 px-4">
                                <h1 className="text-2xl font-bold text-primary-600">Vaultory</h1>
                            </div>
                            <nav className="mt-5 flex-1 px-2 space-y-1">
                                {filteredNavigation.map((item) => (
                                    <Link
                                        key={item.name}
                                        to={item.href}
                                        className={clsx(
                                            location.pathname.startsWith(item.href)
                                                ? 'bg-primary-50 dark:bg-primary-900 text-primary-600 dark:text-primary-200'
                                                : 'text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700',
                                            'group flex items-center px-2 py-2 text-sm font-medium rounded-md'
                                        )}
                                    >
                                        <item.icon
                                            className={clsx(
                                                location.pathname.startsWith(item.href) ? 'text-primary-600 dark:text-primary-200' : 'text-gray-400 dark:text-gray-400 group-hover:text-gray-500',
                                                'mr-3 flex-shrink-0 h-6 w-6'
                                            )}
                                            aria-hidden="true"
                                        />
                                        {item.name}
                                    </Link>
                                ))}
                            </nav>
                        </div>
                        <div className="flex-shrink-0 flex border-t border-gray-200 dark:border-gray-700 p-4">
                            <div className="flex-shrink-0 w-full group block">
                                <div className="flex items-center">
                                    <div className="ml-3">
                                        <p className="text-sm font-medium text-gray-700 dark:text-gray-200">{user?.full_name}</p>
                                        <p className="text-xs font-medium text-gray-500 dark:text-gray-400">{user?.role}</p>
                                    </div>
                                    <button onClick={logout} className="ml-auto text-gray-400 hover:text-gray-500">
                                        <HiLogout className="h-6 w-6" />
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div className="flex flex-col w-0 flex-1 overflow-hidden">
                <div className="md:hidden pl-1 pt-1 sm:pl-3 sm:pt-3">
                    <button
                        type="button"
                        className="-ml-0.5 -mt-0.5 h-12 w-12 inline-flex items-center justify-center rounded-md text-gray-500 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
                        onClick={() => setSidebarOpen(true)}
                    >
                        <span className="sr-only">Open sidebar</span>
                        <HiMenu className="h-6 w-6" aria-hidden="true" />
                    </button>
                </div>

                {/* Top header for dark mode toggle */}
                <div className="flex justify-end p-4 bg-white dark:bg-dark-surface border-b border-gray-200 dark:border-gray-700 md:hidden">
                    <button onClick={toggleDarkMode} className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700">
                        {darkMode ? <HiSun className="h-6 w-6 text-yellow-400" /> : <HiMoon className="h-6 w-6 text-gray-600" />}
                    </button>
                </div>
                <div className="hidden md:flex justify-end p-4 bg-white dark:bg-dark-surface border-b border-gray-200 dark:border-gray-700">
                    <button onClick={toggleDarkMode} className="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700">
                        {darkMode ? <HiSun className="h-6 w-6 text-yellow-400" /> : <HiMoon className="h-6 w-6 text-gray-600" />}
                    </button>
                </div>

                <main className="flex-1 relative z-0 overflow-y-auto focus:outline-none">
                    <div className="py-6">
                        <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
                            <Outlet />
                        </div>
                    </div>
                </main>
            </div>
        </div>
    );
};

export default Layout;
