'use client'

import React, { createContext, useEffect, useState } from 'react'
import axios from 'axios'
import { useRouter } from 'next/navigation'

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
}

export const userContext = createContext<UserContextType>({
    user: null,
    setUser: () => {},
    isAuthenticated: false,
    setIsAuthenticated: () => {},
});

export const UserWrapper = ({ children }: { children: React.ReactNode }) => {
    const router = useRouter();
    const [user, setUser] = useState<User | null>(null);
    const [isAuthenticated, setIsAuthenticated] = useState(false);

    useEffect(() => {
        if (!isAuthenticated) {
            router.push('/auth/sign-in')
        }
    }, [isAuthenticated, router])  
    return (
        <userContext.Provider value={{ user, setUser, isAuthenticated, setIsAuthenticated }}>
            {children}
        </userContext.Provider>
    );
};
