import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import apiClient from '../api/client';
import { useAuth } from './AuthContext';

interface TeamPermission {
    team_id: number;
    team_name: string;
    role: string;
    can_edit: boolean;
    can_delete: boolean;
    can_manage: boolean;
}

interface UserPermissions {
    user_id: number;
    email: string;
    full_name: string;
    teams: TeamPermission[];
    is_org_admin: boolean;
}

interface PermissionContextType {
    permissions: UserPermissions | null;
    loading: boolean;
    refreshPermissions: () => Promise<void>;
    getTeamPermission: (teamId: number) => TeamPermission | null;
    canEdit: (teamId: number) => boolean;
    canDelete: (teamId: number) => boolean;
    canManage: (teamId: number) => boolean;
    isStakeholder: (teamId: number) => boolean;
    isMember: (teamId: number) => boolean;
    isManager: (teamId: number) => boolean;
}

const PermissionContext = createContext<PermissionContextType | undefined>(undefined);

export const PermissionProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { user, isAuthenticated } = useAuth();
    const [permissions, setPermissions] = useState<UserPermissions | null>(null);
    const [loading, setLoading] = useState(true);

    const refreshPermissions = useCallback(async () => {
        if (!isAuthenticated) {
            setPermissions(null);
            setLoading(false);
            return;
        }

        try {
            const res = await apiClient.get('/users/me/permissions');
            setPermissions(res.data);
        } catch (error) {
            console.error('Failed to fetch permissions:', error);
            setPermissions(null);
        } finally {
            setLoading(false);
        }
    }, [isAuthenticated]);

    useEffect(() => {
        refreshPermissions();
    }, [refreshPermissions, user]);

    const getTeamPermission = useCallback((teamId: number): TeamPermission | null => {
        if (!permissions) return null;
        return permissions.teams.find(t => t.team_id === teamId) || null;
    }, [permissions]);

    const canEdit = useCallback((teamId: number): boolean => {
        const perm = getTeamPermission(teamId);
        return perm?.can_edit ?? false;
    }, [getTeamPermission]);

    const canDelete = useCallback((teamId: number): boolean => {
        const perm = getTeamPermission(teamId);
        return perm?.can_delete ?? false;
    }, [getTeamPermission]);

    const canManage = useCallback((teamId: number): boolean => {
        const perm = getTeamPermission(teamId);
        return perm?.can_manage ?? false;
    }, [getTeamPermission]);

    const isStakeholder = useCallback((teamId: number): boolean => {
        const perm = getTeamPermission(teamId);
        return perm?.role === 'stakeholder';
    }, [getTeamPermission]);

    const isMember = useCallback((teamId: number): boolean => {
        const perm = getTeamPermission(teamId);
        return perm?.role === 'member';
    }, [getTeamPermission]);

    const isManager = useCallback((teamId: number): boolean => {
        const perm = getTeamPermission(teamId);
        return perm?.role === 'manager' || perm?.role === 'assistant';
    }, [getTeamPermission]);

    return (
        <PermissionContext.Provider value={{
            permissions,
            loading,
            refreshPermissions,
            getTeamPermission,
            canEdit,
            canDelete,
            canManage,
            isStakeholder,
            isMember,
            isManager,
        }}>
            {children}
        </PermissionContext.Provider>
    );
};

export const usePermissions = () => {
    const context = useContext(PermissionContext);
    if (!context) {
        throw new Error('usePermissions must be used within a PermissionProvider');
    }
    return context;
};
