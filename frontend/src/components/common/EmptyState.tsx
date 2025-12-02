import React from 'react';
import { HiInboxIn } from 'react-icons/hi';

interface EmptyStateProps {
    title: string;
    description?: string;
    icon?: React.ReactNode;
    action?: {
        label: string;
        onClick: () => void;
    };
}

const EmptyState: React.FC<EmptyStateProps> = ({ title, description, icon, action }) => {
    return (
        <div className="text-center py-12">
            <div className="flex justify-center mb-4">
                {icon || <HiInboxIn className="h-12 w-12 text-gray-400 dark:text-gray-600" />}
            </div>
            <h3 className="text-lg font-medium text-gray-900 dark:text-white mb-2">{title}</h3>
            {description && (
                <p className="text-sm text-gray-500 dark:text-gray-400 mb-4">{description}</p>
            )}
            {action && (
                <button
                    onClick={action.onClick}
                    className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
                >
                    {action.label}
                </button>
            )}
        </div>
    );
};

export default EmptyState;
