export interface User {
    id: string;
    email: string;
    full_name: string;
    role: 'Manager' | 'Supervisor' | 'Staff' | 'Auditor';
    warehouse_id?: string;
    company_id: string;
}

export interface LoginResponse {
    access_token: string;
    refresh_token: string;
    user: User;
}

export interface AuthState {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
}
