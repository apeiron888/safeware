import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../../services/api';
import { Item } from '../../types';
import ItemCard from '../../components/common/ItemCard';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import EmptyState from '../../components/common/EmptyState';
import { HiArrowLeft } from 'react-icons/hi';
import { toast } from 'react-toastify';

const AuditorWarehouseDetails: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const [items, setItems] = useState<Item[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchItems = async () => {
            try {
                const res = await api.get(`/auditor/items/warehouse/${id}`);
                setItems(res.data || []);
            } catch (error) {
                console.error('Failed to fetch warehouse items', error);
                toast.error('Failed to load warehouse items');
            } finally {
                setLoading(false);
            }
        };

        if (id) {
            fetchItems();
        }
    }, [id]);

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
            <div className="mb-6">
                <button
                    onClick={() => navigate('/auditor/warehouses')}
                    className="flex items-center text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white mb-4 transition-colors duration-200"
                >
                    <HiArrowLeft className="mr-2" />
                    Back to Warehouses
                </button>
                <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Warehouse Items (Read-Only)</h1>
            </div>

            {/* Items Grid */}
            {items.length === 0 ? (
                <EmptyState
                    title="No items in this warehouse"
                    description="This warehouse currently has no items"
                />
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {items.map((item) => (
                        <ItemCard key={item.id} item={item} readOnly={true} />
                    ))}
                </div>
            )}
        </div>
    );
};

export default AuditorWarehouseDetails;
