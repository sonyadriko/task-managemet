import React, { useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { useTheme } from '../contexts/ThemeContext';
import apiClient from '../api/client';
import Sidebar from '../components/Sidebar';
import './Settings.css';

const Settings: React.FC = () => {
    const { user } = useAuth();
    const { theme, toggleTheme } = useTheme();
    const [orgName, setOrgName] = useState('');
    const [passwordForm, setPasswordForm] = useState({
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
    });
    const [passwordError, setPasswordError] = useState('');
    const [passwordSuccess, setPasswordSuccess] = useState('');
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        const fetchOrgName = async () => {
            if (!user?.organization_id) return;
            try {
                const res = await apiClient.get(`/organizations/${user.organization_id}`);
                setOrgName(res.data?.name || 'Unknown');
            } catch {
                setOrgName('Unknown');
            }
        };
        fetchOrgName();
    }, [user?.organization_id]);

    const handlePasswordChange = async (e: React.FormEvent) => {
        e.preventDefault();
        setPasswordError('');
        setPasswordSuccess('');

        // Validate
        if (passwordForm.newPassword.length < 6) {
            setPasswordError('New password must be at least 6 characters');
            return;
        }

        if (passwordForm.newPassword !== passwordForm.confirmPassword) {
            setPasswordError('New passwords do not match');
            return;
        }

        setLoading(true);
        try {
            await apiClient.post('/auth/change-password', {
                current_password: passwordForm.currentPassword,
                new_password: passwordForm.newPassword
            });
            setPasswordSuccess('Password changed successfully!');
            setPasswordForm({ currentPassword: '', newPassword: '', confirmPassword: '' });
        } catch (error: any) {
            setPasswordError(error.response?.data?.error || 'Failed to change password');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="page-layout">
            <Sidebar activeItem="settings" />
            <main className="main-content">
                <header className="page-header">
                    <div className="header-left">
                        <h1>Settings</h1>
                        <p>Manage your account and preferences</p>
                    </div>
                </header>

                <div className="settings-grid">
                    {/* Profile Section */}
                    <section className="settings-section">
                        <h2>👤 Profile</h2>
                        <div className="settings-card">
                            <div className="profile-header">
                                <div className="profile-avatar">
                                    {user?.full_name?.charAt(0) || 'U'}
                                </div>
                                <div className="profile-info">
                                    <h3>{user?.full_name}</h3>
                                    <p>{user?.email}</p>
                                </div>
                            </div>
                            <div className="form-group">
                                <label className="form-label">Full Name</label>
                                <input type="text" className="form-input" value={user?.full_name || ''} readOnly />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Email</label>
                                <input type="email" className="form-input" value={user?.email || ''} readOnly />
                            </div>
                            <div className="form-group">
                                <label className="form-label">Timezone</label>
                                <input type="text" className="form-input" value={user?.timezone || ''} readOnly />
                            </div>
                        </div>
                    </section>

                    {/* Organization Section */}
                    <section className="settings-section">
                        <h2>🏢 Organization</h2>
                        <div className="settings-card">
                            <div className="info-row">
                                <span className="info-label">Organization</span>
                                <span className="info-value">{orgName || 'Loading...'}</span>
                            </div>
                            <div className="info-row">
                                <span className="info-label">Role</span>
                                <span className="info-value badge badge-primary">Member</span>
                            </div>
                        </div>
                    </section>

                    {/* Security Section */}
                    <section className="settings-section">
                        <h2>🔒 Security</h2>
                        <div className="settings-card">
                            {passwordError && (
                                <div className="alert alert-error">
                                    <span>⚠️</span> {passwordError}
                                </div>
                            )}
                            {passwordSuccess && (
                                <div className="alert alert-success">
                                    <span>✅</span> {passwordSuccess}
                                </div>
                            )}
                            <form onSubmit={handlePasswordChange}>
                                <div className="form-group">
                                    <label className="form-label">Current Password</label>
                                    <input
                                        type="password"
                                        className="form-input"
                                        placeholder="Enter current password"
                                        value={passwordForm.currentPassword}
                                        onChange={(e) => setPasswordForm({ ...passwordForm, currentPassword: e.target.value })}
                                        required
                                    />
                                </div>
                                <div className="form-group">
                                    <label className="form-label">New Password</label>
                                    <input
                                        type="password"
                                        className="form-input"
                                        placeholder="Enter new password (min 6 chars)"
                                        value={passwordForm.newPassword}
                                        onChange={(e) => setPasswordForm({ ...passwordForm, newPassword: e.target.value })}
                                        required
                                        minLength={6}
                                    />
                                </div>
                                <div className="form-group">
                                    <label className="form-label">Confirm New Password</label>
                                    <input
                                        type="password"
                                        className="form-input"
                                        placeholder="Confirm new password"
                                        value={passwordForm.confirmPassword}
                                        onChange={(e) => setPasswordForm({ ...passwordForm, confirmPassword: e.target.value })}
                                        required
                                    />
                                </div>
                                <button type="submit" className="btn btn-primary" disabled={loading}>
                                    {loading ? 'Changing...' : 'Change Password'}
                                </button>
                            </form>
                        </div>
                    </section>

                    {/* Preferences Section */}
                    <section className="settings-section">
                        <h2>🎨 Preferences</h2>
                        <div className="settings-card">
                            <div className="preference-row">
                                <div>
                                    <h4>Dark Mode</h4>
                                    <p>Use dark theme for the interface</p>
                                </div>
                                <label className="toggle">
                                    <input
                                        type="checkbox"
                                        checked={theme === 'dark'}
                                        onChange={toggleTheme}
                                    />
                                    <span className="toggle-slider"></span>
                                </label>
                            </div>
                            <div className="preference-row">
                                <div>
                                    <h4>Email Notifications</h4>
                                    <p>Receive email updates for task changes</p>
                                </div>
                                <label className="toggle">
                                    <input type="checkbox" readOnly />
                                    <span className="toggle-slider"></span>
                                </label>
                            </div>
                        </div>
                    </section>
                </div>
            </main>
        </div>
    );
};

export default Settings;
