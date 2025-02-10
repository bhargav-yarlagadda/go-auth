'use client'

import React, { useContext, useState } from "react";
import { useRouter } from "next/navigation";
import axios from "axios";
import { userContext } from "./UserWrapper";

const SignIn = () => {

  const router = useRouter();
    const {setIsAuthenticated} = useContext(userContext)
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };


  const authenticateUser = async () => {
    try {
        const resp = await axios.post("http://localhost:8080/auth/sign-in", formData);
        const data = resp.data; // No need for `await` here

        if (resp.status === 200) {
            setIsAuthenticated(true);
            router.push('/');
        } else {
            console.error("Unexpected response status:", resp.status);
        }
    } catch (error) {
        console.error("Sign-up failed:", error);
    }
};

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await authenticateUser()
    // TODO: Add authentication logic
  };

  return (
    <div className="flex justify-center items-center h-screen text-black bg-gray-100">
      <div className="bg-white p-6 rounded-lg shadow-md w-80">
        <h2 className="text-2xl font-semibold text-center mb-4">Sign In</h2>
        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label className="block text-gray-700">Email</label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
              required
            />
          </div>
          <div className="mb-4">
            <label className="block text-gray-700">Password</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-400"
              required
            />
          </div>
          <button
            type="submit"
            className="w-full bg-blue-500 text-white py-2 rounded-lg hover:bg-blue-600 transition"
          >
            Sign In
          </button>
        </form>
        <p className="text-center text-sm text-gray-600 mt-4">
          Don't have an account?{" "}
          <button
            onClick={() => router.push("/auth/sign-up")}
            className="text-blue-500 hover:underline"
          >
            Sign Up
          </button>
        </p>
      </div>
    </div>
  );
};

export default SignIn;
