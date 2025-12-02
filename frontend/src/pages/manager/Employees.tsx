import React, { useEffect, useState } from 'react';
import api from '../../services/api';
import { HiUser, HiUserGroup, HiIdentification, HiPencil, HiTrash } from 'react-icons/hi';
import { Dialog } from '@headlessui/react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { toast } from 'react-toastify';
import ConfirmDialog from '../../components/common/ConfirmDialog';

interface Employee {
    id: string;
    full_name: string;
    email: string;
    role: string;
    warehouse_id?: string;
    warehouse_name?: string;
}

interface Warehouse {
    id: string;
    name: string;
    location?: string;
}

const Employees: React.FC = () => {
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [selectedRole, setSelectedRole] = useState<'Supervisor' | 'Staff' | 'Auditor'>('Staff');
    const [editingEmployee, setEditingEmployee] = useState<Employee | null>(null);
    const [deletingEmployee, setDeletingEmployee] = useState<Employee | null>(null);
    const [isDeleting, setIsDeleting] = useState(false);

    const fetchEmployees = async () => {
        try {
            const res = await api.get('/manager/employees');
            setEmployees(Array.isArray(res.data) ? res.data : []);
        } catch (error) {
            console.error("Failed to fetch employees", error);
            setEmployees([]);
        }
    };

    const fetchWarehouses = async () => {
        try {
            const res = await api.get('/manager/summary/warehouses');
            const warehousesData = res.data.warehouses || res.data;
            setWarehouses(Array.isArray(warehousesData) ? warehousesData : []);
        } catch (error) {
            console.error("Failed to fetch warehouses", error);
            setWarehouses([]);
        }
    };

    useEffect(() => {
        fetchEmployees();
        fetchWarehouses();
    }, []);

    const formik = useFormik({
        initialValues: {
            full_name: '',
            email: '',
            password: '',
            warehouse_id: '',
        },
        validationSchema: Yup.object({
            full_name: Yup.string().required('Required'),
            email: Yup.string().email('Invalid email').required('Required'),
            password: Yup.string().test('password-required', 'Min 8 chars', function (value) {
                if (editingEmployee) return true; // Not required when editing
                if (!value) return this.createError({ message: 'Required' });
                return value.length >= 8;
            }),
            warehouse_id: Yup.string().test('warehouse-required', 'Required for this role', function (value) {
                if (selectedRole === 'Auditor') return true; // Not required for Auditor
                return !!value;
            }),
        }),
        onSubmit: async (values) => {
            try {
                const role = editingEmployee ? editingEmployee.role : selectedRole;
                const endpoint = editingEmployee
                    ? `/manager/${role.toLowerCase()}/update/${editingEmployee.id}`
                    : `/manager/${role.toLowerCase()}/create`;

                const payload = editingEmployee && !values.password
                    ? { full_name: values.full_name, email: values.email, warehouse_id: values.warehouse_id }
                    : values;

                const response = await api.post(endpoint, payload);

                // Get warehouse name for display
                const warehouseName = values.warehouse_id
                    ? warehouses.find(w => w.id === values.warehouse_id)?.name
                    : undefined;

                if (editingEmployee) {
                    // Update existing employee in state
                    setEmployees(prev => prev.map(emp =>
                        emp.id === editingEmployee.id
                            ? { ...emp, full_name: values.full_name, email: values.email, warehouse_id: values.warehouse_id, warehouse_name: warehouseName }
                            : emp
                    ));
                } else {
                    // Add the newly created employee to the local state
                    const newEmployee: Employee = {
                        id: response.data.id || response.data.user_id || `temp-${Date.now()}`,
                        full_name: values.full_name,
                        email: values.email,
                        role: role,
                        warehouse_id: values.warehouse_id || undefined,
                        warehouse_name: warehouseName,
                    };
                    setEmployees(prev => [...prev, newEmployee]);
                }

                toast.success(`${role} ${editingEmployee ? 'updated' : 'created'} successfully`);
                setIsModalOpen(false);
                setEditingEmployee(null);
                formik.resetForm();
            } catch (error: any) {
                toast.error(error.response?.data?.error || 'Failed to save employee');
            }
        },
    });

    const handleEditEmployee = (employee: Employee) => {
        setEditingEmployee(employee);
        setSelectedRole(employee.role as 'Supervisor' | 'Staff' | 'Auditor');
        formik.setValues({
            full_name: employee.full_name,
            email: employee.email,
            password: '',
            warehouse_id: employee.warehouse_id || '',
        });
        setIsModalOpen(true);
    };

    const handleDeleteEmployee = (employee: Employee) => {
        setDeletingEmployee(employee);
    };

    const confirmDelete = async () => {
        if (!deletingEmployee) return;
        setIsDeleting(true);
        try {
            // Use appropriate delete endpoint based on role
            const endpoint = `/manager/${deletingEmployee.role.toLowerCase()}/delete/${deletingEmployee.id}`;
            await api.delete(endpoint);

            // Remove from local state immediately
            setEmployees(prev => prev.filter(emp => emp.id !== deletingEmployee.id));

            toast.success('Employee deleted successfully');
            setDeletingEmployee(null);
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to delete employee');
        } finally {
            setIsDeleting(false);
        }
    };

    const openCreateModal = () => {
        setEditingEmployee(null);
        formik.resetForm();
        setIsModalOpen(true);
    };

    // Kanban columns
    const auditors = employees.filter(e => e.role === 'Auditor');
    const supervisors = employees.filter(e => e.role === 'Supervisor');
    const staff = employees.filter(e => e.role === 'Staff');

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Employees</h1>
                <button
                    onClick={openCreateModal}
                    className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 transition-colors duration-200"
                >
                    Add Employee
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Auditors Column */}
                <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg">
                    <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4 flex items-center">
                        <HiIdentification className="mr-2" /> Auditors
                    </h2>
                    <div className="space-y-3">
                        {auditors.map(emp => (
                            <div key={emp.id} className="bg-white dark:bg-dark-surface p-3 rounded shadow">
                                <div className="flex justify-between items-start">
                                    <div className="flex-1">
                                        <p className="font-medium text-gray-900 dark:text-white">{emp.full_name}</p>
                                        <p className="text-sm text-gray-500">{emp.email}</p>
                                    </div>
                                    <div className="flex space-x-2">
                                        <button
                                            onClick={() => handleEditEmployee(emp)}
                                            className="text-gray-400 hover:text-blue-500 transition-colors duration-150"
                                            title="Edit"
                                        >
                                            <HiPencil className="h-4 w-4" />
                                        </button>
                                        <button
                                            onClick={() => handleDeleteEmployee(emp)}
                                            className="text-gray-400 hover:text-red-500 transition-colors duration-150"
                                            title="Delete"
                                        >
                                            <HiTrash className="h-4 w-4" />
                                        </button>
                                    </div>
                                </div>
                            </div>
                        ))}
                        {auditors.length === 0 && <p className="text-gray-500 text-sm">No auditors found</p>}
                    </div>
                </div>

                {/* Supervisors Column */}
                <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg">
                    <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4 flex items-center">
                        <HiUser className="mr-2" /> Supervisors
                    </h2>
                    <div className="space-y-3">
                        {supervisors.map(emp => (
                            <div key={emp.id} className="bg-white dark:bg-dark-surface p-3 rounded shadow">
                                <div className="flex justify-between items-start">
                                    <div className="flex-1">
                                        <p className="font-medium text-gray-900 dark:text-white">{emp.full_name}</p>
                                        <p className="text-sm text-gray-500">{emp.email}</p>
                                        {emp.warehouse_id && <p className="text-xs text-gray-400 mt-1">Warehouse: {warehouses.find(w => w.id === emp.warehouse_id)?.name || emp.warehouse_id}</p>}
                                    </div>
                                    <div className="flex space-x-2">
                                        <button
                                            onClick={() => handleEditEmployee(emp)}
                                            className="text-gray-400 hover:text-blue-500 transition-colors duration-150"
                                            title="Edit"
                                        >
                                            <HiPencil className="h-4 w-4" />
                                        </button>
                                        <button
                                            onClick={() => handleDeleteEmployee(emp)}
                                            className="text-gray-400 hover:text-red-500 transition-colors duration-150"
                                            title="Delete"
                                        >
                                            <HiTrash className="h-4 w-4" />
                                        </button>
                                    </div>
                                </div>
                            </div>
                        ))}
                        {supervisors.length === 0 && <p className="text-gray-500 text-sm">No supervisors found</p>}
                    </div>
                </div>

                {/* Staff Column */}
                <div className="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg">
                    <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4 flex items-center">
                        <HiUserGroup className="mr-2" /> Staff
                    </h2>
                    <div className="space-y-3">
                        {staff.map(emp => (
                            <div key={emp.id} className="bg-white dark:bg-dark-surface p-3 rounded shadow">
                                <div className="flex justify-between items-start">
                                    <div className="flex-1">
                                        <p className="font-medium text-gray-900 dark:text-white">{emp.full_name}</p>
                                        <p className="text-sm text-gray-500">{emp.email}</p>
                                        {emp.warehouse_id && <p className="text-xs text-gray-400 mt-1">Warehouse: {warehouses.find(w => w.id === emp.warehouse_id)?.name || emp.warehouse_id}</p>}
                                    </div>
                                    <div className="flex space-x-2">
                                        <button
                                            onClick={() => handleEditEmployee(emp)}
                                            className="text-gray-400 hover:text-blue-500 transition-colors duration-150"
                                            title="Edit"
                                        >
                                            <HiPencil className="h-4 w-4" />
                                        </button>
                                        <button
                                            onClick={() => handleDeleteEmployee(emp)}
                                            className="text-gray-400 hover:text-red-500 transition-colors duration-150"
                                            title="Delete"
                                        >
                                            <HiTrash className="h-4 w-4" />
                                        </button>
                                    </div>
                                </div>
                            </div>
                        ))}
                        {staff.length === 0 && <p className="text-gray-500 text-sm">No staff found</p>}
                    </div>
                </div>
            </div>

            {/* Add Employee Modal */}
            <Dialog open={isModalOpen} onClose={() => setIsModalOpen(false)} className="fixed z-10 inset-0 overflow-y-auto">
                <div className="flex items-center justify-center min-h-screen">
                    <Dialog.Overlay className="fixed inset-0 bg-black opacity-30" />
                    <div className="relative bg-white dark:bg-dark-surface rounded-lg p-8 max-w-md w-full mx-4">
                        <Dialog.Title className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                            {editingEmployee ? 'Edit Employee' : 'Add New Employee'}
                        </Dialog.Title>

                        <form onSubmit={formik.handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Role</label>
                                <select
                                    value={selectedRole}
                                    onChange={(e) => setSelectedRole(e.target.value as any)}
                                    disabled={!!editingEmployee}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    <option value="Staff">Staff</option>
                                    <option value="Supervisor">Supervisor</option>
                                    <option value="Auditor">Auditor</option>
                                </select>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Full Name</label>
                                <input
                                    type="text"
                                    {...formik.getFieldProps('full_name')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Email</label>
                                <input
                                    type="email"
                                    {...formik.getFieldProps('email')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                            </div>

                            {!editingEmployee && (
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
                                    <input
                                        type="password"
                                        {...formik.getFieldProps('password')}
                                        className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    />
                                    {formik.touched.password && formik.errors.password && (
                                        <div className="text-red-500 text-xs mt-1">{formik.errors.password}</div>
                                    )}
                                </div>
                            )}

                            {selectedRole !== 'Auditor' && (
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Assign Warehouse</label>
                                    <select
                                        {...formik.getFieldProps('warehouse_id')}
                                        className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    >
                                        <option value="">Select Warehouse</option>
                                        {warehouses.map(w => (
                                            <option key={w.id} value={w.id}>{w.name}</option>
                                        ))}
                                    </select>
                                </div>
                            )}

                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    type="button"
                                    onClick={() => setIsModalOpen(false)}
                                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 transition-colors duration-200"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 transition-colors duration-200"
                                >
                                    {editingEmployee ? 'Update' : 'Create'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </Dialog>

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={!!deletingEmployee}
                onClose={() => setDeletingEmployee(null)}
                onConfirm={confirmDelete}
                title="Delete Employee"
                message={`Are you sure you want to delete "${deletingEmployee?.full_name}"? This action cannot be undone.`}
                confirmText="Delete"
                isLoading={isDeleting}
            />
        </div>
    );
};

export default Employees;
