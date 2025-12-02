import React, { useEffect, useState } from 'react';
import api from '../../services/api';
import { HiOfficeBuilding, HiCube, HiCurrencyDollar } from 'react-icons/hi';

interface Summary {
    total_warehouses: number;
    total_items: number;
    total_value: number;
}

const ManagerDashboard: React.FC = () => {
    const [summary, setSummary] = useState<Summary | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchSummary = async () => {
            try {
                // Fetch warehouses and items to calculate stats
                const [warehousesRes, itemsRes] = await Promise.all([
                    api.get('/manager/summary/warehouses'),
                    api.get('/manager/items/all')
                ]);

                // Handle different response structures
                const warehousesData = warehousesRes.data.warehouses || warehousesRes.data;
                const itemsData = itemsRes.data.items || itemsRes.data;

                const warehouses = Array.isArray(warehousesData) ? warehousesData : [];
                const items = Array.isArray(itemsData) ? itemsData : [];

                console.log('Dashboard data:', { warehouses: warehouses.length, items: items.length });

                const totalValue = items.reduce((sum: number, item: any) => {
                    const itemValue = (item.price || 0) * (item.quantity || 0);
                    return sum + itemValue;
                }, 0);

                setSummary({
                    total_warehouses: warehouses.length,
                    total_items: items.length,
                    total_value: totalValue
                });
            } catch (error) {
                console.error("Failed to fetch dashboard data", error);
                // Set default values on error
                setSummary({
                    total_warehouses: 0,
                    total_items: 0,
                    total_value: 0
                });
            } finally {
                setLoading(false);
            }
        };

        fetchSummary();
    }, []);

    if (loading) return <div>Loading stats...</div>;

    return (
        <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-6">Manager Dashboard</h1>

            <div className="grid grid-cols-1 gap-5 sm:grid-cols-3">
                <div className="bg-white dark:bg-dark-surface overflow-hidden shadow rounded-lg">
                    <div className="p-5">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <HiOfficeBuilding className="h-6 w-6 text-gray-400" aria-hidden="true" />
                            </div>
                            <div className="ml-5 w-0 flex-1">
                                <dl>
                                    <dt className="text-sm font-medium text-gray-500 truncate">Total Warehouses</dt>
                                    <dd>
                                        <div className="text-lg font-medium text-gray-900 dark:text-white">{summary?.total_warehouses}</div>
                                    </dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="bg-white dark:bg-dark-surface overflow-hidden shadow rounded-lg">
                    <div className="p-5">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <HiCube className="h-6 w-6 text-gray-400" aria-hidden="true" />
                            </div>
                            <div className="ml-5 w-0 flex-1">
                                <dl>
                                    <dt className="text-sm font-medium text-gray-500 truncate">Total Items</dt>
                                    <dd>
                                        <div className="text-lg font-medium text-gray-900 dark:text-white">{summary?.total_items}</div>
                                    </dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="bg-white dark:bg-dark-surface overflow-hidden shadow rounded-lg">
                    <div className="p-5">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <HiCurrencyDollar className="h-6 w-6 text-gray-400" aria-hidden="true" />
                            </div>
                            <div className="ml-5 w-0 flex-1">
                                <dl>
                                    <dt className="text-sm font-medium text-gray-500 truncate">Total Value</dt>
                                    <dd>
                                        <div className="text-lg font-medium text-gray-900 dark:text-white">${summary?.total_value.toFixed(2)}</div>
                                    </dd>
                                </dl>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ManagerDashboard;
