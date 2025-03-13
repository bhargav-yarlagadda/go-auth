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
    const [loading, setLoading] = useState(true); // Fix: Ensure loading state prevents redirection

    useEffect(() => {
        const storedToken = localStorage.getItem("token");
        if (storedToken) {
            setToken(storedToken);
        }
        setLoading(false); // Fix: Prevent redirect before token is retrieved
    }, []);

    useEffect(() => {
        const validateToken = async () => {
            if (!token) return;
            try {
                const response = await axios.get("http://localhost:8080/auth/validate", {
                    headers: { Authorization: `Bearer ${token}` }
                });
                setUser(response.data);
                setIsAuthenticated(true);
            } catch (error) {
                logout();
            }
        };

        if (token) validateToken();
    }, [token]);

    // ✅ Fix: Only redirect if loading is false and token is absent
    useEffect(() => {
        if (!loading && !token) {
            router.push("/auth/sign-in");
        }
    }, [loading, token]); 

    const logout = () => {
        localStorage.removeItem("token");
        setUser(null);
        setToken(null);
        setIsAuthenticated(false);
        router.push("/auth/sign-in");
    };

    return (
        <userContext.Provider value={{ user, setUser, isAuthenticated, setIsAuthenticated, token, setToken, logout }}>
            {!loading && children}
        </userContext.Provider>
    );
};
