import React, { useEffect, useState } from 'react';
import api from '../../services/api';
import { Item } from '../../types';
import ItemCard from '../../components/common/ItemCard';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import EmptyState from '../../components/common/EmptyState';
import ConfirmDialog from '../../components/common/ConfirmDialog';
import { Dialog } from '@headlessui/react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { toast } from 'react-toastify';
import { HiPlus } from 'react-icons/hi';

const StaffDashboard: React.FC = () => {
    const [items, setItems] = useState<Item[]>([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [editingItem, setEditingItem] = useState<Item | null>(null);
    const [deletingItem, setDeletingItem] = useState<Item | null>(null);
    const [isDeleting, setIsDeleting] = useState(false);

    const fetchItems = async () => {
        try {
            const res = await api.get('/staff/items');
            const itemsData = res.data.items || res.data;
            setItems(Array.isArray(itemsData) ? itemsData : []);
        } catch (error) {
            console.error('Failed to fetch items', error);
            toast.error('Failed to load items');
            setItems([]);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchItems();
    }, []);

    const formik = useFormik({
        initialValues: {
            sku: '',
            name: '',
            quality: 'New',
            quantity: 0,
            price: 0,
            department: '',
        },
        validationSchema: Yup.object({
            sku: Yup.string().required('Required'),
            name: Yup.string().required('Required'),
            quality: Yup.string().required('Required'),
            quantity: Yup.number().min(0, 'Must be positive').required('Required'),
            price: Yup.number().min(0, 'Must be positive').required('Required'),
            department: Yup.string(),
        }),
        onSubmit: async (values) => {
            try {
                if (editingItem) {
                    await api.put(`/staff/item/update/${editingItem.id}`, values);
                    toast.success('Item updated successfully');
                } else {
                    await api.post('/staff/item/add', values);
                    toast.success('Item created successfully');
                }
                setIsModalOpen(false);
                setEditingItem(null);
                formik.resetForm();
                fetchItems();
            } catch (error: any) {
                toast.error(error.response?.data?.error || 'Failed to save item');
            }
        },
    });

    const handleEdit = (item: Item) => {
        setEditingItem(item);
        formik.setValues({
            sku: item.sku || '',
            name: item.name,
            quality: item.quality || 'New',
            quantity: item.quantity,
            price: item.price,
            department: item.department || '',
        });
        setIsModalOpen(true);
    };

    const handleDelete = (item: Item) => {
        setDeletingItem(item);
    };

    const confirmDelete = async () => {
        if (!deletingItem) return;
        setIsDeleting(true);
        try {
            await api.delete(`/staff/item/remove/${deletingItem.id}`);
            toast.success('Item deleted successfully');
            setDeletingItem(null);
            fetchItems();
        } catch (error: any) {
            toast.error(error.response?.data?.error || 'Failed to delete item');
        } finally {
            setIsDeleting(false);
        }
    };

    const openCreateModal = () => {
        setEditingItem(null);
        formik.resetForm();
        setIsModalOpen(true);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    return (
        <div>
            {/* Header */}
            <div className="flex justify-between items-center mb-6">
                <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Staff Dashboard - Warehouse Items</h1>
                <button
                    onClick={openCreateModal}
                    className="flex items-center px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 transition-colors duration-200"
                >
                    <HiPlus className="mr-2" />
                    Add Item
                </button>
            </div>

            {/* Items Grid */}
            {items.length === 0 ? (
                <EmptyState
                    title="No items in your warehouse"
                    description="Get started by adding your first item"
                    action={{
                        label: 'Add Item',
                        onClick: openCreateModal,
                    }}
                />
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {items.map((item) => (
                        <ItemCard
                            key={item.id}
                            item={item}
                            onEdit={handleEdit}
                            onDelete={handleDelete}
                        />
                    ))}
                </div>
            )}

            {/* Item Form Modal */}
            <Dialog open={isModalOpen} onClose={() => setIsModalOpen(false)} className="fixed z-10 inset-0 overflow-y-auto">
                <div className="flex items-center justify-center min-h-screen">
                    <Dialog.Overlay className="fixed inset-0 bg-black opacity-30" />
                    <div className="relative bg-white dark:bg-dark-surface rounded-lg p-8 max-w-md w-full mx-4">
                        <Dialog.Title className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                            {editingItem ? 'Edit Item' : 'Add New Item'}
                        </Dialog.Title>

                        <form onSubmit={formik.handleSubmit} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">SKU</label>
                                <input
                                    type="text"
                                    {...formik.getFieldProps('sku')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                                {formik.touched.sku && formik.errors.sku && (
                                    <div className="text-red-500 text-xs mt-1">{formik.errors.sku}</div>
                                )}
                            </div>

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
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Quality</label>
                                <select
                                    {...formik.getFieldProps('quality')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                >
                                    <option value="New">New</option>
                                    <option value="Used">Used</option>
                                    <option value="Damaged">Damaged</option>
                                </select>
                                {formik.touched.quality && formik.errors.quality && (
                                    <div className="text-red-500 text-xs mt-1">{formik.errors.quality}</div>
                                )}
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Quantity</label>
                                    <input
                                        type="number"
                                        {...formik.getFieldProps('quantity')}
                                        className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    />
                                    {formik.touched.quantity && formik.errors.quantity && (
                                        <div className="text-red-500 text-xs mt-1">{formik.errors.quantity}</div>
                                    )}
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Price</label>
                                    <input
                                        type="number"
                                        step="0.01"
                                        {...formik.getFieldProps('price')}
                                        className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                    />
                                    {formik.touched.price && formik.errors.price && (
                                        <div className="text-red-500 text-xs mt-1">{formik.errors.price}</div>
                                    )}
                                </div>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300">Department (Optional)</label>
                                <input
                                    type="text"
                                    {...formik.getFieldProps('department')}
                                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
                                />
                            </div>

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
                                    {editingItem ? 'Update' : 'Create'}
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            </Dialog>

            {/* Delete Confirmation */}
            <ConfirmDialog
                isOpen={!!deletingItem}
                onClose={() => setDeletingItem(null)}
                onConfirm={confirmDelete}
                title="Delete Item"
                message={`Are you sure you want to delete "${deletingItem?.name}"? This action cannot be undone.`}
                confirmText="Delete"
                isLoading={isDeleting}
            />
        </div>
    );
};

export default StaffDashboard;
