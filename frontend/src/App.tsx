import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import Login from './pages/auth/Login';
import Register from './pages/auth/Register';

import Layout from './components/layout/Layout';
import ManagerDashboard from './pages/manager/Dashboard';
import Employees from './pages/manager/Employees';
import Warehouses from './pages/manager/Warehouses';
import WarehouseDetails from './pages/manager/WarehouseDetails';
import AuditLogs from './pages/manager/AuditLogs';
import SupervisorDashboard from './pages/supervisor/Dashboard';
import SupervisorEmployees from './pages/supervisor/Employees';
import StaffDashboard from './pages/staff/Dashboard';
import AuditorDashboard from './pages/auditor/Dashboard';
import AuditorWarehouses from './pages/auditor/Warehouses';
import AuditorWarehouseDetails from './pages/auditor/WarehouseDetails';
import AuditorAuditLogs from './pages/auditor/AuditLogs';

const ProtectedRoute: React.FC<{ children: React.ReactNode; roles?: string[] }> = ({ children, roles }) => {
    const { user, isAuthenticated, isLoading } = useAuth();

    if (isLoading) {
        return <div className="flex items-center justify-center min-h-screen">Loading...</div>;
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" />;
    }

    if (roles && user && !roles.includes(user.role)) {
        return <Navigate to="/" />; // Or unauthorized page
    }

    return <>{children}</>;
};

function App() {
    return (
        <AuthProvider>
            <Router>
                <Routes>
                    <Route path="/login" element={<Login />} />
                    <Route path="/signup" element={<Register />} />

                    <Route element={<Layout />}>
                        {/* Manager Routes */}
                        <Route path="/manager/dashboard" element={<ProtectedRoute roles={['Manager']}><ManagerDashboard /></ProtectedRoute>} />
                        <Route path="/manager/employees" element={<ProtectedRoute roles={['Manager']}><Employees /></ProtectedRoute>} />
                        <Route path="/manager/warehouses" element={<ProtectedRoute roles={['Manager']}><Warehouses /></ProtectedRoute>} />
                        <Route path="/manager/logs" element={<ProtectedRoute roles={['Manager']}><AuditLogs /></ProtectedRoute>} />
                        <Route path="/manager/warehouse/:id" element={<ProtectedRoute roles={['Manager']}><WarehouseDetails /></ProtectedRoute>} />

                        {/* Supervisor Routes */}
                        <Route path="/supervisor/dashboard" element={<ProtectedRoute roles={['Supervisor']}><SupervisorDashboard /></ProtectedRoute>} />
                        <Route path="/supervisor/employees" element={<ProtectedRoute roles={['Supervisor']}><SupervisorEmployees /></ProtectedRoute>} />

                        {/* Staff Routes */}
                        <Route path="/staff/dashboard" element={<ProtectedRoute roles={['Staff']}><StaffDashboard /></ProtectedRoute>} />

                        {/* Auditor Routes */}
                        <Route path="/auditor/dashboard" element={<ProtectedRoute roles={['Auditor']}><AuditorDashboard /></ProtectedRoute>} />
                        <Route path="/auditor/warehouses" element={<ProtectedRoute roles={['Auditor']}><AuditorWarehouses /></ProtectedRoute>} />
                        <Route path="/auditor/warehouse/:id" element={<ProtectedRoute roles={['Auditor']}><AuditorWarehouseDetails /></ProtectedRoute>} />
                        <Route path="/auditor/logs" element={<ProtectedRoute roles={['Auditor']}><AuditorAuditLogs /></ProtectedRoute>} />
                    </Route>

                    <Route path="/" element={<Navigate to="/login" />} />
                </Routes>
                <ToastContainer position="top-right" autoClose={3000} />
            </Router>
        </AuthProvider>
    );
}

export default App;
