'use client';

import React, { createContext, useEffect, useState } from 'react';
import axios from 'axios';
import { useRouter } from 'next/navigation';

// Define User Type
interface User {
    id?: string;
    name?: string;
    email?: string;
}

// Define Context Type
interface UserContextType {
    user: User | null;
    setUser: React.Dispatch<React.SetStateAction<User | null>>;
    isAuthenticated: boolean;
    setIsAuthenticated: React.Dispatch<React.SetStateAction<boolean>>;
    token: string | null;
    setToken: React.Dispatch<React.SetStateAction<string | null>>;
    logout: () => void;
}

export const userContext = createContext<UserContextType>({
    user: null,
    setUser: () => {},
    isAuthenticated: false,
    setIsAuthenticated: () => {},
    token: null,
    setToken: () => {},
    logout: () => {},
});

export const UserWrapper = ({ children }: { children: React.ReactNode }) => {
    const router = useRouter();
    const [isAuthenticated, setIsAuthenticated] = useState(false);
    const [user, setUser] = useState<User | null>(null);
    const [token, setToken] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const initializeAuth = async () => {
            const storedToken = localStorage.getItem('token');
            if (!storedToken) {
                setLoading(false);
                return;
            }
            setToken(storedToken);
            
            try {
                const { data } = await axios.get('http://localhost:8080/auth/validate', {
                    headers: { Authorization: `Bearer ${storedToken}` },
                });
                setUser(data);
                setIsAuthenticated(true);
            } catch {
                logout();
            } finally {
                setLoading(false);
            }
        };

        initializeAuth();
    }, []);

    useEffect(() => {
        if (!loading && !token) {
            router.push('/auth/sign-in');
        }
    }, [loading, token, router]);

    const logout = () => {
        localStorage.removeItem('token');
        setUser(null);
        setToken(null);
        setIsAuthenticated(false);
        router.push('/auth/sign-in');
    };

    return (
        <userContext.Provider value={{ user, setUser, isAuthenticated, setIsAuthenticated, token, setToken, logout }}>
            {!loading && children}
        </userContext.Provider>
    );
};
