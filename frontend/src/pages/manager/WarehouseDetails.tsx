import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/api';
import { Item, Warehouse } from '../../types';
import ItemCard from '../../components/common/ItemCard';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import EmptyState from '../../components/common/EmptyState';
import ConfirmDialog from '../../components/common/ConfirmDialog';
import { Dialog } from '@headlessui/react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { toast } from 'react-toastify';
import { HiArrowLeft, HiPlus, HiUser, HiTrash, HiArrowUp, HiArrowDown } from 'react-icons/hi';

interface Employee {
    id: string;
    full_name: string;
    email: string;
    role: string;
    warehouse_id?: string;
}

const WarehouseDetails: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [warehouse, setWarehouse] = useState<Warehouse | null>(null);
    const [items, setItems] = useState<Item[]>([]);
    const [employees, setEmployees] = useState<Employee[]>([]);
    const [loading, setLoading] = useState(true);
    const [isItemModalOpen, setIsItemModalOpen] = useState(false);
    const [isEmployeeModalOpen, setIsEmployeeModalOpen] = useState(false);
    const [editingItem, setEditingItem] = useState<Item | null>(null);
    const [deletingItem, setDeletingItem] = useState<Item | null>(null);
    const [deletingEmployee, setDeletingEmployee] = useState<Employee | null>(null);
    const [isDeleting, setIsDeleting] = useState(false);

    const fetchWarehouseDetails = async () => {
        try {
            // Fetch warehouse info
            const warehousesRes = await api.get('/manager/summary/warehouses');
            const warehousesData = warehousesRes.data.warehouses || warehousesRes.data;
            const currentWarehouse = Array.isArray(warehousesData)
                ? warehousesData.find((w: Warehouse) => w.id === id)
                : null;
            setWarehouse(currentWarehouse || null);

            // Fetch items for this warehouse
            const itemsRes = await api.get(`/manager/items/warehouse/${id}`);
            const itemsData = itemsRes.data.items || itemsRes.data;
            setItems(Array.isArray(itemsData) ? itemsData : []);

            // Fetch employees assigned to this warehouse
            const employeesRes = await api.get('/manager/employees');
            const allEmployees = Array.isArray(employeesRes.data) ? employeesRes.data : [];
            const warehouseEmployees = allEmployees.filter((emp: Employee) =>
                emp.warehouse_id === id && (emp.role === 'Supervisor' || emp.role === 'Staff')
            );
            setEmployees(warehouseEmployees);
        } catch (error) {
            console.error('Failed to fetch warehouse details', error);
            toast.error('Failed to load warehouse details');
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (id) {
            fetchWarehouseDetails();
        }
    }, [id]);

    const itemFormik = useFormik({
        initialValues: {
            sku: '',
            name: '',
            quality: 'New',
            quantity: 0,
            price: 0,
            department: '',
            batch: '',
        },
        validationSchema: Yup.object({
            sku: Yup.string().required('Required'),
            name: Yup.string().required('Required'),
            quality: Yup.string().required('Required'),
            quantity: Yup.number().min(0, 'Must be positive').required('Required'),
            price: Yup.number().min(0, 'Must be positive').required('Required'),
            department: Yup.string(),
            batch: Yup.string(),
        }),
        onSubmit: async (values) => {
            try {
                if (editingItem) {
                    await api.put(`/manager/item/update/${editingItem.id}`, values);
                    setItems(prev => prev.map(i => i.id === editingItem.id ? { ...i, ...values } : i));
                    toast.success('Item updated successfully');
                } else {
                    const response = await api.post('/manager/item/create', {
                        ...values,
                        warehouse_id: id!
                    });
                    if (response.data) {
                        const newItem: Item = {
                            id: response.data.id || response.data._id || response.data.item_id || `temp-${Date.now()}`,
                            ...values,
                            warehouse_id: id!
                        };
                        setItems(prev => [...prev, newItem]);
                        // Refresh to get complete item data from server
                        await fetchWarehouseDetails();
                    }
                    toast.success('Item created successfully');
                }
                setIsItemModalOpen(false);
                setEditingItem(null);
                itemFormik.resetForm();
            } catch (error: any) {
                toast.error(error.response?.data?.error || 'Failed to save item');
            }
        },
    });

    const employeeFormik = useFormik({
        initialValues: {
            full_name: '',
            email: '',
            password: '',
            role: 'Staff',
        },
        validationSchema: Yup.object({
            full_name: Yup.string().required('Required'),
            email: Yup.string().email('Invalid email').required('Required'),
            password: Yup.string().min(8, 'Min 8 characters').required('Required'),
            role: Yup.string().oneOf(['Staff', 'Supervisor']).required('Required'),
        }),
        onSubmit: async (values) => {
            try {
                const endpoint = `/manager/${values.role.toLowerCase()}/create`;
                const response = await api.post(endpoint, {
                    ...values,
                    warehouse_id: id,
                });

                if (response.data) {
                    const newEmployee: Employee = {
                        id: response.data.user_id || `temp-${Date.now()}`,
                        full_name: values.full_name,
                        email: values.email,
                        role: values.role,
                        warehouse_id: id,
                    };
                    setEmployees(prev => [...prev, newEmployee]);
                    toast.success(`${values.role} added successfully`);
                }
                setIsEmployeeModalOpen(false);
                employeeFormik.resetForm();
            } catch (error: any) {
                toast.error(error.response?.data?.error || 'Failed to add employee');
            }
        },
    });

    const handleEditItem = (item: Item) => {
        setEditingItem(item);
        itemFormik.setValues({
            sku: item.sku || '',
            name: item.name,
            quality: item.quality || 'New',
            quantity: item.quantity,
            price: item.price,
            department: item.department || '',
            batch: '',
        });
        setIsItemModalOpen(true);
    };

    const handleDeleteItem = (item: Item) => {
        setDeletingItem(item);
    };

    const confirmDeleteItem = async () => {
        if (!deletingItem) return;
        setIsDeleting(true);
        try {
            await api.delete(`/manager/item/remove/${deletingItem.id}`);
            setItems(prev => prev.filter(i => i.id !== deletingItem.id));
            toast.success('Item deleted successfully');
            setDeletingItem(null);
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to delete item');
        } finally {
            setIsDeleting(false);
        }
    };

    const handleDeleteEmployee = (employee: Employee) => {
        setDeletingEmployee(employee);
    };

    const confirmDeleteEmployee = async () => {
        if (!deletingEmployee) return;
        setIsDeleting(true);
        try {
            const endpoint = `/manager/${deletingEmployee.role.toLowerCase()}/delete/${deletingEmployee.id}`;
            await api.delete(endpoint);
            setEmployees(prev => prev.filter(emp => emp.id !== deletingEmployee.id));
            toast.success('Employee removed successfully');
            setDeletingEmployee(null);
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to remove employee');
        } finally {
            setIsDeleting(false);
        }
    };

    const handlePromote = async (employee: Employee) => {
        if (employee.role === 'Supervisor') {
            toast.info('Employee is already a Supervisor');
            return;
        }

        try {
            // This would require a promote endpoint in the backend
            // For now, we can show a message
            toast.info('Promotion feature coming soon');
            // TODO: Implement promotion endpoint
        } catch (error: any) {
            toast.error('Failed to promote employee');
        }
    };

    const handleDemote = async (employee: Employee) => {
        if (employee.role === 'Staff') {
            toast.info('Employee is already Staff');
            return;
        }

        try {
            // This would require a demote endpoint in the backend
            toast.info('Demotion feature coming soon');
            // TODO: Implement demotion endpoint
        } catch (error: any) {
            toast.error('Failed to demote employee');
        }
    };

    const openCreateItemModal = () => {
        setEditingItem(null);
        itemFormik.resetForm();
        setIsItemModalOpen(true);
    };

    const openAddEmployeeModal = () => {
        employeeFormik.resetForm();
        setIsEmployeeModalOpen(true);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    if (!warehouse) {
        return (
            <div className="text-center py-12">
                <h2 className="text-xl font-semibold text-gray-900 dark:text-white">Warehouse not found</h2>
                <button
                    onClick={() => navigate('/manager/warehouses')}
                    className="mt-4 text-primary-600 hover:text-primary-700"
                >
                    Back to Warehouses
                </button>
            </div>
        );
    }

    const supervisors = employees.filter(emp => emp.role === 'Supervisor');
    const staff = employees.filter(emp => emp.role === 'Staff');

    return (
        <div className="flex gap-6">
            {/* Main Content Area */}
            <div className="flex-1">
                {/* Header */}
                <div className="mb-6">
                    <button
                        onClick={() => navigate('/manager/warehouses')}
                        className="flex items-center text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white mb-4 transition-colors duration-200"
                    >
                        <HiArrowLeft className="mr-2" />
                        Back to Warehouses
                    </button>
                    <div className="flex justify-between items-start">
                        <div>
                            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">{warehouse.name}</h1>
                            <p className="text-gray-500 dark:text-gray-400 mt-1">{warehouse.location}</p>
                        </div>
                        <button
                            onClick={openCreateItemModal}
                            className="flex items-center px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 transition-colors duration-200"
                        >
                            <HiPlus className="mr-2" />
                            Add Item
                        </button>
                    </div>
                </div>

                {/* Items Grid */}
                {items.length === 0 ? (
                    <EmptyState
                        title="No items in this warehouse"
                        description="Get started by adding your first item"
                        action={{
                            label: 'Add Item',
                            onClick: openCreateItemModal,
                        }}
                    />
                ) : (
                    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {items.map((item) => (
                            <ItemCard
                                key={item.id}
                                item={item}
                                onEdit={handleEditItem}
                                onDelete={handleDeleteItem}
                            />
                        ))}
                    </div>
                )}
            </div>

            {/* Employee Sidebar */}
            <div className="w-80 flex-shrink-0">
                <div className="bg-white dark:bg-dark-surface rounded-lg shadow-md p-6 sticky top-6">
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
                            <HiUser className="mr-2" />
                            Employees
                        </h2>
                        <button
                            onClick={openAddEmployeeModal}
                            className="p-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 transition-colors duration-200"
                            title="Add Employee"
                        >
                            <HiPlus className="h-4 w-4" />
                        </button>
                    </div>

                    {/* Supervisors Section */}
                    <div className="mb-6">
                        <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Supervisors</h3>
                        {supervisors.length === 0 ? (
                            <p className="text-sm text-gray-500 dark:text-gray-400 italic">No supervisors assigned</p>
                        ) : (
                            <div className="space-y-2">
                                {supervisors.map(emp => (
                                    <div key={emp.id} className="bg-gray-50 dark:bg-gray-700 p-3 rounded-md">
                                        <div className="flex justify-between items-start">
                                            <div className="flex-1 min-w-0">
                                                <p className="font-medium text-gray-900 dark:text-white text-sm truncate">{emp.full_name}</p>
                                                <p className="text-xs text-gray-500 dark:text-gray-400 truncate">{emp.email}</p>
                                            </div>
                                            <div className="flex space-x-1 ml-2">
                                                <button
                                                    onClick={() => handleDemote(emp)}
                                                    className="p-1 text-gray-400 hover:text-orange-500 transition-colors"
                                                    title="Demote to Staff"
                                                >
                                                    <HiArrowDown className="h-4 w-4" />
                                                </button>
                                                <button
                                                    onClick={() => handleDeleteEmployee(emp)}
                                                    className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                                                    title="Remove Employee"
                                                >
                                                    <HiTrash className="h-4 w-4" />
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>

                    {/* Staff Section */}
                    <div>
                        <h3 className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Staff</h3>
                        {staff.length === 0 ? (
                            <p className="text-sm text-gray-500 dark:text-gray-400 italic">No staff assigned</p>
                        ) : (
                            <div className="space-y-2">
                                {staff.map(emp => (
                                    <div key={emp.id} className="bg-gray-50 dark:bg-gray-700 p-3 rounded-md">
                                        <div className="flex justify-between items-start">
                                            <div className="flex-1 min-w-0">
                                                <p className="font-medium text-gray-900 dark:text-white text-sm truncate">{emp.full_name}</p>
                                                <p className="text-xs text-gray-500 dark:text-gray-400 truncate">{emp.email}</p>
                                            </div>
                                            <div className="flex space-x-1 ml-2">
                                                <button
                                                    onClick={() => handlePromote(emp)}
                                                    className="p-1 text-gray-400 hover:text-green-500 transition-colors"
                                                    title="Promote to Supervisor"
                                                >
                                                    <HiArrowUp className="h-4 w-4" />
                                                </button>
                                                <button
                                                    onClick={() => handleDeleteEmployee(emp)}
                                                    className="p-1 text-gray-400 hover:text-red-500 transition-colors"
                                                    title="Remove Employee"
                                                >
                                                    <HiTrash className="h-4 w-4" />
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {/* Item Form Modal */}
            <Dialog open={isItemModalOpen} onClose={() => setIsItemModalOpen(false)} className="fixed z-10 inset-0 overflow-y-auto">
                <div className="flex items-center justify-center min-h-screen">
                    <Dialog.Overlay className="fixed inset-0 bg-black opacity-30" />
                    <div className="relative bg-white dark:bg-dark-surface rounded-lg p-8 max-w-md w-full mx-4">
                        <Dialog.Title className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                            {editingItem ? 'Edit Item' : 'Add New Item'}
                        </Dialog.Title>

                        <form onSubmit={itemFormik.handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">SKU</label>
                                <input
                                    type="text"
                                    {...itemFormik.getFieldProps('sku')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {itemFormik.touched.sku && itemFormik.errors.sku && (
                                    <div className="text-red-500 text-xs mt-1">{itemFormik.errors.sku}</div>
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Name</label>
                                <input
                                    type="text"
                                    {...itemFormik.getFieldProps('name')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {itemFormik.touched.name && itemFormik.errors.name && (
                                    <div className="text-red-500 text-xs mt-1">{itemFormik.errors.name}</div>
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Quality</label>
                                <select
                                    {...itemFormik.getFieldProps('quality')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                >
                                    <option value="New">New</option>
                                    <option value="Used">Used</option>
                                    <option value="Damaged">Damaged</option>
                                </select>
                                {itemFormik.touched.quality && itemFormik.errors.quality && (
                                    <div className="text-red-500 text-xs mt-1">{itemFormik.errors.quality}</div>
                                )}
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Quantity</label>
                                    <input
                                        type="number"
                                        {...itemFormik.getFieldProps('quantity')}
                                        className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    />
                                    {itemFormik.touched.quantity && itemFormik.errors.quantity && (
                                        <div className="text-red-500 text-xs mt-1">{itemFormik.errors.quantity}</div>
                                    )}
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Price</label>
                                    <input
                                        type="number"
                                        step="0.01"
                                        {...itemFormik.getFieldProps('price')}
                                        className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    />
                                    {itemFormik.touched.price && itemFormik.errors.price && (
                                        <div className="text-red-500 text-xs mt-1">{itemFormik.errors.price}</div>
                                    )}
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Department (Optional)</label>
                                <input
                                    type="text"
                                    {...itemFormik.getFieldProps('department')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Batch (Optional)</label>
                                <input
                                    type="text"
                                    {...itemFormik.getFieldProps('batch')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                            </div>

                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    type="button"
                                    onClick={() => setIsItemModalOpen(false)}
                                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 transition-colors duration-200"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 transition-colors duration-200"
                                >
                                    {editingItem ? 'Update' : 'Create'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </Dialog>

            {/* Add Employee Modal */}
            <Dialog open={isEmployeeModalOpen} onClose={() => setIsEmployeeModalOpen(false)} className="fixed z-10 inset-0 overflow-y-auto">
                <div className="flex items-center justify-center min-h-screen">
                    <Dialog.Overlay className="fixed inset-0 bg-black opacity-30" />
                    <div className="relative bg-white dark:bg-dark-surface rounded-lg p-8 max-w-md w-full mx-4">
                        <Dialog.Title className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                            Add Employee to Warehouse
                        </Dialog.Title>

                        <form onSubmit={employeeFormik.handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Full Name</label>
                                <input
                                    type="text"
                                    {...employeeFormik.getFieldProps('full_name')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {employeeFormik.touched.full_name && employeeFormik.errors.full_name && (
                                    <div className="text-red-500 text-xs mt-1">{employeeFormik.errors.full_name}</div>
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Email</label>
                                <input
                                    type="email"
                                    {...employeeFormik.getFieldProps('email')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {employeeFormik.touched.email && employeeFormik.errors.email && (
                                    <div className="text-red-500 text-xs mt-1">{employeeFormik.errors.email}</div>
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Password</label>
                                <input
                                    type="password"
                                    {...employeeFormik.getFieldProps('password')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {employeeFormik.touched.password && employeeFormik.errors.password && (
                                    <div className="text-red-500 text-xs mt-1">{employeeFormik.errors.password}</div>
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Role</label>
                                <select
                                    {...employeeFormik.getFieldProps('role')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                >
                                    <option value="Staff">Staff</option>
                                    <option value="Supervisor">Supervisor</option>
                                </select>
                            </div>

                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    type="button"
                                    onClick={() => setIsEmployeeModalOpen(false)}
                                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 transition-colors duration-200"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 transition-colors duration-200"
                                >
                                    Add Employee
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </Dialog>

            {/* Delete Item Confirmation */}
            <ConfirmDialog
                isOpen={!!deletingItem}
                onClose={() => setDeletingItem(null)}
                onConfirm={confirmDeleteItem}
                title="Delete Item"
                message={`Are you sure you want to delete "${deletingItem?.name}"? This action cannot be undone.`}
                confirmText="Delete"
                isLoading={isDeleting}
            />

            {/* Delete Employee Confirmation */}
            <ConfirmDialog
                isOpen={!!deletingEmployee}
                onClose={() => setDeletingEmployee(null)}
                onConfirm={confirmDeleteEmployee}
                title="Remove Employee"
                message={`Are you sure you want to remove "${deletingEmployee?.full_name}" from this warehouse?`}
                confirmText="Remove"
                isLoading={isDeleting}
            />
        </div>
    );
};

export default WarehouseDetails;
