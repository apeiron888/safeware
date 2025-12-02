import React from 'react';
import { HiPencil, HiTrash } from 'react-icons/hi';
import { Item } from '../../types';

interface ItemCardProps {
    item: Item;
    onEdit?: (item: Item) => void;
    onDelete?: (item: Item) => void;
    readOnly?: boolean;
}

const ItemCard: React.FC<ItemCardProps> = ({ item, onEdit, onDelete, readOnly = false }) => {
    return (
        <div className="bg-white dark:bg-dark-surface rounded-lg shadow hover:shadow-md transition-shadow duration-200 overflow-hidden">
            <div className="p-4">
                <div className="flex justify-between items-start mb-2">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white truncate">
                        {item.name}
                    </h3>
                    {!readOnly && (onEdit || onDelete) && (
                        <div className="flex space-x-2 ml-2">
                            {onEdit && (
                                <button
                                    onClick={() => onEdit(item)}
                                    className="text-gray-400 hover:text-blue-500 transition-colors duration-150"
                                    aria-label="Edit item"
                                >
                                    <HiPencil className="h-5 w-5" />
                                </button>
                            )}
                            {onDelete && (
                                <button
                                    onClick={() => onDelete(item)}
                                    className="text-gray-400 hover:text-red-500 transition-colors duration-150"
                                    aria-label="Delete item"
                                >
                                    <HiTrash className="h-5 w-5" />
                                </button>
                            )}
                        </div>
                    )}
                </div>

                {item.sku && (
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-2">
                        SKU: {item.sku}
                    </p>
                )}

                {item.department && (
                    <p className="text-sm text-gray-600 dark:text-gray-400 mb-3">
                        {item.department}
                    </p>
                )}

                <div className="space-y-2">
                    {item.quality && (
                        <div className="flex justify-between items-center">
                            <span className="text-sm text-gray-500 dark:text-gray-400">Quality:</span>
                            <span className={`text-sm font-medium px-2 py-1 rounded ${item.quality === 'New' ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200' :
                                    item.quality === 'Used' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200' :
                                        'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-200'
                                }`}>
                                {item.quality}
                            </span>
                        </div>
                    )}
                    <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-500 dark:text-gray-400">Quantity:</span>
                        <span className="text-sm font-medium text-gray-900 dark:text-white">
                            {item.quantity}
                        </span>
                    </div>
                    <div className="flex justify-between items-center">
                        <span className="text-sm text-gray-500 dark:text-gray-400">Price:</span>
                        <span className="text-sm font-semibold text-primary-600 dark:text-primary-400">
                            ${item.price.toFixed(2)}
                        </span>
                    </div>
                    <div className="border-t border-gray-200 dark:border-gray-700 pt-2 mt-2">
                        <div className="flex justify-between items-center">
                            <span className="text-sm font-medium text-gray-500 dark:text-gray-400">Total Value:</span>
                            <span className="text-base font-bold text-secondary-600 dark:text-secondary-400">
                                ${(item.quantity * item.price).toFixed(2)}
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ItemCard;
