import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../services/api';
import { Warehouse } from '../../types';
import LoadingSpinner from '../../components/common/LoadingSpinner';
import EmptyState from '../../components/common/EmptyState';
import { HiOfficeBuilding, HiLocationMarker, HiEye } from 'react-icons/hi';
import { toast } from 'react-toastify';

const AuditorWarehouses: React.FC = () => {
    const navigate = useNavigate();
    const [warehouses, setWarehouses] = useState<Warehouse[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchWarehouses = async () => {
            try {
                const res = await api.get('/auditor/warehouses');
                const warehousesData = res.data.warehouses || res.data;
                setWarehouses(Array.isArray(warehousesData) ? warehousesData : []);
            } catch (error) {
                console.error('Failed to fetch warehouses', error);
                toast.error('Failed to load warehouses');
                setWarehouses([]);
            } finally {
                setLoading(false);
            }
        };

        fetchWarehouses();
    }, []);

    if (loading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <LoadingSpinner size="lg" />
            </div>
        );
    }

    return (
        <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-6">Warehouses (Read-Only)</h1>

            {warehouses.length === 0 ? (
                <EmptyState
                    title="No warehouses found"
                    description="There are no warehouses to display at this time"
                />
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {warehouses.map((warehouse) => (
                        <div
                            key={warehouse.id}
                            className="bg-white dark:bg-dark-surface rounded-lg shadow hover:shadow-md transition-shadow duration-200 overflow-hidden cursor-pointer"
                            onClick={() => navigate(`/auditor/warehouse/${warehouse.id}`)}
                        >
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
                                {warehouse.items_count !== undefined && (
                                    <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                                        <div className="flex justify-between items-center text-sm">
                                            <span className="text-gray-500 dark:text-gray-400">Total Items:</span>
                                            <span className="font-medium text-gray-900 dark:text-white">{warehouse.items_count}</span>
                                        </div>
                                        {warehouse.total_value !== undefined && (
                                            <div className="flex justify-between items-center text-sm mt-2">
                                                <span className="text-gray-500 dark:text-gray-400">Total Value:</span>
                                                <span className="font-semibold text-primary-600 dark:text-primary-400">
                                                    ${warehouse.total_value.toFixed(2)}
                                                </span>
                                            </div>
                                        )}
                                    </div>
                                )}
                                <div className="mt-4 flex items-center text-sm font-medium text-primary-600 dark:text-primary-400">
                                    <HiEye className="mr-1" />
                                    View Details
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default AuditorWarehouses;
