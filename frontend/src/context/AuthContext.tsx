import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import api from '../services/api';
import { User, LoginResponse } from '../types/auth';

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (tokens: LoginResponse) => void;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const initAuth = async () => {
            const token = localStorage.getItem('access_token');
            if (token) {
                try {
                    // Decode token to get user info (assuming claims are in token)
                    // Or fetch profile from backend
                    // const decoded: any = jwtDecode(token);

                    // For MVP, we might need to fetch the full profile if token doesn't have everything
                    // But let's try to use what we have or fetch /users/me
                    try {
                        const response = await api.get('/users/me');
                        setUser(response.data);
                    } catch (error) {
                        console.error("Failed to fetch profile", error);
                        // If fetch fails but token is valid, maybe just logout
                        logout();
                    }
                } catch (error) {
                    console.error("Invalid token", error);
                    logout();
                }
            }
            setIsLoading(false);
        };

        initAuth();
    }, []);

    const login = (data: LoginResponse) => {
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token);
        setUser(data.user);
    };

    const logout = () => {
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        setUser(null);
        // Optional: Call backend logout endpoint
        api.post('/auth/logout').catch(err => console.error("Logout error", err));
        window.location.href = '/login';
    };

    return (
        <AuthContext.Provider value={{ user, isAuthenticated: !!user, isLoading, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
