import apiClient from './client';

export interface RegisterData {
    email: string;
    password: string;
    full_name: string;
    organization_id: number;
    timezone?: string;
}

export interface LoginData {
    email: string;
    password: string;
}

export interface AuthResponse {
    token: string;
    user: {
        id: number;
        email: string;
        full_name: string;
        organization_id: number;
        timezone: string;
    };
}

export const authAPI = {
    register: async (data: RegisterData): Promise<AuthResponse> => {
        const response = await apiClient.post('/auth/register', data);
        return response.data;
    },

    login: async (data: LoginData): Promise<AuthResponse> => {
        const response = await apiClient.post('/auth/login', data);
        return response.data;
    },

    logout: async (): Promise<void> => {
        await apiClient.post('/auth/logout');
    },
};
