import React, { useEffect, useState } from 'react';
import api from '../../services/api';
import { HiSearch, HiRefresh } from 'react-icons/hi';

interface AuditLog {
    id: string;
    action: string;
    resource_type: string;
    user_id: string;
    username: string;
    timestamp: string;
    status?: string;
    details?: any;
}

const AuditorAuditLogs: React.FC = () => {
    const [logs, setLogs] = useState<AuditLog[]>([]);
    const [loading, setLoading] = useState(true);
    const [filterAction, setFilterAction] = useState('');
    const [filterResource, setFilterResource] = useState('');

    const fetchLogs = async () => {
        setLoading(true);
        try {
            const params = new URLSearchParams();
            if (filterAction) params.append('action', filterAction);
            if (filterResource) params.append('resource_type', filterResource);

            const res = await api.get(`/auditor/audit-logs?${params.toString()}`);
            const logsData = res.data.logs || res.data;
            setLogs(Array.isArray(logsData) ? logsData : []);
        } catch (error) {
            console.error("Failed to fetch logs", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchLogs();
    }, []);

    const formatDetails = (details: any) => {
        if (!details) return '-';
        if (typeof details === 'string') return details;

        // Extract key information from details object
        const parts: string[] = [];
        if (details.email) parts.push(`Email: ${details.email}`);
        if (details.role) parts.push(`Role: ${details.role}`);
        if (details.name) parts.push(`Name: ${details.name}`);
        if (details.warehouse_name) parts.push(`Warehouse: ${details.warehouse_name}`);
        if (details.sku) parts.push(`SKU: ${details.sku}`);
        if (details.quantity !== undefined) parts.push(`Qty: ${details.quantity}`);
        if (details.warehouse_id) parts.push(`Warehouse ID: ${details.warehouse_id}`);

        return parts.length > 0 ? parts.join(', ') : JSON.stringify(details);
    };

    const formatTime = (timestamp: string) => {
        try {
            const date = new Date(timestamp);
            return date.toLocaleString('en-US', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit'
            });
        } catch {
            return timestamp;
        }
    };

    const getActionBadgeColor = (action: string) => {
        switch (action?.toUpperCase()) {
            case 'LOGIN':
                return 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200';
            case 'CREATE':
                return 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200';
            case 'UPDATE':
                return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200';
            case 'DELETE':
                return 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200';
            case 'READ':
                return 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200';
            default:
                return 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200';
        }
    };

    return (
        <div>
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-6 gap-4">
                <h1 className="text-2xl font-semibold text-gray-900 dark:text-white">Audit Logs (Read-Only)</h1>

                <div className="flex flex-wrap gap-2">
                    <input
                        type="text"
                        placeholder="Filter by Action..."
                        value={filterAction}
                        onChange={(e) => setFilterAction(e.target.value)}
                        className="rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white sm:text-sm"
                    />
                    <input
                        type="text"
                        placeholder="Filter by Resource..."
                        value={filterResource}
                        onChange={(e) => setFilterResource(e.target.value)}
                        className="rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white sm:text-sm"
                    />
                    <button
                        onClick={fetchLogs}
                        className="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                    >
                        <HiSearch className="mr-2" /> Filter
                    </button>
                    <button
                        onClick={() => { setFilterAction(''); setFilterResource(''); fetchLogs(); }}
                        className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600 dark:hover:bg-gray-600"
                    >
                        <HiRefresh className="mr-2" /> Reset
                    </button>
                </div>
            </div>

            <div className="bg-white dark:bg-dark-surface shadow overflow-hidden sm:rounded-lg">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                        <thead className="bg-gray-50 dark:bg-gray-800">
                            <tr>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                    Time
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                    Action
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                    Resource
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                    User
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                    Status
                                </th>
                                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                    Details
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white dark:bg-dark-surface divide-y divide-gray-200 dark:divide-gray-700">
                            {logs.map((log) => (
                                <tr key={log.id}>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                        {formatTime(log.timestamp)}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${getActionBadgeColor(log.action)}`}>
                                            {log.action}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                        {log.resource_type}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                        {log.username || log.user_id}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                                        {log.status === 'SUCCESS' ? (
                                            <span className="text-green-600 dark:text-green-400">✓ {log.status}</span>
                                        ) : log.status === 'FAILED' ? (
                                            <span className="text-red-600 dark:text-red-400">✗ {log.status}</span>
                                        ) : (
                                            <span className="text-gray-500">-</span>
                                        )}
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500 dark:text-gray-400 max-w-xs truncate">
                                        {formatDetails(log.details)}
                                    </td>
                                </tr>
                            ))}
                            {logs.length === 0 && !loading && (
                                <tr>
                                    <td colSpan={6} className="px-6 py-4 text-center text-sm text-gray-500 dark:text-gray-400">
                                        No logs found.
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
};

export default AuditorAuditLogs;
