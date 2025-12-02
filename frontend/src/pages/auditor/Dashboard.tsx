import React from 'react';
import { useNavigate } from 'react-router-dom';
import { HiOfficeBuilding, HiClipboardList } from 'react-icons/hi';

const AuditorDashboard: React.FC = () => {
    const navigate = useNavigate();

    return (
        <div>
            <h1 className="text-2xl font-semibold text-gray-900 dark:text-white mb-6">Auditor Dashboard</h1>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl">
                {/* Warehouses Card */}
                <button
                    onClick={() => navigate('/auditor/warehouses')}
                    className="group bg-white dark:bg-dark-surface rounded-lg shadow hover:shadow-lg transition-all duration-200 p-8 text-left"
                >
                    <div className="flex items-center justify-center w-16 h-16 bg-primary-100 dark:bg-primary-900 rounded-lg mb-4 group-hover:scale-110 transition-transform duration-200">
                        <HiOfficeBuilding className="h-8 w-8 text-primary-600 dark:text-primary-400" />
                    </div>
                    <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">View Warehouses</h2>
                    <p className="text-gray-600 dark:text-gray-400">
                        Browse all warehouses and their inventory in read-only mode
                    </p>
                </button>

                {/* Audit Logs Card */}
                <button
                    onClick={() => navigate('/auditor/logs')}
                    className="group bg-white dark:bg-dark-surface rounded-lg shadow hover:shadow-lg transition-all duration-200 p-8 text-left"
                >
                    <div className="flex items-center justify-center w-16 h-16 bg-secondary-100 dark:bg-secondary-900 rounded-lg mb-4 group-hover:scale-110 transition-transform duration-200">
                        <HiClipboardList className="h-8 w-8 text-secondary-600 dark:text-secondary-400" />
                    </div>
                    <h2 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">View Audit Logs</h2>
                    <p className="text-gray-600 dark:text-gray-400">
                        Review all system activity and changes in the audit trail
                    </p>
                </button>
            </div>
        </div>
    );
};

export default AuditorDashboard;
