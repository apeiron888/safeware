import React, { useEffect, useState } from 'react';
import api from '../../services/api';
import { HiOfficeBuilding, HiLocationMarker, HiTrash, HiPencil } from 'react-icons/hi';
import { Dialog } from '@headlessui/react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { toast } from 'react-toastify';
import { Link } from 'react-router-dom';

interface Warehouse {
    id: string;
    name: string;
    location: string;
    capacity?: number;
    supervisor_id?: string;
}

const Warehouses: React.FC = () => {
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingWarehouse, setEditingWarehouse] = useState<Warehouse | null>(null);

    const fetchWarehouses = async () => {
        try {
            const res = await api.get('/manager/summary/warehouses');
            const warehousesData = res.data.warehouses || res.data;
            setWarehouses(Array.isArray(warehousesData) ? warehousesData : []);
        } catch (error) {
            console.error("Failed to fetch warehouses", error);
            setWarehouses([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchWarehouses();
    }, []);

    const handleDelete = async (id: string) => {
        if (window.confirm('Are you sure you want to delete this warehouse?')) {
            try {
                await api.delete(`/manager/warehouse/delete/${id}`);
                // Remove from local state
                setWarehouses(prev => prev.filter(w => w.id !== id));
                toast.success('Warehouse deleted successfully');
            } catch (error: any) {
                toast.error(error.response?.data?.error || 'Failed to delete warehouse');
            }
        }
    };

    const formik = useFormik({
        initialValues: {
            name: '',
            location: '',
        },
        validationSchema: Yup.object({
            name: Yup.string().required('Required'),
            location: Yup.string().required('Required'),
        }),
        onSubmit: async (values) => {
            try {
                if (editingWarehouse) {
                    await api.patch(`/manager/warehouse/update/${editingWarehouse.id}`, values);
                    // Update local state
                    setWarehouses(prev => prev.map(w =>
                        w.id === editingWarehouse.id ? { ...w, ...values } : w
                    ));
                    toast.success('Warehouse updated successfully');
                } else {
                    const response = await api.post('/manager/warehouse/create', values);
                    // Add to local state
                    if (response.data) {
                        const newWarehouse = {
                            id: response.data.id || response.data.warehouse_id || `temp-${Date.now()}`,
                            ...values
                        };
                        setWarehouses(prev => [...prev, newWarehouse]);
                    }
                    toast.success('Warehouse created successfully');
                }
                setIsModalOpen(false);
                setEditingWarehouse(null);
                formik.resetForm();
            } catch (error: any) {
                toast.error(error.response?.data?.error || 'Failed to save warehouse');
            }
        },
    });

    const openEditModal = (warehouse: Warehouse) => {
        setEditingWarehouse(warehouse);
        formik.setValues({
            name: warehouse.name,
            location: warehouse.location,
        });
        setIsModalOpen(true);
    };

    const openCreateModal = () => {
        setEditingWarehouse(null);
        formik.resetForm();
        setIsModalOpen(true);
    };

    return (
        <div>
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Warehouses</h1>
                <button
                    onClick={openCreateModal}
                    className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
                >
                    Add Warehouse
                </button>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {warehouses.map((warehouse) => (
                    <div key={warehouse.id} className="bg-white dark:bg-dark-surface rounded-lg shadow overflow-hidden hover:shadow-md transition-shadow duration-200">
                        <div className="p-5">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center">
                                    <div className="flex-shrink-0 bg-primary-100 dark:bg-primary-900 rounded-md p-3">
                                        <HiOfficeBuilding className="h-6 w-6 text-primary-600 dark:text-primary-200" />
                                    </div>
                                    <div className="ml-4">
                                        <h3 className="text-lg font-medium text-gray-900 dark:text-white">{warehouse.name}</h3>
                                        <div className="flex items-center text-sm text-gray-500 dark:text-gray-400 mt-1">
                                            <HiLocationMarker className="mr-1" /> {warehouse.location}
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="mt-4 border-t border-gray-200 dark:border-gray-700 pt-4 flex justify-between items-center">
                                <Link to={`/manager/warehouse/${warehouse.id}`} className="text-sm font-medium text-primary-600 hover:text-primary-500">
                                    View Details
                                </Link>
                                <div className="flex space-x-2">
                                    <button onClick={() => openEditModal(warehouse)} className="text-gray-400 hover:text-blue-500">
                                        <HiPencil className="h-5 w-5" />
                                    </button>
                                    <button onClick={() => handleDelete(warehouse.id)} className="text-gray-400 hover:text-red-500">
                                        <HiTrash className="h-5 w-5" />
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                ))}
                {warehouses.length === 0 && !loading && (
                    <div className="col-span-full text-center py-10 text-gray-500">
                        No warehouses found. Create one to get started.
                    </div>
                )}
            </div>

            {/* Modal */}
            <Dialog open={isModalOpen} onClose={() => setIsModalOpen(false)} className="fixed z-10 inset-0 overflow-y-auto">
                <div className="flex items-center justify-center min-h-screen">
                    <Dialog.Overlay className="fixed inset-0 bg-black opacity-30" />
                    <div className="relative bg-white dark:bg-dark-surface rounded-lg p-8 max-w-md w-full mx-4">
                        <Dialog.Title className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                            {editingWarehouse ? 'Edit Warehouse' : 'Add New Warehouse'}
                        </Dialog.Title>

                        <form onSubmit={formik.handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Name</label>
                                <input
                                    type="text"
                                    {...formik.getFieldProps('name')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {formik.touched.name && formik.errors.name && (
                                    <div className="text-red-500 text-xs mt-1">{formik.errors.name}</div>
                                )}
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Location</label>
                                <input
                                    type="text"
                                    {...formik.getFieldProps('location')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {formik.touched.location && formik.errors.location && (
                                    <div className="text-red-500 text-xs mt-1">{formik.errors.location}</div>
                                )}
                            </div>

                            <div className="flex justify-end space-x-3 mt-6">
                                <button
                                    type="button"
                                    onClick={() => setIsModalOpen(false)}
                                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700"
                                >
                                    {editingWarehouse ? 'Update' : 'Create'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </Dialog>
        </div>
    );
};

export default Warehouses;
